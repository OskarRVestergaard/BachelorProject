package PowPeer

import (
	"errors"
	"github.com/OskarRVestergaard/BachelorProject/production/Message"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy/PoW"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"github.com/google/uuid"
	"sync"
	"time"
)

/*
This is a single Peer that both listens to and sends messages
CURRENTLY IT ASSUMES THAT A PEER NEVER LEAVES AND TCP CONNECTIONS DON'T DROP
*/

type PoWPeer struct {
	signatureStrategy          signature_strategy.SignatureInterface
	lotteryStrategy            lottery_strategy.LotteryInterface //TODO Change to just use proof of work, and remove strategy, also proof of work should also send slot number along
	publicToSecret             chan map[string]string
	unfinalizedTransactions    chan []models.SignedPaymentTransaction
	blockTreeChan              chan PoWblockchain.Blocktree
	unhandledBlocks            chan PoWblockchain.Block
	unhandledMessages          chan Message.Message
	hardness                   int
	maximumTransactionsInBlock int
	network                    network.Network
	stopMiningSignal           chan struct{}
	isMiningMutex              sync.Mutex
	startTime                  time.Time
}

func (p *PoWPeer) RunPeer(IpPort string, startTime time.Time) {
	p.startTime = startTime
	p.signatureStrategy = signature_strategy.ECDSASig{}
	p.lotteryStrategy = &PoW.PoW{}
	address, err := network.StringToAddress(IpPort)
	if err != nil {
		panic("Could not parse IpPort: " + err.Error())
	}
	p.network = network.Network{}
	messagesFromNetwork := p.network.StartNetwork(address)

	p.stopMiningSignal = make(chan struct{})

	p.unfinalizedTransactions = make(chan []models.SignedPaymentTransaction, 1)
	p.unfinalizedTransactions <- make([]models.SignedPaymentTransaction, 0, 100)
	p.publicToSecret = make(chan map[string]string, 1)
	p.publicToSecret <- make(map[string]string)
	p.blockTreeChan = make(chan PoWblockchain.Blocktree, 1)
	newBlockTree, blockTreeCreationWentWell := PoWblockchain.NewBlocktree(PoWblockchain.CreateGenesisBlock())
	if !blockTreeCreationWentWell {
		panic("Could not generate new blocktree")
	}
	p.unhandledBlocks = make(chan PoWblockchain.Block, 20)
	p.hardness = newBlockTree.GetHead().BlockData.Hardness
	p.maximumTransactionsInBlock = constants.BlockPaymentAmountLimit
	p.unhandledMessages = make(chan Message.Message, 50)
	p.blockTreeChan <- newBlockTree

	go p.blockHandlerLoop()
	go p.messageHandlerLoop(messagesFromNetwork)
}

func (p *PoWPeer) CreateAccount() string {
	secretKey, publicKey := p.signatureStrategy.KeyGen()
	keys := <-p.publicToSecret
	keys[publicKey] = secretKey
	p.publicToSecret <- keys
	return publicKey
}

func (p *PoWPeer) GetAddress() network.Address {
	return p.network.GetAddress()
}

func (p *PoWPeer) Connect(ip string, port int) {
	addr := network.Address{
		Ip:   ip,
		Port: port,
	}
	ownIpPort := p.network.GetAddress().ToString()
	print(ownIpPort + " Connecting to " + addr.ToString() + "\n")
	err := p.network.SendMessageTo(Message.Message{MessageType: constants.JoinMessage, MessageSender: ownIpPort}, addr)

	if err != nil {
		panic(err.Error())
	}
}

func (p *PoWPeer) FloodSignedTransaction(from string, to string, amount int) {
	secretSigningKey, foundSecretKey := p.getSecretKey(from)
	if !foundSecretKey {
		return
	}
	trans := models.SignedPaymentTransaction{Id: uuid.New(), From: from, To: to, Amount: amount, Signature: nil}
	trans.SignTransaction(p.signatureStrategy, secretSigningKey)
	ipPort := p.network.GetAddress().ToString()
	msg := Message.Message{MessageType: constants.SignedTransaction, MessageSender: ipPort, SignedTransaction: trans}
	p.addTransaction(trans)
	p.network.FloodMessageToAllKnown(msg)
}

func (p *PoWPeer) messageHandlerLoop(incomingMessages chan Message.Message) {
	for {
		msg := <-incomingMessages
		p.handleMessage(msg)
	}
}

func (p *PoWPeer) handleMessage(msg Message.Message) {
	msgType := (msg).MessageType

	switch msgType {
	case constants.SignedTransaction:
		if utils.TransactionHasCorrectSignature(p.signatureStrategy, msg.SignedTransaction) {
			p.addTransaction(Message.MakeDeepCopyOfPayment(msg.SignedTransaction))
		}
	case constants.JoinMessage:

	case constants.BlockDelivery:
		for _, block := range msg.PoWMessageBlocks {
			p.unhandledBlocks <- block
		}
	default:
		println(p.network.GetAddress().ToString() + ": received a UNKNOWN message type ( " + msg.MessageType + " ) from: " + msg.MessageSender)
	}
}

/*
getSecretKey

returns the secret key associated with a given public key and return a boolean indicating whether the key is known
*/
func (p *PoWPeer) getSecretKey(pk string) (secretKey string, isKnownKey bool) {
	publicToSecret := <-p.publicToSecret
	secretSigningKey, foundSecretKey := publicToSecret[pk]
	p.publicToSecret <- publicToSecret
	if !foundSecretKey {
		return "", false
	}
	return secretSigningKey, true
}

func (p *PoWPeer) GetBlockTree() interface{} {
	//TODO Not at all thread safe to use it this way, fine if used for reading during testing
	blocktree := <-p.blockTreeChan
	p.blockTreeChan <- blocktree
	return blocktree
}

func (p *PoWPeer) StartMining(_ int) error {
	noActiveMiner := p.isMiningMutex.TryLock()
	if !noActiveMiner {
		return errors.New("peer is already mining")
	}
	secretKeys := <-p.publicToSecret
	verificationKey := utils.GetSomeKey(secretKeys)
	p.publicToSecret <- secretKeys
	blocktree := <-p.blockTreeChan
	newHeadHashes := blocktree.SubScribeToGetHead()
	head := blocktree.GetHead()
	initialHash := head.HashOfBlock()
	winningDraws := make(chan lottery_strategy.WinningLotteryParams, 10)
	p.lotteryStrategy.StartNewMiner(verificationKey, p.hardness, initialHash, newHeadHashes, winningDraws, p.stopMiningSignal)
	go p.blockCreatingLoop(winningDraws)

	p.blockTreeChan <- blocktree
	return nil
}

func (p *PoWPeer) StopMining() error {
	noActiveMiner := p.isMiningMutex.TryLock()
	if noActiveMiner {
		p.isMiningMutex.Unlock()
		return errors.New("peer is already not mining")
	}
	p.stopMiningSignal <- struct{}{}
	p.isMiningMutex.Unlock()
	return nil
}

func (p *PoWPeer) createBlock(verificationKey string, slot int, draw lottery_strategy.WinningLotteryParams, blocktree PoWblockchain.Blocktree) (newBlock PoWblockchain.Block, isEmpty bool) {
	//TODO Need to check that the draw is correct
	secretKey, foundSk := p.getSecretKey(verificationKey)
	if !foundSk {
		panic("Tried to create a block but peer did not have the associated SecretKey")
	}
	parentHash := draw.ParentHash
	unfinalizedTransactions := <-p.unfinalizedTransactions
	allTransactionsToAdd := blocktree.GetTransactionsNotInTree(unfinalizedTransactions)
	p.unfinalizedTransactions <- unfinalizedTransactions

	var transactionsToAdd []models.SignedPaymentTransaction
	if len(allTransactionsToAdd) <= p.maximumTransactionsInBlock {
		transactionsToAdd = allTransactionsToAdd
	}
	if len(allTransactionsToAdd) > p.maximumTransactionsInBlock {
		transactionsToAdd = make([]models.SignedPaymentTransaction, p.maximumTransactionsInBlock)
		for i := 0; i < p.maximumTransactionsInBlock; i++ {
			transactionsToAdd[i] = allTransactionsToAdd[i]
			//This could maybe cause starvation of transactions, if not enough blocks are made to saturate transaction demand
		}
	}
	//
	resultBlock := PoWblockchain.Block{
		IsGenesis: false,
		Vk:        verificationKey,
		Slot:      slot,
		Draw:      draw,
		BlockData: PoWblockchain.BlockData{
			Transactions: transactionsToAdd,
		},
		ParentHash: parentHash,
		Signature:  nil,
	}
	resultBlock.SignBlock(p.signatureStrategy, secretKey)
	if resultBlock.HasCorrectSignature(p.signatureStrategy) {
		return resultBlock, false
	} else {
		panic("Something went wrong, created block but gave it a wrong signature")
	}
}

func (p *PoWPeer) sendBlockWithTransactions(draw lottery_strategy.WinningLotteryParams) {
	secretKeys := <-p.publicToSecret
	verificationKey := utils.GetSomeKey(secretKeys) //todo maybe make sure that it is the same public key that was used for the draw
	p.publicToSecret <- secretKeys
	blocktree := <-p.blockTreeChan
	extendedOnSlot := blocktree.HashToBlock(draw.ParentHash).Slot
	slot := utils.CalculateSlot(p.startTime)
	for slot <= extendedOnSlot {
		time.Sleep(constants.SlotLength / 10)
		slot = utils.CalculateSlot(p.startTime)
	}
	blockWithTransactions, isEmpty := p.createBlock(verificationKey, slot, draw, blocktree)
	if isEmpty {
		p.blockTreeChan <- blocktree
		return
	}
	msg := Message.Message{
		MessageType:      constants.BlockDelivery,
		MessageSender:    p.network.GetAddress().ToString(),
		PoWMessageBlocks: []PoWblockchain.Block{blockWithTransactions},
	}
	for _, block := range msg.PoWMessageBlocks {
		p.unhandledBlocks <- block
	}
	p.blockTreeChan <- blocktree
	p.network.FloodMessageToAllKnown(msg)
}

func (p *PoWPeer) blockHandlerLoop() {
	for {
		blockToHandle := <-p.unhandledBlocks
		go p.handleBlock(blockToHandle)
	}
}

func (p *PoWPeer) verifyBlock(block PoWblockchain.Block) bool {
	//TODO Needs to verify that the transactions are not already present too (just like the sender did), since someone not following the protocol could exploit this
	//TODO This is potentially very slow, but could be faster using dynamic programming in the case the chain best chain does not switch often
	if !block.HasCorrectSignature(p.signatureStrategy) {
		return false
	}
	if !p.verifyTransactions(block.BlockData.Transactions) {
		return false
	}
	if block.Draw.Vk != block.Vk {
		return false
	}
	if block.Draw.ParentHash != block.ParentHash {
		return false //TODO Instance of new block (slot2) being sent with an old draw (slot1)
	}
	if !p.lotteryStrategy.Verify(block.Vk, block.ParentHash, p.hardness, block.Draw.Counter) {
		return false
	}
	return true
}

func (p *PoWPeer) verifyTransactions(transactions []models.SignedPaymentTransaction) bool {
	for _, transaction := range transactions {
		transactionSignatureIsCorrect := utils.TransactionHasCorrectSignature(p.signatureStrategy, transaction)
		if !transactionSignatureIsCorrect {
			return false
		}
	}
	return true
}

func (p *PoWPeer) handleBlock(block PoWblockchain.Block) {
	if !p.verifyBlock(block) {
		return
	}
	blocktree := <-p.blockTreeChan
	block = Message.MakeDeepCopyOfPoWBlock(block)
	var t = blocktree.AddBlock(block)
	switch t {
	case -3:
		//Slot number is not greater than parent
		p.blockTreeChan <- blocktree
	case -2:
		//Block with isGenesis true, not a real block and should be ignored
		p.blockTreeChan <- blocktree
	case -1:
		//Block is in tree already and can be ignored
		p.blockTreeChan <- blocktree
	case 0:
		//Parent is not in the tree, try to add later
		//TODO Maybe have another slice that are blocks which are waiting for parents to be added,
		//TODO such that they can be added immediately follow the parents addition to the tree (in case 1)

		p.blockTreeChan <- blocktree
		time.Sleep(1000 * time.Millisecond) //Needs to be enough time for the other block to arrive
		p.unhandledBlocks <- block
	case 1:
		//Block successfully added to the tree
		p.blockTreeChan <- blocktree
	default:
		p.blockTreeChan <- blocktree
		panic("addBlockReturnValueNotUnderstood")
	}
}

func (p *PoWPeer) addTransaction(t models.SignedPaymentTransaction) {
	unfinalizedTransactions := <-p.unfinalizedTransactions
	unfinalizedTransactions = append(unfinalizedTransactions, t)
	p.unfinalizedTransactions <- unfinalizedTransactions
}

func (p *PoWPeer) blockCreatingLoop(wins chan lottery_strategy.WinningLotteryParams) {
	for {
		newWin := <-wins
		go p.sendBlockWithTransactions(newWin)
	}
}

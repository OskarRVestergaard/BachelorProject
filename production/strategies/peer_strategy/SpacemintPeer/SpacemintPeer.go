package SpacemintPeer

import (
	"errors"
	"github.com/OskarRVestergaard/BachelorProject/Task1"
	"github.com/OskarRVestergaard/BachelorProject/production/Message"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/SpaceMintBlockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy/PoSpace"
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

type PoSpacePeer struct {
	signatureStrategy          signature_strategy.SignatureInterface
	lotteryStrategy            *PoSpace.PoSpace
	publicToSecret             chan map[string]string
	knownCommitments           chan map[string]SpaceMintBlockchain.SpaceCommitment
	unfinalizedTransactions    chan SpaceMintBlockchain.SpacemintTransactions
	blockTreeChan              chan SpaceMintBlockchain.Blocktree
	unhandledBlocks            chan SpaceMintBlockchain.Block
	unhandledMessages          chan Message.Message
	maximumPaymentsInBlock     int
	maximumSpaceCommitsInBlock int
	maximumPenaltiesInBlock    int
	network                    network.Network
	stopMiningSignal           chan struct{}
	isMiningMutex              sync.Mutex
	startTime                  time.Time
}

func (p *PoSpacePeer) RunPeer(IpPort string, startTime time.Time) {
	p.startTime = startTime
	p.signatureStrategy = signature_strategy.ECDSASig{}
	p.lotteryStrategy = &PoSpace.PoSpace{}
	address, err := network.StringToAddress(IpPort)
	if err != nil {
		panic("Could not parse IpPort: " + err.Error())
	}
	p.network = network.Network{}
	messagesFromNetwork := p.network.StartNetwork(address)

	p.stopMiningSignal = make(chan struct{})

	p.unfinalizedTransactions = make(chan SpaceMintBlockchain.SpacemintTransactions, 1)
	p.unfinalizedTransactions <- SpaceMintBlockchain.SpacemintTransactions{
		Payments:         []models.SignedPaymentTransaction{},
		SpaceCommitments: []SpaceMintBlockchain.SpaceCommitment{},
		Penalties:        []SpaceMintBlockchain.Penalty{},
	}
	p.knownCommitments = make(chan map[string]SpaceMintBlockchain.SpaceCommitment, 1)
	p.knownCommitments <- make(map[string]SpaceMintBlockchain.SpaceCommitment)
	p.publicToSecret = make(chan map[string]string, 1)
	p.publicToSecret <- make(map[string]string)
	p.blockTreeChan = make(chan SpaceMintBlockchain.Blocktree, 1)
	newBlockTree, blockTreeCreationWentWell := SpaceMintBlockchain.NewBlocktree(SpaceMintBlockchain.CreateGenesisBlock())
	if !blockTreeCreationWentWell {
		panic("Could not generate new blocktree")
	}
	p.unhandledBlocks = make(chan SpaceMintBlockchain.Block, 20)
	p.maximumPaymentsInBlock = constants.BlockPaymentAmountLimit
	p.maximumSpaceCommitsInBlock = constants.BlockSpaceCommitAmountLimit
	p.maximumPenaltiesInBlock = constants.BlockPenaltyAmountLimit
	p.unhandledMessages = make(chan Message.Message, 50)
	p.blockTreeChan <- newBlockTree

	go p.blockHandlerLoop()
	go p.messageHandlerLoop(messagesFromNetwork)
}

func (p *PoSpacePeer) CreateAccount() string {
	secretKey, publicKey := p.signatureStrategy.KeyGen()
	keys := <-p.publicToSecret
	keys[publicKey] = secretKey
	p.publicToSecret <- keys
	return publicKey
}

func (p *PoSpacePeer) GetAddress() network.Address {
	return p.network.GetAddress()
}

func (p *PoSpacePeer) Connect(ip string, port int) {
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

func (p *PoSpacePeer) floodSpaceCommit(commitment sha256.HashValue, id uuid.UUID, n int, pk string) {
	spaceTransaction := SpaceMintBlockchain.SpaceCommitment{
		Id:         id,
		N:          n,
		PublicKey:  pk,
		Commitment: commitment,
	}
	ipPort := p.network.GetAddress().ToString()
	msg := Message.Message{MessageType: constants.SpaceCommitTransaction, MessageSender: ipPort, SpaceCommitment: spaceTransaction}
	p.addSpaceCommit(Message.MakeDeepCopyOfSpaceCommit(spaceTransaction))
	p.network.FloodMessageToAllKnown(Message.MakeDeepCopyOfMessage(msg))
}

func (p *PoSpacePeer) FloodSignedTransaction(from string, to string, amount int) {
	secretSigningKey, foundSecretKey := p.getSecretKey(from)
	if !foundSecretKey {
		return
	}
	payment := models.SignedPaymentTransaction{Id: uuid.New(), From: from, To: to, Amount: amount, Signature: nil}
	payment.SignTransaction(p.signatureStrategy, secretSigningKey)
	ipPort := p.network.GetAddress().ToString()
	msg := Message.Message{MessageType: constants.SignedTransaction, MessageSender: ipPort, SignedTransaction: payment}
	p.addPayment(payment)
	p.network.FloodMessageToAllKnown(msg)
}

func (p *PoSpacePeer) messageHandlerLoop(incomingMessages chan Message.Message) {
	for {
		msg := <-incomingMessages
		p.handleMessage(Message.MakeDeepCopyOfMessage(msg))
	}
}

func (p *PoSpacePeer) handleMessage(msg Message.Message) {
	msgType := (msg).MessageType
	switch msgType {
	case constants.SignedTransaction:
		if utils.TransactionHasCorrectSignature(p.signatureStrategy, msg.SignedTransaction) {
			p.addPayment(msg.SignedTransaction)
		}
	case constants.SpaceCommitTransaction:
		p.addSpaceCommit(msg.SpaceCommitment)
	case constants.JoinMessage:

	case constants.BlockDelivery:
		for _, block := range msg.SpaceMintBlocks {
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
func (p *PoSpacePeer) getSecretKey(pk string) (secretKey string, isKnownKey bool) {
	publicToSecret := <-p.publicToSecret
	secretSigningKey, foundSecretKey := publicToSecret[pk]
	p.publicToSecret <- publicToSecret
	if !foundSecretKey {
		return "", false
	}
	return secretSigningKey, true
}

/*
GetBlockTree

For testing only, it is NOT thread safe, but is called after the blocktree does not change anymore, it can still be useful for testing
*/
func (p *PoSpacePeer) GetBlockTree() interface{} {
	blocktree := <-p.blockTreeChan
	p.blockTreeChan <- blocktree
	return blocktree
}

func (p *PoSpacePeer) startBlocksToMinePasser(initialMiningLocation PoSpace.MiningLocation, newMiningLocations chan PoSpace.MiningLocation) chan PoSpace.MiningLocation {
	mostRecentMiningLocation := make(chan PoSpace.MiningLocation, 1)
	mostRecentMiningLocation <- initialMiningLocation
	go func() {
		for {
			newLocation := <-newMiningLocations
			_ = <-mostRecentMiningLocation
			mostRecentMiningLocation <- newLocation
		}
	}()
	slotNotifier := utils.StartTimeSlotUpdater(p.startTime)
	hashesSendToMiner := make(chan PoSpace.MiningLocation, 10)
	go func() {
		for {
			_ = <-slotNotifier //TODO, Send slot along with hashes, instead of letting the miner calculate the slot itself
			hashToMineOn := <-mostRecentMiningLocation
			hashesSendToMiner <- hashToMineOn
			mostRecentMiningLocation <- hashToMineOn
		}
	}()
	return hashesSendToMiner
}

func (p *PoSpacePeer) StartMining(n int) error {
	noActiveMiner := p.isMiningMutex.TryLock()
	if !noActiveMiner {
		return errors.New("peer is already mining")
	}
	secretKeys := <-p.publicToSecret
	verificationKey := utils.GetSomeKey(secretKeys)
	p.publicToSecret <- secretKeys
	blocktree := <-p.blockTreeChan
	newMiningLocations := blocktree.SubScribeToGetHead(n)
	head := blocktree.GetHead()
	initialMiningLocation := blocktree.GetMiningLocation(head.HashOfBlock(), n)
	winningDraws := make(chan PoSpace.LotteryDraw, 10)
	poSpaceParameters := Task1.GenerateParameters()
	blocksToMiner := p.startBlocksToMinePasser(initialMiningLocation, newMiningLocations)
	commitment := p.lotteryStrategy.StartNewMiner(poSpaceParameters, verificationKey, 0, initialMiningLocation, blocksToMiner, winningDraws, p.stopMiningSignal)
	p.floodSpaceCommit(commitment, poSpaceParameters.Id, n, verificationKey)
	go p.blockCreatingLoop(winningDraws)

	p.blockTreeChan <- blocktree
	return nil
}

func (p *PoSpacePeer) StopMining() error {
	noActiveMiner := p.isMiningMutex.TryLock()
	if noActiveMiner {
		p.isMiningMutex.Unlock()
		return errors.New("peer is already not mining")
	}
	p.stopMiningSignal <- struct{}{}
	p.isMiningMutex.Unlock()
	return nil
}

func (p *PoSpacePeer) createBlock(verificationKey string, slot int, draw PoSpace.LotteryDraw, blocktree SpaceMintBlockchain.Blocktree) (newBlock SpaceMintBlockchain.Block, isEmpty bool) {
	//TODO Need to check that the draw is correct
	secretKey, foundSk := p.getSecretKey(verificationKey)
	if !foundSk {
		panic("Tried to create a block but peer did not have the associated SecretKey")
	}
	parentHash := draw.ParentHash
	unfinalizedTransactions := <-p.unfinalizedTransactions
	allTransactionsToAdd := blocktree.GetTransactionsNotInTree(unfinalizedTransactions)
	p.unfinalizedTransactions <- unfinalizedTransactions

	var paymentsToAdd = allTransactionsToAdd.Payments
	if len(allTransactionsToAdd.Payments) > p.maximumPaymentsInBlock {
		paymentsToAdd = allTransactionsToAdd.Payments[:p.maximumPaymentsInBlock]
	}
	var SpaceCommitsToAdd = allTransactionsToAdd.SpaceCommitments
	if len(allTransactionsToAdd.SpaceCommitments) > p.maximumSpaceCommitsInBlock {
		SpaceCommitsToAdd = allTransactionsToAdd.SpaceCommitments[:p.maximumSpaceCommitsInBlock]
	}
	//Same for penalty

	//
	resultBlock := SpaceMintBlockchain.Block{
		IsGenesis:  false,
		ParentHash: parentHash,
		HashSubBlock: SpaceMintBlockchain.HashSubBlock{
			Slot:                          slot,
			SignatureOnParentHashSubBlock: nil,
			Draw:                          draw,
		},
		TransactionSubBlock: SpaceMintBlockchain.TransactionSubBlock{
			Slot: slot,
			Transactions: SpaceMintBlockchain.SpacemintTransactions{
				Payments:         paymentsToAdd,
				SpaceCommitments: SpaceCommitsToAdd,
				Penalties:        []SpaceMintBlockchain.Penalty{},
			},
		},
		SignatureSubBlock: SpaceMintBlockchain.SignatureSubBlock{
			Slot:                                  slot,
			SignatureOnCurrentTransactionSubBlock: nil,
			SignatureOnParentSignatureSubBlock:    nil,
		},
	}

	parentBlock, isEmpty := blocktree.HashToBlock(parentHash)
	if isEmpty {
		panic("Something went wrong druing block creation, tried to create a block with no valid parent!")
	}
	resultBlock.SignBlock(parentBlock, p.signatureStrategy, secretKey)
	return resultBlock, false
}

func (p *PoSpacePeer) sendBlockWithTransactions(draw PoSpace.LotteryDraw) {
	secretKeys := <-p.publicToSecret
	verificationKey := utils.GetSomeKey(secretKeys) //todo maybe make sure that it is the same public key that was used for the draw
	p.publicToSecret <- secretKeys
	blocktree := <-p.blockTreeChan
	nod, isEmpty := blocktree.HashToBlock(draw.ParentHash)
	if isEmpty {
		panic("Trying to send block with no valid parent")
	}
	extendedOnSlot := nod.TransactionSubBlock.Slot
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
		MessageType:     constants.BlockDelivery,
		MessageSender:   p.network.GetAddress().ToString(),
		SpaceMintBlocks: []SpaceMintBlockchain.Block{blockWithTransactions},
	}
	for _, block := range msg.SpaceMintBlocks {
		p.unhandledBlocks <- block
	}
	p.blockTreeChan <- blocktree
	p.network.FloodMessageToAllKnown(msg)
}

func (p *PoSpacePeer) blockHandlerLoop() {
	for {
		blockToHandle := <-p.unhandledBlocks
		go p.handleBlock(blockToHandle)
	}
}

func (p *PoSpacePeer) verifyBlock(block SpaceMintBlockchain.Block) bool {
	//Ideally this also needs to verify that the transactions are not already present too (just like the sender did), since someone not following the protocol could exploit this
	//TODO This is potentially very slow, but could be faster using dynamic programming in the case the chain best chain does not switch often
	blockTree := <-p.blockTreeChan
	parentBlock, isEmpty := blockTree.HashToBlock(block.ParentHash)
	if isEmpty {
		p.unhandledBlocks <- block
		p.blockTreeChan <- blockTree
		time.Sleep(200 * time.Millisecond)
		return false
	}
	if !block.HasConsistentSignaturesAndSlots(parentBlock, p.signatureStrategy) {
		p.blockTreeChan <- blockTree
		return false
	}
	if !p.verifyTransactions(block.TransactionSubBlock.Transactions.Payments) {
		p.blockTreeChan <- blockTree
		return false
	}
	if !block.ParentHash.Equals(block.HashSubBlock.Draw.ParentHash) {
		p.blockTreeChan <- blockTree
		return false
	}
	knownCommitments := <-p.knownCommitments
	commitmentOfProof, isKnown := knownCommitments[block.HashSubBlock.Draw.Vk]
	if !isKnown {
		//We have not heard about this peer allocation this space, maybe it has just not arrived yet, is a full system, they would have to be delivered long before (contained is a prior block)
		p.unhandledBlocks <- block
		p.blockTreeChan <- blockTree
		p.knownCommitments <- knownCommitments
		time.Sleep(200 * time.Millisecond)
		return false
	}
	chalA, chalB := blockTree.GetChallengesForExtendingOnBlockWithHash(block.ParentHash, commitmentOfProof.N)
	location := PoSpace.MiningLocation{
		Slot:          block.HashSubBlock.Slot - 1,
		ParentHash:    block.ParentHash,
		ChallengeSetP: chalA,
		ChallengeSetV: chalB,
	}
	//TODO Discuss and handle real parameters
	prm := Task1.GenerateParameters()
	prm.Id = commitmentOfProof.Id
	if !p.lotteryStrategy.Verify(prm, block.HashSubBlock.Draw, location, commitmentOfProof.Commitment) {
		p.blockTreeChan <- blockTree
		p.knownCommitments <- knownCommitments
		return false
	}
	p.blockTreeChan <- blockTree
	p.knownCommitments <- knownCommitments
	return true
}

func (p *PoSpacePeer) verifyTransactions(transactions []models.SignedPaymentTransaction) bool {
	for _, transaction := range transactions {
		transactionSignatureIsCorrect := utils.TransactionHasCorrectSignature(p.signatureStrategy, transaction)
		if !transactionSignatureIsCorrect {
			return false
		}
	}
	return true
}

func (p *PoSpacePeer) handleBlock(block SpaceMintBlockchain.Block) {
	if !p.verifyBlock(block) {
		return
	}
	blocktree := <-p.blockTreeChan
	block = Message.MakeDeepCopyOfPoSBlock(block)
	knownCommitments := <-p.knownCommitments
	commitmentOfProof := knownCommitments[block.HashSubBlock.Draw.Vk]
	var t = blocktree.AddBlock(block, int64(commitmentOfProof.N))
	p.knownCommitments <- knownCommitments
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
		//This case is actually avoided because of the verification step earlier
		p.blockTreeChan <- blocktree
		time.Sleep(200 * time.Millisecond) //Needs to be enough time for the other block to arrive
		p.unhandledBlocks <- block
	case 1:
		//Block successfully added to the tree
		p.blockTreeChan <- blocktree
	default:
		p.blockTreeChan <- blocktree
		panic("addBlockReturnValueNotUnderstood")
	}
}

func (p *PoSpacePeer) addTransaction(t SpaceMintBlockchain.SpacemintTransactions) {
	unfinalizedTransactions := <-p.unfinalizedTransactions
	newTransaction := Message.MakeDeepCopyOfTransaction(t)
	unfinalizedTransactions.Payments = append(unfinalizedTransactions.Payments, newTransaction.Payments...)
	unfinalizedTransactions.SpaceCommitments = append(unfinalizedTransactions.SpaceCommitments, newTransaction.SpaceCommitments...)
	unfinalizedTransactions.Penalties = append(unfinalizedTransactions.Penalties, newTransaction.Penalties...)
	p.unfinalizedTransactions <- unfinalizedTransactions
}

func (p *PoSpacePeer) blockCreatingLoop(wins chan PoSpace.LotteryDraw) {
	for {
		newWin := <-wins
		go p.sendBlockWithTransactions(newWin)
	}
}

func (p *PoSpacePeer) addPayment(payment models.SignedPaymentTransaction) {
	transactionToAdd := SpaceMintBlockchain.SpacemintTransactions{
		Payments:         []models.SignedPaymentTransaction{payment},
		SpaceCommitments: []SpaceMintBlockchain.SpaceCommitment{},
		Penalties:        []SpaceMintBlockchain.Penalty{},
	}
	p.addTransaction(transactionToAdd)
}

func (p *PoSpacePeer) addSpaceCommit(spaceCommitment SpaceMintBlockchain.SpaceCommitment) {
	//"unsafe" fix to avoid the miner startup problem, space commits are accepted immediate, instead of just being added to the chain
	currentlyKnownCommitments := <-p.knownCommitments
	currentlyKnownCommitments[spaceCommitment.PublicKey] = spaceCommitment
	p.knownCommitments <- currentlyKnownCommitments
	//The above should be safe to omit/delete if finalization is made

	transactionToAdd := SpaceMintBlockchain.SpacemintTransactions{
		Payments:         []models.SignedPaymentTransaction{},
		SpaceCommitments: []SpaceMintBlockchain.SpaceCommitment{spaceCommitment},
		Penalties:        []SpaceMintBlockchain.Penalty{},
	}
	p.addTransaction(transactionToAdd)
}

package PowPeer

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"sync"
	"time"
)

/*
This is a single Peer that both listens to and sends messages
CURRENTLY IT ASSUMES THAT A PEER NEVER LEAVES AND TCP CONNECTIONS DON'T DROP
*/

type PoWPeer struct {
	signatureStrategy          signature_strategy.SignatureInterface
	lotteryStrategy            lottery_strategy.LotteryInterface
	publicToSecret             chan map[string]string
	unfinalizedTransactions    chan []blockchain.SignedTransaction
	blockTreeChan              chan blockchain.Blocktree
	unhandledBlocks            chan blockchain.Block
	unhandledMessages          chan blockchain.Message
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
	p.lotteryStrategy = &lottery_strategy.PoW{}
	address, err := network.StringToAddress(IpPort)
	if err != nil {
		panic("Could not parse IpPort: " + err.Error())
	}
	p.network = network.Network{}
	messagesFromNetwork := p.network.StartNetwork(address)

	p.stopMiningSignal = make(chan struct{})

	p.unfinalizedTransactions = make(chan []blockchain.SignedTransaction, 1)
	p.unfinalizedTransactions <- make([]blockchain.SignedTransaction, 0, 100)
	p.publicToSecret = make(chan map[string]string, 1)
	p.publicToSecret <- make(map[string]string)
	p.blockTreeChan = make(chan blockchain.Blocktree, 1)
	newBlockTree, blockTreeCreationWentWell := blockchain.NewBlocktree(blockchain.CreateGenesisBlock())
	if !blockTreeCreationWentWell {
		panic("Could not generate new blocktree")
	}
	p.unhandledBlocks = make(chan blockchain.Block, 20)
	p.hardness = newBlockTree.GetHead().BlockData.Hardness
	p.maximumTransactionsInBlock = constants.BlockSize
	p.unhandledMessages = make(chan blockchain.Message, 50)
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

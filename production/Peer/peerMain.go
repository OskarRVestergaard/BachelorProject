package Peer

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"sync"
)

/*
This is a single Peer that both listens to and sends messages
CURRENTLY IT ASSUMES THAT A PEER NEVER LEAVES AND TCP CONNECTIONS DON'T DROP
*/

type Peer struct {
	signatureStrategy          signature_strategy.SignatureInterface
	lotteryStrategy            lottery_strategy.LotteryInterface
	Ledger                     *models.Ledger
	validMutex                 sync.Mutex
	PublicToSecret             map[string]string
	unfinalizedTransMutex      sync.Mutex
	unfinalizedTransactions    []blockchain.SignedTransaction
	blockTreeChan              chan blockchain.Blocktree
	unhandledBlocks            chan blockchain.Block
	unhandledMessages          chan blockchain.Message
	hardness                   int
	maximumTransactionsInBlock int
	network                    network.Network
}

func (p *Peer) RunPeer(IpPort string) {
	p.signatureStrategy = signature_strategy.ECDSASig{}
	p.lotteryStrategy = &lottery_strategy.PoW{}
	address, err := network.StringToAddress(IpPort)
	if err != nil {
		panic("Could not parse IpPort: " + err.Error())
	}
	p.network = network.Network{}
	messagesFromNetwork := p.network.StartNetwork(address)

	p.Ledger = MakeLedger()

	p.PublicToSecret = make(map[string]string)
	p.blockTreeChan = make(chan blockchain.Blocktree, 1)
	newBlockTree, blockTreeCreationWentWell := blockchain.NewBlocktree(blockchain.CreateGenesisBlock())
	if !blockTreeCreationWentWell {
		panic("Could not generate new blocktree")
	}
	p.unhandledBlocks = make(chan blockchain.Block, 20)
	p.hardness = newBlockTree.GetHead().BlockData.Hardness
	p.maximumTransactionsInBlock = 10
	p.unhandledMessages = make(chan blockchain.Message, 50)
	p.blockTreeChan <- newBlockTree

	go p.blockHandlerLoop()
	go p.messageHandlerLoop(messagesFromNetwork)
}

package Peer

import (
	"encoding/gob"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
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
	IpPort                     string
	ActiveConnections          map[string]models.Void
	Encoders                   map[string]*gob.Encoder
	Ledger                     *models.Ledger
	decoderMutex               sync.Mutex
	acMutex                    sync.Mutex
	encMutex                   sync.Mutex
	floodMutex                 sync.Mutex
	validMutex                 sync.Mutex
	PublicToSecret             map[string]string
	unfinalizedTransMutex      sync.Mutex
	unfinalizedTransactions    []blockchain.SignedTransaction
	blockTreeChan              chan blockchain.Blocktree
	unhandledBlocks            chan blockchain.Block
	unhandledMessages          chan blockchain.Message
	hardness                   int
	maximumTransactionsInBlock int
	MessageLog                 []blockchain.Message
	logMutex                   sync.Mutex
	sentMutex                  sync.Mutex
	receivedMutex              sync.Mutex
	sentCounter                int
	receivedCounter            int
}

func (p *Peer) RunPeer(IpPort string) {
	p.signatureStrategy = signature_strategy.ECDSASig{}
	p.lotteryStrategy = &lottery_strategy.PoW{}
	p.IpPort = IpPort
	p.acMutex.Lock()
	p.ActiveConnections = make(map[string]models.Void)
	p.acMutex.Unlock()
	p.Ledger = MakeLedger()
	p.encMutex.Lock()
	p.Encoders = make(map[string]*gob.Encoder)
	p.encMutex.Unlock()
	p.AddIpPort(IpPort)
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

	go p.startBlockHandler()
	go p.startListener()
	go p.messageHandlerLoop()
}

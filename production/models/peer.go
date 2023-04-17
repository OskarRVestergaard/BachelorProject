package models

import (
	"encoding/gob"
	"fmt"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/models/messages"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"github.com/google/uuid"
	"io"
	"math/big"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

/*
This is a single Peer that both listens to and sends messages
CURRENTLY IT ASSUMES THAT A PEER NEVER LEAVES AND TCP CONNECTIONS DON'T DROP
*/

type Void struct {
}

var member Void

type Message struct {
	MessageType       string
	MessageSender     string
	SignedTransaction messages.SignedTransaction
	MessageBlocks     []blockchain.Block
	PeerMap           map[string]Void
}

type Peer struct {
	signatureStrategy       signature_strategy.SignatureInterface
	lotteryStrategy         lottery_strategy.LotteryInterface
	IpPort                  string
	ActiveConnections       map[string]Void
	Encoders                map[string]*gob.Encoder
	Ledger                  *Ledger
	acMutex                 sync.Mutex
	encMutex                sync.Mutex
	floodMutex              sync.Mutex
	validMutex              sync.Mutex
	PublicToSecret          map[string]string
	unfinalizedTransMutex   sync.Mutex
	unfinalizedTransactions []messages.SignedTransaction
	blockTreeMutex          sync.Mutex
	blockTree               *blockchain.Blocktree
	unhandledBlocks         chan blockchain.Block
	//TODO FinalizedBlockChain (Is actually not needed, just add transaction to ledger, and do cleanup on unfinilized tranactions and blocktree)
}

func (p *Peer) RunPeer(IpPort string) {
	//p.signatureStrategy = signature_strategy.RSASig{}
	p.signatureStrategy = signature_strategy.ECDSASig{}
	p.lotteryStrategy = lottery_strategy.PoW{}
	p.IpPort = IpPort
	p.acMutex.Lock()
	p.ActiveConnections = make(map[string]Void)
	p.acMutex.Unlock()
	p.Ledger = MakeLedger()
	p.encMutex.Lock()
	p.Encoders = make(map[string]*gob.Encoder)
	p.encMutex.Unlock()
	p.AddIpPort(IpPort)
	p.PublicToSecret = make(map[string]string)
	p.blockTreeMutex.Lock()
	p.blockTree = blockchain.NewBlocktree(blockchain.CreateGenesisBlock())
	p.blockTreeMutex.Unlock()
	p.unhandledBlocks = make(chan blockchain.Block, 20)

	time.Sleep(1500 * time.Millisecond)
	go p.startBlockHandler()
	go p.startListener()
}

func (p *Peer) startBlockHandler() {
	for {
		blockToHandle := <-p.unhandledBlocks
		go p.handleBlock(blockToHandle)
	}
}

func (p *Peer) handleBlock(block blockchain.Block) {
	//TODO The check are currently made here, this can hurt performance since some part might be done multiple times for a given block
	hasCorrectSignature := block.HasCorrectSignature(p.signatureStrategy)
	if !hasCorrectSignature {
		return
	}
	//Other checks? //TODO HERE!!!!!!
	p.blockTreeMutex.Lock()
	var t = p.blockTree.AddBlock(block)
	switch t {
	case -2:
		//Block with isGenesis true, not a real block and should be ignored
	case -1:
		//Block is in tree already and can be ignored
	case 0:
		//Parent is not in the tree, try to add later
		//TODO Maybe have another slice that are blocks which are waiting for parents to be added,
		//TODO such that they can be added immediately follow the parents addition to the tree (in case 1)
		p.blockTreeMutex.Unlock()
		time.Sleep(200 * time.Millisecond)
		p.blockTreeMutex.Lock()
		p.unhandledBlocks <- block
	case 1:
		//Block successfully added to the tree
	default:
		p.blockTreeMutex.Unlock()
		panic("addBlockReturnValueNotUnderstood")
	}
	p.blockTreeMutex.Unlock()
}

func (p *Peer) startListener() {
	ln, _ := net.Listen("tcp", p.IpPort)
	for {
		conn, _ := ln.Accept()
		p.AddIpPort(conn.LocalAddr().String())
		go p.Receiver(conn)
	}
}

func (p *Peer) FloodSignedTransaction(from string, to string, amount int) {
	p.floodMutex.Lock()

	t := messages.SignedTransaction{Id: uuid.New(), From: from, To: to, Amount: amount, Signature: big.NewInt(1000000)}

	p.validMutex.Lock()
	msg := Message{MessageType: constants.SignedTransaction, MessageSender: p.IpPort, SignedTransaction: t}

	// TODO: WOW MAGI SOM LAVER SIGNED TRANSACTION TIL EN BESKED DER KAN HASHES BURDE MÅSKE FIXES ORDENTLIGT PÅ ET TIDSPUNKT :D MVH Winther Wonderboy
	hashedMessage := hash_strategy.HashSignedTransactionToByteArrayWowSoCool(msg.SignedTransaction)
	publicKey := msg.SignedTransaction.From
	if secretSigningKey, ok := p.PublicToSecret[from]; ok {
		signatureToAssign := p.signatureStrategy.Sign(hashedMessage, secretSigningKey)
		msg.SignedTransaction.Signature = signatureToAssign
	}
	signature := msg.SignedTransaction.Signature
	if p.signatureStrategy.Verify(publicKey, hashedMessage, signature) {
		p.addTransaction(msg.SignedTransaction)
	} else {
		p.Ledger.Mutex.Lock()
		p.Ledger.UTA++
		p.Ledger.Mutex.Unlock()
	}
	p.validMutex.Unlock()
	p.FloodMessage(msg)
	p.floodMutex.Unlock()

}

func (p *Peer) validTransaction(from string, amount int) bool {
	if amount == 0 {
		println("Invalid SignedTransaction with the amount 0")
		return false
	} else if p.Ledger.Accounts[from] < amount {
		println("Account should hold the transaction amount")
		return false
	}
	return true
}

func (p *Peer) FloodMessage(msg Message) {

	p.acMutex.Lock()
	ac := p.ActiveConnections
	p.acMutex.Unlock()
	for e := range ac {
		if e != p.IpPort {
			err := p.SendMessageTo(e, msg)
			if err != nil {
				println(err.Error())
			}
		}
	}
}

func (p *Peer) SendMessageTo(ipPort string, msg Message) error {
	var enc *gob.Encoder

	p.encMutex.Lock()
	if val, isIn := p.Encoders[ipPort]; isIn {
		enc = val
	} else {
		var conn, err = net.Dial("tcp", ipPort)
		if err != nil {
			p.encMutex.Unlock()
			return err
		}
		enc = gob.NewEncoder(conn)
		p.Encoders[ipPort] = enc
	}

	err := enc.Encode(msg)
	if err != nil {
		p.encMutex.Unlock()
		return err
	}
	p.encMutex.Unlock()
	return nil
}

func (p *Peer) AddIpPort(ipPort string) {
	p.acMutex.Lock()
	p.ActiveConnections[ipPort] = member
	p.acMutex.Unlock()
}

func (p *Peer) Receiver(conn net.Conn) {

	msg := &Message{}
	dec := gob.NewDecoder(conn)
	for {
		err := dec.Decode(msg)
		savedMsg := *msg
		if err == io.EOF {
			err2 := conn.Close()
			print(err2.Error())
			return
		}
		if err != nil {
			println(err.Error())
			err2 := conn.Close()
			print(err2.Error())
			return
		}
		handled := savedMsg
		p.handleMessage(handled)
	}
}

func (p *Peer) handleMessage(msg Message) {
	msgType := (msg).MessageType

	switch msgType {
	case constants.SignedTransaction:
		p.validMutex.Lock()

		// TODO: WOW MAGI SOM LAVER SIGNED TRANSACTION TIL EN BESKED DER KAN HASHES BURDE MÅSKE FIXES ORDENTLIGT PÅ ET TIDSPUNKT :D MVH Winther Wonderboy
		hashedMessage := hash_strategy.HashSignedTransactionToByteArrayWowSoCool(msg.SignedTransaction)
		publicKey := msg.SignedTransaction.From
		signature := msg.SignedTransaction.Signature
		if p.signatureStrategy.Verify(publicKey, hashedMessage, signature) {
			p.addTransaction((msg).SignedTransaction)
		} else {
			p.Ledger.Mutex.Lock()
			p.Ledger.UTA++
			p.Ledger.Mutex.Unlock()
		}
		p.validMutex.Unlock()
	case constants.JoinMessage:
		p.AddIpPort((msg).MessageSender)
	case constants.GetPeersMessage:
		p.acMutex.Lock()
		ac := p.ActiveConnections
		p.acMutex.Unlock()
		err := p.SendMessageTo((msg).MessageSender, Message{MessageType: constants.PeerMapDelivery, MessageSender: p.IpPort, PeerMap: ac})
		if err != nil {
			println(err.Error())
		}
	case constants.PeerMapDelivery:
		for e := range (msg).PeerMap {
			p.AddIpPort(e)
		}
		p.FloodMessage(Message{MessageType: constants.JoinMessage, MessageSender: p.IpPort})
	case constants.BlockDelivery:
		for _, block := range msg.MessageBlocks {
			p.unhandledBlocks <- block
		}
	default:
		println(p.IpPort + ": received a UNKNOWN message type from: " + (msg).MessageSender)
	}
}

func (p *Peer) PrintLedger() {
	p.Ledger.Mutex.Lock()
	l := p.Ledger.Accounts
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	println()
	print("Ledger of: " + p.IpPort + ": ")
	for _, value := range keys {
		print("[" + strconv.Itoa(l[value]) + "]")
	}
	println(" with the amount of SignedTransactions: " + strconv.Itoa(p.Ledger.TA) + " and deniedTransactions: " + strconv.Itoa(p.Ledger.UTA))
	p.Ledger.Mutex.Unlock()
}

func MakeLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]int)
	return ledger
}

func (p *Peer) addTransaction(t messages.SignedTransaction) {
	p.unfinalizedTransMutex.Lock()
	p.unfinalizedTransactions = append(p.unfinalizedTransactions, t)
	p.unfinalizedTransMutex.Unlock()
}

func (p *Peer) UpdateLedger(transactions []*messages.SignedTransaction) {

	p.Ledger.Mutex.Lock()
	for _, trans := range transactions {
		p.Ledger.TA = p.Ledger.TA + 1
		p.Ledger.Accounts[trans.From] -= trans.Amount
		p.Ledger.Accounts[trans.To] += trans.Amount
	}
	p.Ledger.Mutex.Unlock()

}

func (p *Peer) Connect(ip string, port int) {
	ipPort := ip + ":" + strconv.Itoa(port)
	err := p.SendMessageTo(ipPort, Message{MessageType: constants.GetPeersMessage, MessageSender: p.IpPort})

	if err != nil {
		println(err.Error())
		p.printNewNetworkStarted()
	}
}

// TODO Rename or remove
func (p *Peer) printNewNetworkStarted() {
	println("Network started")
	println("********************************************************************")
	println("Host IP: " + p.IpPort)
	println("********************************************************************")
}

func (p *Peer) CreateAccount() string {
	secretKey, publicKey := p.signatureStrategy.KeyGen()
	p.PublicToSecret[publicKey] = secretKey
	return publicKey
}

// CreateBalanceOnLedger for testing only
// TODO Remove
func (p *Peer) CreateBalanceOnLedger(pk string, amount int) {

	p.Ledger.Mutex.Lock()
	p.Ledger.Accounts[pk] += amount
	p.Ledger.Mutex.Unlock()
}

// TODO FIX LATER
func (p *Peer) Mine() {
	fmt.Println("We're there. All I can see are turtle tracks. Whaddaya say we give Bowser the old Brooklyn one-two?")
	var hasPotentialWinner bool
	for k := range p.PublicToSecret {
		hasPotentialWinner, _ = p.lotteryStrategy.Mine(k, "PrevHash")
	}

	if hasPotentialWinner {
		// do block stuff IDK
	}
}

func (p *Peer) SendFakeBlockWithTransactions() {
	var publicKey = utils.GetSomeKey(p.PublicToSecret)
	var secretKey = p.PublicToSecret[publicKey]
	var headBlock = p.blockTree.GetHead()
	var headBlockHash = headBlock.HashOfBlock()
	var blockWithCurrentlyUnhandledTransactions = blockchain.Block{
		IsGenesis: false,
		Vk:        publicKey,
		Slot:      1,
		Draw:      "+4 cards",
		BlockData: blockchain.BlockData{
			Transactions: p.unfinalizedTransactions, //TODO Should only add not already added transactions (ones not in the chain)
		},
		ParentHash: headBlockHash,
		Signature:  big.Int{},
	}
	blockWithCurrentlyUnhandledTransactions.CalculateSignature(p.signatureStrategy, secretKey)
	var msg = Message{
		MessageType:   constants.BlockDelivery,
		MessageSender: p.IpPort,
		MessageBlocks: []blockchain.Block{blockWithCurrentlyUnhandledTransactions},
	}
	go p.FloodMessage(msg)
	for _, block := range msg.MessageBlocks {
		p.unhandledBlocks <- block
	}
}

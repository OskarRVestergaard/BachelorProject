package models

import (
	"encoding/gob"
	"fmt"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/models/messages"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
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

var debugging bool

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
	SignatureStrategy       signature_strategy.SignatureInterface
	LotteryStrategy         lottery_strategy.LotteryInterface
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
	UnfinalizedTransactions []*messages.SignedTransaction
	blockMutex              sync.Mutex
	blockTree               *blockchain.Blocktree
	unhandledBlocks         []blockchain.Block
	//TODO FinilizedBlockChain (Is actually not needed, just add transaction to ledger, and do cleanup on unfinilized tranactions and blocktree)
}

func (p *Peer) RunPeer(IpPort string) {
	//p.SignatureStrategy = signature_strategy.RSASig{}
	p.SignatureStrategy = signature_strategy.ECDSASig{}
	p.LotteryStrategy = lottery_strategy.PoW{}
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
	p.blockMutex.Lock()
	p.blockMutex.Unlock()

	p.blockTree = blockchain.NewBlocktree(blockchain.CreateGenesisBlock())

	time.Sleep(2500 * time.Millisecond)
	go p.StartListener()
}

func (p *Peer) StartListener() {
	ln, _ := net.Listen("tcp", p.IpPort)
	for {
		conn, _ := ln.Accept()
		p.AddIpPort(conn.LocalAddr().String())
		go p.Receiver(conn)
	}
}

func (p *Peer) FloodSignedTransaction(from string, to string, amount int) {
	debug(p.IpPort + " called doSignedTransaction")

	p.floodMutex.Lock()

	t := messages.SignedTransaction{Id: GenerateId(), From: from, To: to, Amount: amount, Signature: big.NewInt(1000000)}

	p.validMutex.Lock()
	msg := Message{MessageType: constants.SignedTransaction, MessageSender: p.IpPort, SignedTransaction: t}

	// TODO: WOW MAGI SOM LAVER SIGNED TRANSACTION TIL EN BESKED DER KAN HASHES BURDE MÅSKE FIXES ORDENTLIGT PÅ ET TIDSPUNKT :D MVH Winther Wonderboy
	hashedMessage := hash_strategy.HashSignedTransactionToByteArrayWowSoCool(msg.SignedTransaction)
	publicKey := msg.SignedTransaction.From
	if val, ok := p.PublicToSecret[from]; ok {
		signatureToAssign := p.SignatureStrategy.Sign(hashedMessage, val)
		msg.SignedTransaction.Signature = signatureToAssign
	}
	signature := msg.SignedTransaction.Signature
	if p.SignatureStrategy.Verify(publicKey, hashedMessage, signature) {
		p.UpdateUncontrolledTransactions(msg.SignedTransaction)
	} else {
		p.Ledger.Mutex.Lock()
		p.Ledger.UTA++
		p.Ledger.Mutex.Unlock()
	}
	p.validMutex.Unlock()
	p.FloodMessage(msg)
	p.floodMutex.Unlock()

}

func GenerateId() uuid.UUID {
	Id := uuid.New()
	return Id
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
	debug(p.IpPort + " called sendMessageTo")
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
	debug(p.IpPort + " called addIpPort and adding: " + ipPort)
	p.acMutex.Lock()
	p.ActiveConnections[ipPort] = member
	p.acMutex.Unlock()
}

func (p *Peer) Receiver(conn net.Conn) {
	debug(p.IpPort + " called receiver")

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
		if p.SignatureStrategy.Verify(publicKey, hashedMessage, signature) {
			p.UpdateUncontrolledTransactions((msg).SignedTransaction)
		} else {
			p.Ledger.Mutex.Lock()
			p.Ledger.UTA++
			p.Ledger.Mutex.Unlock()
		}
		p.validMutex.Unlock()
	case constants.JoinMessage: //"joinMessage":
		debug(p.IpPort + ": received a join message from: " + (msg).MessageSender)
		p.AddIpPort((msg).MessageSender)
	case constants.GetPeersMessage: //"getPeersMessage":
		debug(p.IpPort + ": received a getPeers message from: " + (msg).MessageSender)
		p.acMutex.Lock()
		ac := p.ActiveConnections
		p.acMutex.Unlock()
		err := p.SendMessageTo((msg).MessageSender, Message{MessageType: constants.PeerMapDelivery, MessageSender: p.IpPort, PeerMap: ac})
		if err != nil {
			println(err.Error())
		}
	case constants.PeerMapDelivery: //"peerMapDelivery":
		debug(p.IpPort + ": received a peerMapDelivery message from: " + (msg).MessageSender)
		for e := range (msg).PeerMap {
			p.AddIpPort(e)
			debug("added: " + e)
		}
		p.FloodMessage(Message{MessageType: constants.JoinMessage, MessageSender: p.IpPort})
	case constants.BlockDelivery:
		print("BLOCK DELIVERY FROM: " + (msg).MessageSender)
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

func (p *Peer) UpdateUncontrolledTransactions(t messages.SignedTransaction) {
	debug(p.IpPort + " called updateLedger")
	p.unfinalizedTransMutex.Lock()
	p.UnfinalizedTransactions = append(p.UnfinalizedTransactions, &t)
	p.unfinalizedTransMutex.Unlock()
}

func (p *Peer) UpdateLedger(transactions []*messages.SignedTransaction) {
	debug(p.IpPort + " called updateLedger")

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
		p.startNewNetwork()
	}
}

// TODO Rename or remove
func (p *Peer) startNewNetwork() {
	println("Network started")
	println("********************************************************************")
	println("Host IP: " + p.IpPort)
	println("********************************************************************")
}

func (p *Peer) PrintActiveCons() {
	println("Peer: " + p.IpPort + " has the following connections: ")
	p.acMutex.Lock()
	ac := p.ActiveConnections
	p.acMutex.Unlock()
	for e := range ac {
		println(e)
	}
}

func debug(msg string) {
	if debugging {
		println(msg)
	}
}

func (p *Peer) CreateAccount() string {

	secretKey, publicKey := p.SignatureStrategy.KeyGen()

	p.PublicToSecret[publicKey] = secretKey

	return publicKey
}

// CreateBalanceOnLedger for testing only
// TODO Remove
func (p *Peer) CreateBalanceOnLedger(pk string, amount int) {

	debug(p.IpPort + " called updateLedger")

	p.Ledger.Mutex.Lock()
	p.Ledger.Accounts[pk] += amount

	p.Ledger.Mutex.Unlock()
}

func (p *Peer) FloodBlocks(slotNumber int) {

}

// TODO FIX LATER
func (p *Peer) Mine() {
	fmt.Println("We're there. All I can see are turtle tracks. Whaddaya say we give Bowser the old Brooklyn one-two?")
	var hasPotentialWinner bool
	for k := range p.PublicToSecret {
		hasPotentialWinner, _ = p.LotteryStrategy.Mine(k, "PrevHash")
	}

	if hasPotentialWinner {
		// do block stuff IDK
	}
}

//func MakeBlock(transactions []*structs.SignedTransaction, prevHash string) Block {
//	//TODO add maximum blockSize
//	var b Block
//	//b.slotNumber = slot
//	b.PreviousHash = prevHash
//	//b.TransactionsLog = transactions
//	b.Transactions = transactions
//	b.Hash = calculateHash(b.PreviousHash, b.Transactions)
//	slot += 1
//	return b
//
//}

//TODO: empty unctontrolled list when sending block + receiver should remove doublicates. ledger should be updated correct

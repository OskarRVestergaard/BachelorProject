package peer

import (
	"encoding/gob"
	"example.com/packages/block"
	"example.com/packages/models"
	"example.com/packages/signature_strategy"
	"example.com/packages/utils"
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

type Ledger struct {
	Accounts map[string]int
	mutex    sync.Mutex
	TA       int
	Uta      int
}

type Void struct {
}

var member Void

type Message struct {
	MessageType       string
	MessageSender     string
	SignedTransaction models.SignedTransaction
	PeerMap           map[string]Void
}

type Peer struct {
	SignatureStrategy signature_strategy.SignatureInterface
	IpPort            string
	ActiveConnections map[string]Void
	Encoders          map[string]*gob.Encoder
	Ledger            *Ledger
	acMutex           sync.Mutex
	encMutex          sync.Mutex
	floodMutex        sync.Mutex
	validMutex        sync.Mutex
	PublicToSecret    map[string]string
	VisibleBlockchain map[int]*block.Block
}

func (p *Peer) RunPeer(IpPort string) {
	p.SignatureStrategy = signature_strategy.RSASig{}
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
	p.VisibleBlockchain = make(map[int]*block.Block)
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
	t := models.SignedTransaction{From: from, To: to, Amount: amount, Signature: big.NewInt(1000000)}

	p.validMutex.Lock()
	//msg := Message{"SignedTransaction", p.IpPort, t, map[string]Void{}}
	msg := Message{utils.SignedTransaction, p.IpPort, t, map[string]Void{}}

	if val, ok := p.PublicToSecret[from]; ok {

		signature := p.SignatureStrategy.Sign(msg.SignedTransaction, val)
		msg.SignedTransaction.Signature = signature
	}

	if p.SignatureStrategy.Verify(msg.SignedTransaction) {
		p.UpdateLedger(msg.SignedTransaction)
	} else {
		p.Ledger.mutex.Lock()
		p.Ledger.Uta++
		p.Ledger.mutex.Unlock()
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
	case utils.SignedTransaction:
		p.validMutex.Lock()
		if p.SignatureStrategy.Verify(msg.SignedTransaction) {
			p.UpdateLedger((msg).SignedTransaction)
		} else {
			p.Ledger.mutex.Lock()
			p.Ledger.Uta++
			p.Ledger.mutex.Unlock()
		}
		p.validMutex.Unlock()
	case utils.JoinMessage: //"joinMessage":
		debug(p.IpPort + ": received a join message from: " + (msg).MessageSender)
		p.AddIpPort((msg).MessageSender)
	case utils.GetPeersMessage: //"getPeersMessage":
		debug(p.IpPort + ": received a getPeers message from: " + (msg).MessageSender)
		p.acMutex.Lock()
		ac := p.ActiveConnections
		p.acMutex.Unlock()
		err := p.SendMessageTo((msg).MessageSender, Message{MessageType: utils.PeerMapDelivery, MessageSender: p.IpPort, SignedTransaction: models.SignedTransaction{Signature: big.NewInt(0)}, PeerMap: ac})
		if err != nil {
			println(err.Error())
		}
	case utils.PeerMapDelivery: //"peerMapDelivery":
		debug(p.IpPort + ": received a peerMapDelivery message from: " + (msg).MessageSender)
		//TODO FIX: THIS ASSUMES THAT THE ONLY TIME THIS MESSAGE IS RECEIVED IS WHEN CONNECTING TO A NETWORK (since we flood joinMessage)
		for e := range (msg).PeerMap {
			p.AddIpPort(e)
			debug("added: " + e)
		}
		p.FloodMessage(Message{MessageType: utils.JoinMessage, MessageSender: p.IpPort, SignedTransaction: models.SignedTransaction{Signature: big.NewInt(0)}, PeerMap: map[string]Void{}})
	default:
		println(p.IpPort + ": received a UNKNOWN message type from: " + (msg).MessageSender)
	}
}

func (p *Peer) PrintLedger() {
	p.Ledger.mutex.Lock()
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
	println(" with the amount of SignedTransactions: " + strconv.Itoa(p.Ledger.TA) + " and deniedTransactions: " + strconv.Itoa(p.Ledger.Uta))
	p.Ledger.mutex.Unlock()
}

func MakeLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]int)
	return ledger
}

func (p *Peer) UpdateLedger(t models.SignedTransaction) {
	debug(p.IpPort + " called updateLedger")

	p.Ledger.mutex.Lock()

	p.Ledger.TA = p.Ledger.TA + 1
	p.Ledger.Accounts[t.From] -= t.Amount
	p.Ledger.Accounts[t.To] += t.Amount

	p.Ledger.mutex.Unlock()
}

func (p *Peer) Connect(ip string, port int) {
	ipPort := ip + ":" + strconv.Itoa(port)
	err := p.SendMessageTo(ipPort, Message{MessageType: utils.GetPeersMessage, MessageSender: p.IpPort, SignedTransaction: models.SignedTransaction{Signature: big.NewInt(0)}, PeerMap: map[string]Void{}})

	if err != nil {
		println(err.Error())
		p.startNewNetwork()
	}
}

func (p *Peer) startNewNetwork() {
	p.VisibleBlockchain = makeGenesisBlockchain() //"asd" // = makeGenesisBlockchain()
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

//----------------------

func (p *Peer) CreateAccount() string {

	n, d, e := p.SignatureStrategy.KeyGen()

	publicKey := n.String() + ";" + e.String() + ";"
	secretKey := n.String() + ";" + d.String() + ";"

	p.PublicToSecret[publicKey] = secretKey

	return publicKey
}
func makeGenesisBlockchain() map[int]*block.Block {
	genesisBlock := &block.Block{
		SlotNumber:      0,
		Hash:            "GenesisBlock",
		PreviousHash:    "GenesisBlock",
		TransactionsLog: nil,
	}
	var blockChain = make(map[int]*block.Block)
	blockChain[genesisBlock.SlotNumber] = genesisBlock
	//blockChain = append(blockChain, (genesisBlock))
	return blockChain
}

// CreateBalanceOnLedger for testing only
func (p *Peer) CreateBalanceOnLedger(pk string, amount int) {

	debug(p.IpPort + " called updateLedger")

	p.Ledger.mutex.Lock()
	p.Ledger.Accounts[pk] += amount

	p.Ledger.mutex.Unlock()
}

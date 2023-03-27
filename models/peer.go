package models

import (
	"encoding/gob"
	"example.com/packages/strategies/hash_strategy"
	lottery_strategy2 "example.com/packages/strategies/lottery_strategy"
	signature_strategy2 "example.com/packages/strategies/signature_strategy"
	"example.com/packages/structs"
	"example.com/packages/utils/Constants"
	"fmt"
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
	SignedTransaction structs.SignedTransaction
	MessageBlocks     []Block
	PeerMap           map[string]Void
}

type Peer struct {
	SignatureStrategy        signature_strategy2.SignatureInterface
	LotteryStrategy          lottery_strategy2.LotteryInterface
	IpPort                   string
	ActiveConnections        map[string]Void
	Encoders                 map[string]*gob.Encoder
	Ledger                   *Ledger
	acMutex                  sync.Mutex
	encMutex                 sync.Mutex
	floodMutex               sync.Mutex
	validMutex               sync.Mutex
	PublicToSecret           map[string]string
	uncontrolledTransMutex   sync.Mutex
	UncontrolledTransactions []*structs.SignedTransaction
	GenesisBlock             []*Block
	blockMutex               sync.Mutex
	Blocks                   []Block
}

func (p *Peer) RunPeer(IpPort string) {
	//p.SignatureStrategy = signature_strategy.RSASig{}
	p.SignatureStrategy = signature_strategy2.ECDSASig{}
	p.LotteryStrategy = lottery_strategy2.PoW{}
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

	t := structs.SignedTransaction{Id: GenerateId(), From: from, To: to, Amount: amount, Signature: big.NewInt(1000000)}

	p.validMutex.Lock()
	msg := Message{MessageType: Constants.SignedTransaction, MessageSender: p.IpPort, SignedTransaction: t}

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
		p.Ledger.mutex.Lock()
		p.Ledger.Uta++
		p.Ledger.mutex.Unlock()
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
	case Constants.SignedTransaction:
		p.validMutex.Lock()

		// TODO: WOW MAGI SOM LAVER SIGNED TRANSACTION TIL EN BESKED DER KAN HASHES BURDE MÅSKE FIXES ORDENTLIGT PÅ ET TIDSPUNKT :D MVH Winther Wonderboy
		hashedMessage := hash_strategy.HashSignedTransactionToByteArrayWowSoCool(msg.SignedTransaction)
		publicKey := msg.SignedTransaction.From
		signature := msg.SignedTransaction.Signature
		if p.SignatureStrategy.Verify(publicKey, hashedMessage, signature) {
			p.UpdateUncontrolledTransactions((msg).SignedTransaction)
		} else {
			p.Ledger.mutex.Lock()
			p.Ledger.Uta++
			p.Ledger.mutex.Unlock()
		}
		p.validMutex.Unlock()
	case Constants.JoinMessage: //"joinMessage":
		debug(p.IpPort + ": received a join message from: " + (msg).MessageSender)
		p.AddIpPort((msg).MessageSender)
	case Constants.GetPeersMessage: //"getPeersMessage":
		debug(p.IpPort + ": received a getPeers message from: " + (msg).MessageSender)
		p.acMutex.Lock()
		ac := p.ActiveConnections
		p.acMutex.Unlock()
		//Also sends genesis block
		err := p.SendMessageTo((msg).MessageSender, Message{MessageType: Constants.PeerMapDelivery, MessageSender: p.IpPort, SignedTransaction: structs.SignedTransaction{Signature: big.NewInt(0)}, MessageBlocks: []Block{*makeGenesisBlock()}, PeerMap: ac})
		if err != nil {
			println(err.Error())
		}
	case Constants.PeerMapDelivery: //"peerMapDelivery":
		debug(p.IpPort + ": received a peerMapDelivery message from: " + (msg).MessageSender)
		//TODO FIX: THIS ASSUMES THAT THE ONLY TIME THIS MESSAGE IS RECEIVED IS WHEN CONNECTING TO A NETWORK (since we flood joinMessage)
		//TODO should not just append blocks on genesis
		for e := range (msg).PeerMap {
			p.AddIpPort(e)
			debug("added: " + e)
		}

		for e := range (msg).MessageBlocks {
			p.UpdateBlock(msg.MessageBlocks[e])
		}

		p.FloodMessage(Message{MessageType: Constants.JoinMessage, MessageSender: p.IpPort, SignedTransaction: structs.SignedTransaction{Signature: big.NewInt(0)}})
	case Constants.Block:
		for e := range (msg).MessageBlocks {
			p.UpdateLedger(msg.MessageBlocks[e].Transactions)
		}
		p.GenesisBlock = append(p.GenesisBlock, &msg.MessageBlocks[0])
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

func (p *Peer) UpdateUncontrolledTransactions(t structs.SignedTransaction) {
	debug(p.IpPort + " called updateLedger")

	p.uncontrolledTransMutex.Lock()
	p.UncontrolledTransactions = append(p.UncontrolledTransactions, &t)
	p.uncontrolledTransMutex.Unlock()
}

func (p *Peer) UpdateLedger(transactions []*structs.SignedTransaction) {
	debug(p.IpPort + " called updateLedger")

	p.Ledger.mutex.Lock()
	for _, trans := range transactions {
		p.Ledger.TA = p.Ledger.TA + 1
		p.Ledger.Accounts[trans.From] -= trans.Amount
		p.Ledger.Accounts[trans.To] += trans.Amount
	}
	p.Ledger.mutex.Unlock()

}

func (p *Peer) Connect(ip string, port int) {
	ipPort := ip + ":" + strconv.Itoa(port)
	err := p.SendMessageTo(ipPort, Message{MessageType: Constants.GetPeersMessage, MessageSender: p.IpPort})

	if err != nil {
		println(err.Error())
		p.startNewNetwork()
	}
}

func (p *Peer) startNewNetwork() {
	p.blockMutex.Lock()
	p.GenesisBlock = append(p.GenesisBlock, makeGenesisBlock())
	p.blockMutex.Unlock()
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

func makeGenesisBlock() *Block {

	genesisBlock := &Block{
		SlotNumber:   0,
		Hash:         "GenesisBlock",
		PreviousHash: "GenesisBlock",
		Transactions: nil,
	}
	return genesisBlock
}

// CreateBalanceOnLedger for testing only
func (p *Peer) CreateBalanceOnLedger(pk string, amount int) {

	debug(p.IpPort + " called updateLedger")

	p.Ledger.mutex.Lock()
	p.Ledger.Accounts[pk] += amount

	p.Ledger.mutex.Unlock()
}

func (p *Peer) UpdateBlock(b Block) {
	//debug(p.IpPort + " called addIpPort and adding: " + &b.Hash)
	p.blockMutex.Lock()
	p.GenesisBlock = append(p.GenesisBlock, &b)
	p.blockMutex.Unlock()
}

func (p *Peer) FloodBlocks(slotNumber int) {
	//t := structs.SignedTransaction{Id: GenerateId(), From: from, To: to, Amount: amount, Signature: big.NewInt(1000000)}
	var block1 Block
	block1.SlotNumber = slotNumber
	block1.PreviousHash = "Nothing"
	for b := 0; b < Constants.BlockSize; b++ {
		if len(p.UncontrolledTransactions) >= Constants.BlockSize {
			block1.Transactions = append(block1.Transactions, p.UncontrolledTransactions[b])
		}
	}
	block1.Hash = ConvertToString(block1.Transactions)
	var blockArray = []Block{block1}

	var t = Message{
		MessageType: Constants.Block,
		//MessageSender:     "",
		//SignedTransaction: structs.SignedTransaction{},
		MessageBlocks: blockArray,
		//PeerMap:           nil,
	}
	p.FloodMessage(t)
	//p.SendMessageTo((msg).MessageSender, Message{MessageType: utils.PeerMapDelivery, MessageSender: p.IpPort, SignedTransaction: structs.SignedTransaction{Signature: big.NewInt(0)}, MessageBlocks: []block.Block{*makeGenesisBlock()}, PeerMap: ac})
	//block.Block{
	//	SlotNumber:   slotNumber,
	//	Hash:         "",
	//	PreviousHash: "",
	//	Transactions: nil,

}

func (p *Peer) Mine() {
	fmt.Println("We're there. All I can see are turtle tracks. Whaddaya say we give Bowser the old Brooklyn one-two?")
	var hasPotentialWinner bool
	for k := range p.PublicToSecret {
		hasPotentialWinner, _ = p.LotteryStrategy.Mine(k, p.GenesisBlock[len(p.GenesisBlock)-1].PreviousHash)
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

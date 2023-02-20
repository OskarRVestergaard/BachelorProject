package peer

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"io"
	"math/big"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var debugging bool

/*
This is a single Peer that both listens to and sends messages
CURRENTLY IT ASSUMES THAT PEER NEVER LEAVE AND TCP CONNECTIONS DONT DROP
*/

type Ledger struct {
	Accounts map[string]int
	mutex    sync.Mutex
	TA       int
	Uta      int
}

type Peer struct {
	IpPort            string
	ActiveConnections map[string]Void
	Encoders          map[string]*gob.Encoder
	Ledger            *Ledger
	acMutex           sync.Mutex
	encMutex          sync.Mutex
	floodMutex        sync.Mutex
	validMutex        sync.Mutex
	PublicToSecret    map[string]string
}

type Void struct {
}

var member Void

type SignedTransaction struct {
	From      string
	To        string
	Amount    int
	Signature *big.Int
}

type Message struct {
	MessageType       string
	MessageSender     string
	SignedTransaction SignedTransaction
	PeerMap           map[string]Void
}

func (p *Peer) RunPeer(IpPort string) {
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
	if amount == 0 {
		println("Invalid SignedTransaction with the amount 0")
	} else {
		p.floodMutex.Lock()
		t := SignedTransaction{from, to, amount, big.NewInt(1000000)}

		p.validMutex.Lock()
		msg := Message{"SignedTransaction", p.IpPort, t, map[string]Void{}}

		if val, ok := p.PublicToSecret[from]; ok {
			signature := CreateSigniture(msg.SignedTransaction, val)
			msg.SignedTransaction.Signature = signature
		}

		if ValidateSignature(msg.SignedTransaction) {
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
			conn.Close()
			return
		}
		if err != nil {
			println(err.Error())
			conn.Close()
			return
		}
		handled := savedMsg
		p.handleMessage(handled)
	}
}

func (p *Peer) handleMessage(msg Message) {
	msgType := (msg).MessageType

	switch msgType {
	case "SignedTransaction":
		p.validMutex.Lock()
		if ValidateSignature(msg.SignedTransaction) {
			p.UpdateLedger((msg).SignedTransaction)
		} else {
			p.Ledger.mutex.Lock()
			p.Ledger.Uta++
			p.Ledger.mutex.Unlock()
		}
		p.validMutex.Unlock()
	case "joinMessage":
		debug(p.IpPort + ": received a join message from: " + (msg).MessageSender)
		p.AddIpPort((msg).MessageSender)
	case "getPeersMessage":
		debug(p.IpPort + ": received a getPeers message from: " + (msg).MessageSender)
		p.acMutex.Lock()
		ac := p.ActiveConnections
		p.acMutex.Unlock()
		err := p.SendMessageTo((msg).MessageSender, Message{"peerMapDelivery", p.IpPort, SignedTransaction{"", "", 0, big.NewInt(0)}, ac})
		if err != nil {
			println(err.Error())
		}
	case "peerMapDelivery":
		debug(p.IpPort + ": received a peerMapDelivery message from: " + (msg).MessageSender)
		//TODO FIX: THIS ASSUMES THAT THE ONLY TIME THIS MESSAGE IS RECEIVED IS WHEN CONNECTING TO A NETWORK (since we flood joinMessage)
		for e := range (msg).PeerMap {
			p.AddIpPort(e)
			debug("added: " + e)
		}
		p.FloodMessage(Message{"joinMessage", p.IpPort, SignedTransaction{"", "", 0, big.NewInt(0)}, map[string]Void{}})
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

func (p *Peer) UpdateLedger(t SignedTransaction) {
	debug(p.IpPort + " called updateLedger")

	p.Ledger.mutex.Lock()

	p.Ledger.TA = p.Ledger.TA + 1
	p.Ledger.Accounts[t.From] -= t.Amount
	p.Ledger.Accounts[t.To] += t.Amount

	p.Ledger.mutex.Unlock()
}

func (p *Peer) Connect(ip string, port int) {
	ipPort := ip + ":" + strconv.Itoa(port)
	err := p.SendMessageTo(ipPort, Message{"getPeersMessage", p.IpPort, SignedTransaction{"", "", 0, big.NewInt(0)}, map[string]Void{}})

	if err != nil {
		println(err.Error())
		p.startNewNetwork()
	}
}

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

//----------------------

func KeyGen(k int) (*big.Int, *big.Int, *big.Int) {
	e := big.NewInt(3)

	b := k / 2

	if k%2 != 0 {
		b += 1
	}

	p, _ := rand.Prime(rand.Reader, b)
	q, _ := rand.Prime(rand.Reader, b)
	n := big.NewInt(0)

	n = n.Mul(p, q)

	l := big.NewInt(0)
	l2 := big.NewInt(0)

	q_minus_one := big.NewInt(0)
	q_minus_one = q_minus_one.Sub(q, big.NewInt(1))
	p_minus_one := big.NewInt(0)
	p_minus_one = p_minus_one.Sub(p, big.NewInt(1))

	for {
		if l.GCD(nil, nil, e, q_minus_one).Cmp(big.NewInt(1)) == 0 {
			break
		}

		q, _ = rand.Prime(rand.Reader, b)
		q_minus_one = q_minus_one.Sub(q, big.NewInt(1))

	}

	for {
		if l2.GCD(nil, nil, e, p_minus_one).Cmp(big.NewInt(1)) == 0 {
			break
		}

		p, _ = rand.Prime(rand.Reader, b)
		p_minus_one = p_minus_one.Sub(p, big.NewInt(1))
	}

	pq_minus_ones := big.NewInt(0)
	pq_minus_ones = pq_minus_ones.Mul(p_minus_one, q_minus_one)

	n = n.Mul(p, q)

	d := big.NewInt(0)
	d = d.Exp(e, big.NewInt(-1), pq_minus_ones)

	return n, d, e

}

func (p *Peer) CreateAccount() string {
	n, d, e := KeyGen(2048)

	publicKey := n.String() + ";" + e.String() + ";"
	secretKey := n.String() + ";" + d.String() + ";"

	p.PublicToSecret[publicKey] = secretKey

	return publicKey
}

func CreateSigniture(transaction SignedTransaction, secretKey string) *big.Int {
	n, d := splitKey(secretKey)

	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
	t = strings.Replace(t, ";", "", -1)

	hashed := hashMe(t)
	sign := Decrypt(hashed, n, d)
	return sign
}

func splitKey(key string) (*big.Int, *big.Int) {
	splitkey := strings.Split(key, ";")
	n_string := splitkey[0]
	de_string := splitkey[1]

	n := big.NewInt(0)
	de := big.NewInt(0)

	n, _ = n.SetString(n_string, 10)
	de, _ = de.SetString(de_string, 10)

	return n, de
}

func ValidateSignature(transaction SignedTransaction) bool {
	signature := transaction.Signature

	pk := transaction.From
	n, e := splitKey(pk)
	unsigned := Encrypt(signature, n, e)

	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
	t = strings.Replace(t, ";", "", -1)

	//append signature to message
	hashed := hashMe(t)

	return (hashed.Cmp(unsigned) == 0)
}

func Encrypt(msg *big.Int, n *big.Int, e *big.Int) *big.Int {
	res := big.NewInt(0)
	res = res.Exp(msg, e, n)
	return res
}

func Decrypt(cipher *big.Int, n *big.Int, d *big.Int) *big.Int {
	res := big.NewInt(0)
	res = res.Exp(cipher, d, n)
	return res
}

func hashMe(msg string) *big.Int {
	h := sha256.New()
	h.Write([]byte(msg))

	hm := new(big.Int).SetBytes(h.Sum(nil))

	return hm
}

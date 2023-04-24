package Peer

import (
	"encoding/gob"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"github.com/google/uuid"
	"io"
	"net"
	"strconv"
)

// TODO Rename or remove
func (p *Peer) printNewNetworkStarted() {
	println("Network started")
	println("********************************************************************")
	println("Host IP: " + p.IpPort)
	println("********************************************************************")
}

func (p *Peer) Connect(ip string, port int) {
	ipPort := ip + ":" + strconv.Itoa(port)
	print(p.IpPort + " Connecting to " + ipPort + "\n")
	err := p.SendMessageTo(ipPort, blockchain.Message{MessageType: constants.GetPeersMessage, MessageSender: p.IpPort})

	if err != nil {
		println(err.Error())
		p.printNewNetworkStarted()
	}
}

func (p *Peer) FloodMessage(msg blockchain.Message) {
	p.logMessage(msg)
	p.acMutex.Lock()
	ac := p.ActiveConnections
	for e := range ac {
		if e != p.IpPort {
			err := p.SendMessageTo(e, msg)
			if err != nil {
				println(err.Error())
			}
		}
	}
	p.acMutex.Unlock()
}

func (p *Peer) SendMessageTo(ipPort string, msg blockchain.Message) error {
	var enc *gob.Encoder
	p.sentMutex.Lock()
	p.sentCounter = p.sentCounter + 1
	p.sentMutex.Unlock()
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

	err := enc.Encode(utils.MakeDeepCopyOfMessage(msg))
	if err != nil {
		p.encMutex.Unlock()
		return err
	}
	p.encMutex.Unlock()
	return nil
}

func (p *Peer) AddIpPort(ipPort string) {
	p.acMutex.Lock()
	p.ActiveConnections[ipPort] = models.Void{}
	p.acMutex.Unlock()
}

func (p *Peer) Receiver(conn net.Conn) { //Also a loop

	dec := gob.NewDecoder(conn)
	for {
		newMsg := &blockchain.Message{}
		p.decoderMutex.Lock()
		err := dec.Decode(newMsg)
		p.receivedMutex.Lock()
		p.receivedCounter = p.receivedCounter + 1
		p.receivedMutex.Unlock()
		if err == io.EOF {
			err2 := conn.Close()
			print(err2.Error())
			p.encMutex.Unlock()
			return
		}
		if err != nil {
			println(err.Error())
			err2 := conn.Close()
			print(err2.Error())
			p.encMutex.Unlock()
			return
		}
		msgCopy := utils.MakeDeepCopyOfMessage(*newMsg)
		p.logMessage(msgCopy)
		p.unhandledMessages <- msgCopy
		p.decoderMutex.Unlock()
	}
}

func (p *Peer) messageHandlerLoop() {
	for {
		msg := <-p.unhandledMessages
		p.handleMessage(msg)
	}
}

func (p *Peer) handleMessage(msg blockchain.Message) {
	msgType := (msg).MessageType

	switch msgType {
	case constants.SignedTransaction:
		p.validMutex.Lock()
		if utils.TransactionHasCorrectSignature(p.signatureStrategy, msg.SignedTransaction) {
			deepCopyOfTransaction := utils.MakeDeepCopyOfTransaction(msg.SignedTransaction)
			p.addTransaction(deepCopyOfTransaction)
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
		err := p.SendMessageTo((msg).MessageSender, blockchain.Message{MessageType: constants.PeerMapDelivery, MessageSender: p.IpPort, PeerMap: ac})
		p.acMutex.Unlock()
		if err != nil {
			println(err.Error())
		}
	case constants.PeerMapDelivery:
		for e := range (msg).PeerMap {
			p.AddIpPort(e)
		}
		p.FloodMessage(blockchain.Message{MessageType: constants.JoinMessage, MessageSender: p.IpPort})
	case constants.BlockDelivery:
		for _, block := range msg.MessageBlocks {
			p.unhandledBlocks <- block
		}
	default:
		println(p.IpPort + ": received a UNKNOWN message type from: " + (msg).MessageSender)
	}
}

func (p *Peer) startListener() { //Go NewTCPListenerLoop
	ln, err2 := net.Listen("tcp", p.IpPort)
	if err2 != nil {
		panic(err2.Error())
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("Error happened for listener: " + err.Error())
		}
		p.AddIpPort(conn.LocalAddr().String())
		go p.Receiver(conn)
	}
}

func (p *Peer) FloodSignedTransaction(from string, to string, amount int) {
	p.floodMutex.Lock()
	trans := blockchain.SignedTransaction{Id: uuid.New(), From: from, To: to, Amount: amount, Signature: nil}

	p.validMutex.Lock()

	secretSigningKey, foundSecretKey := p.PublicToSecret[from]
	if foundSecretKey {
		trans.SignTransaction(p.signatureStrategy, secretSigningKey)
	}
	msg := blockchain.Message{MessageType: constants.SignedTransaction, MessageSender: p.IpPort, SignedTransaction: trans}
	p.addTransaction(trans)
	p.validMutex.Unlock()
	p.FloodMessage(msg)
	p.floodMutex.Unlock()
}

func (p *Peer) logMessage(msg blockchain.Message) {
	p.logMutex.Lock()
	p.MessageLog = append(p.MessageLog, msg)
	p.logMutex.Unlock()
}

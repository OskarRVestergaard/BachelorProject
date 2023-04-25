package Peer

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"github.com/google/uuid"
)

func (p *Peer) GetAddress() network.Address {
	return p.network.GetAddress()
}

func (p *Peer) Connect(ip string, port int) {
	addr := network.Address{
		Ip:   ip,
		Port: port,
	}
	ownIpPort := p.network.GetAddress().ToString()
	print(ownIpPort + " Connecting to " + addr.ToString() + "\n")
	err := p.network.SendMessageTo(blockchain.Message{MessageType: constants.GetPeersMessage, MessageSender: ownIpPort}, addr)

	if err != nil {
		panic(err.Error())
	}
}

func (p *Peer) messageHandlerLoop(incomingMessages chan blockchain.Message) {
	for {
		msg := <-incomingMessages
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
		}
		p.validMutex.Unlock()
	case constants.JoinMessage:

	case constants.GetPeersMessage:

	case constants.PeerMapDelivery:

	case constants.BlockDelivery:
		for _, block := range msg.MessageBlocks {
			p.unhandledBlocks <- block
		}
	default:
		println(p.network.GetAddress().ToString() + ": received a UNKNOWN message type ( " + msg.MessageType + " ) from: " + msg.MessageSender)
	}
}

func (p *Peer) FloodSignedTransaction(from string, to string, amount int) {
	trans := blockchain.SignedTransaction{Id: uuid.New(), From: from, To: to, Amount: amount, Signature: nil}

	p.validMutex.Lock()

	secretSigningKey, foundSecretKey := p.PublicToSecret[from]
	if foundSecretKey {
		trans.SignTransaction(p.signatureStrategy, secretSigningKey)
	}
	ipPort := p.network.GetAddress().ToString()
	msg := blockchain.Message{MessageType: constants.SignedTransaction, MessageSender: ipPort, SignedTransaction: trans}
	p.addTransaction(trans)
	p.validMutex.Unlock()
	p.network.FloodMessageToAllKnown(msg)
}

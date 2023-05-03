package PowPeer

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
	err := p.network.SendMessageTo(blockchain.Message{MessageType: constants.JoinMessage, MessageSender: ownIpPort}, addr)

	if err != nil {
		panic(err.Error())
	}
}

func (p *Peer) FloodSignedTransaction(from string, to string, amount int) {
	secretSigningKey, foundSecretKey := p.getSecretKey(from)
	if !foundSecretKey {
		return
	}
	trans := blockchain.SignedTransaction{Id: uuid.New(), From: from, To: to, Amount: amount, Signature: nil}
	trans.SignTransaction(p.signatureStrategy, secretSigningKey)
	ipPort := p.network.GetAddress().ToString()
	msg := blockchain.Message{MessageType: constants.SignedTransaction, MessageSender: ipPort, SignedTransaction: trans}
	p.addTransaction(trans)
	p.network.FloodMessageToAllKnown(msg)
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
		if utils.TransactionHasCorrectSignature(p.signatureStrategy, msg.SignedTransaction) {
			p.addTransaction(utils.MakeDeepCopyOfTransaction(msg.SignedTransaction))
		}
	case constants.JoinMessage:

	case constants.BlockDelivery:
		for _, block := range msg.MessageBlocks {
			p.unhandledBlocks <- block
		}
	default:
		println(p.network.GetAddress().ToString() + ": received a UNKNOWN message type ( " + msg.MessageType + " ) from: " + msg.MessageSender)
	}
}

/*
getSecretKey

returns the secret key associated with a given public key and return a boolean indicating whether the key is known
*/
func (p *Peer) getSecretKey(pk string) (secretKey string, isKnownKey bool) {
	publicToSecret := <-p.publicToSecret
	secretSigningKey, foundSecretKey := publicToSecret[pk]
	p.publicToSecret <- publicToSecret
	if !foundSecretKey {
		return "", false
	}
	return secretSigningKey, true
}

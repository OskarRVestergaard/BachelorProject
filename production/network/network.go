package network

import (
	"encoding/gob"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"io"
	"net"
)

/*
Public
*/

type Network struct {
	ownAddress       Address
	encoders         chan map[Address]*gob.Encoder
	incomingMessages chan blockchain.Message
	outgoingMessages chan outgoingMessage
}

func (network *Network) GetAddress() Address {
	return network.ownAddress
}

func (network *Network) StartNetwork(add Address) (receivedMessages chan blockchain.Message) {
	network.ownAddress = add
	connectionChannel := make(chan net.Conn, 4)
	network.incomingMessages = make(chan blockchain.Message, 20)
	network.outgoingMessages = make(chan outgoingMessage, 20)
	network.encoders = make(chan map[Address]*gob.Encoder, 1)
	encodersMap := make(map[Address]*gob.Encoder)
	network.encoders <- encodersMap
	go network.listenerLoop(connectionChannel)
	go network.connDelegationLoop(connectionChannel)
	go network.senderLoop()
	return network.incomingMessages
}

/*
SendMessageTo

might be blocking if the network is busy sending messages
*/
func (network *Network) SendMessageTo(message blockchain.Message, address Address) error {
	if !network.isKnownAddress(address) {
		var conn, err = net.Dial("tcp", address.ToString())
		if err != nil {
			return err
		}

		enc := gob.NewEncoder(conn)
		encoders := <-network.encoders
		encoders[address] = enc
		network.encoders <- encoders
	}
	msg := outgoingMessage{
		message:     message,
		destination: address,
	}
	network.outgoingMessages <- msg
	return nil
}

func (network *Network) FloodMessageToAllKnown(message blockchain.Message) {
	encoders := <-network.encoders
	for address, _ := range encoders {
		msg := outgoingMessage{
			message:     message,
			destination: address,
		}
		network.outgoingMessages <- msg
	}
	network.encoders <- encoders
}

/*
Private
*/

type outgoingMessage struct {
	message     blockchain.Message
	destination Address
}

func (network *Network) isKnownAddress(address Address) bool {
	encoders := <-network.encoders
	_, found := encoders[address]
	network.encoders <- encoders
	return found
}

func (network *Network) senderLoop() {
	for {
		msg := <-network.outgoingMessages
		go network.handleSendMessage(msg)
	}
}

func (network *Network) handleSendMessage(message outgoingMessage) {
	encoders := <-network.encoders
	encoder, hasEncoder := encoders[message.destination]
	if !hasEncoder {
		panic("Tried to send message to unknown party")
	}
	err := encoder.Encode(utils.MakeDeepCopyOfMessage(message.message))
	if err != nil {
		panic("Something went wrong during message encoding")
	}
	network.encoders <- encoders
}

func (network *Network) listenerLoop(connections chan net.Conn) {
	ln, err := net.Listen("tcp", network.ownAddress.ToString())
	if err != nil {
		panic(err.Error())
	}
	for {
		conn, err2 := ln.Accept()
		if err2 != nil {
			panic("Error happened for listener: " + err2.Error())
		}
		//Maybe check if this is an already establish connection
		connections <- conn
	}
}

func (network *Network) connDelegationLoop(connections chan net.Conn) {
	for {
		conn := <-connections
		go network.connectionReceiverLoop(conn)
	}
}

func (network *Network) connectionReceiverLoop(conn net.Conn) {
	remoteAddress, connErr := ConnToRemoteAddress(conn)
	if network.isKnownAddress(remoteAddress) {
		return
	}
	if connErr != nil {
		panic("Address given by connection failed parsing")
	}
	dec := gob.NewDecoder(conn)

	for {
		newMsg := &blockchain.Message{}
		err := dec.Decode(newMsg)
		if err == io.EOF {
			closingError := conn.Close()
			println(closingError.Error())
			return
		}
		if err != nil {
			println(err.Error())
			closingError := conn.Close()
			println(closingError.Error())
			return
		}
		network.incomingMessages <- utils.MakeDeepCopyOfMessage(*newMsg)
	}
}

package network

import (
	"encoding/gob"
	"github.com/OskarRVestergaard/BachelorProject/production/Message"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"io"
	"net"
)

/*
Public
*/

type Network struct {
	ownAddress       Address
	encoders         chan map[Address]*gob.Encoder
	incomingMessages chan Message.Message
	outgoingMessages chan outgoingMessage
}

func (network *Network) GetAddress() Address {
	return network.ownAddress
}

func (network *Network) EstablishConnectionTo(address Address) error {
	encoders := <-network.encoders
	if !network.isKnownAddress(address, encoders) {
		var conn, err = net.Dial("tcp", address.ToString())
		if err != nil {
			return err
		}

		enc := gob.NewEncoder(conn)
		encoders[address] = enc
		network.encoders <- encoders
	}
	return nil
}

func (network *Network) StartNetwork(address Address) (receivedMessages chan Message.Message) {
	network.ownAddress = address
	connectionChannel := make(chan net.Conn, 4)
	network.incomingMessages = make(chan Message.Message, 50)
	network.outgoingMessages = make(chan outgoingMessage, 100)
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
func (network *Network) SendMessageTo(message Message.Message, address Address) error {
	encoders := <-network.encoders
	if !network.isKnownAddress(address, encoders) {
		var conn, err = net.Dial("tcp", address.ToString())
		if err != nil {
			return err
		}

		enc := gob.NewEncoder(conn)
		encoders[address] = enc
		network.encoders <- encoders
	}
	msg := outgoingMessage{
		message:     Message.MakeDeepCopyOfMessage(message),
		destination: address,
	}
	network.outgoingMessages <- msg
	return nil
}

func (network *Network) FloodMessageToAllKnown(message Message.Message) {
	if message.MessageType == constants.BlockDelivery {
		print("test")
	}
	encoders := <-network.encoders
	for address, _ := range encoders {
		msg := outgoingMessage{
			message:     Message.MakeDeepCopyOfMessage(message),
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
	message     Message.Message
	destination Address
}

func (network *Network) isKnownAddress(address Address, encoders map[Address]*gob.Encoder) bool {
	_, found := encoders[address]
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
	err := encoder.Encode(Message.MakeDeepCopyOfMessage(message.message))
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
	dec := gob.NewDecoder(conn)

	for {
		newMsg := &Message.Message{}
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
		network.incomingMessages <- Message.MakeDeepCopyOfMessage(*newMsg)
	}
}

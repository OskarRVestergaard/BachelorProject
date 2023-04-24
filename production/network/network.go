package network

import (
	"encoding/gob"
	"errors"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"io"
	"net"
)

type Network struct {
	ownAddress       Address
	encoders         chan map[Address]*gob.Encoder
	incomingMessages chan blockchain.Message
	outgoingMessages chan outgoingMessage
}

type outgoingMessage struct {
	message     blockchain.Message
	destination Address
}

/*
Public methods
*/

func (network *Network) GetAddress() Address {
	return network.ownAddress
}

func (network *Network) StartNetwork(add Address) (incomingMessages chan blockchain.Message) {
	network.ownAddress = add
	connectionChannel := make(chan net.Conn, 4)
	network.incomingMessages = make(chan blockchain.Message, 20)
	network.outgoingMessages = make(chan outgoingMessage, 20)
	network.encoders = make(chan map[Address]*gob.Encoder, 1)
	go network.listenerLoop(connectionChannel)
	go network.connDelegationLoop(connectionChannel)
	go network.senderLoop()
	//Maybe more here? A message sender loop probably
	return incomingMessages
}

func (network *Network) SendMessageTo(message blockchain.Message, address Address) {
	msg := outgoingMessage{
		message:     message,
		destination: address,
	}
	network.outgoingMessages <- msg
}

func (network *Network) FloodMessageToAllKnown(message blockchain.Message) {
	encoders := <-network.encoders
	for address, _ := range encoders {
		network.SendMessageTo(message, address)
	}
	network.encoders <- encoders
}

/*
Private methods
*/

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

func (network *Network) establishConnectionWith(address Address) error {
	//This check is to avoid sending this unnecessary message, it must still be checked again in handle connection
	if network.isKnownAddress(address) {
		return errors.New("the network already has a connection with that address")
	}
	var conn, err = net.Dial("tcp", address.ToString())
	if err != nil {
		return err
	}
	network.handleConnection(conn)
	return nil
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
		//CHECK THAT THIS IS NOT ALREADY A KNOWN CONNECTION
		connections <- conn
	}
}

func (network *Network) connDelegationLoop(connections chan net.Conn) {
	for {
		conn := <-connections
		go network.handleConnection(conn)
	}
}

func (network *Network) handleConnection(conn net.Conn) {
	remoteAddress, connErr := ConnToRemoteAddress(conn)
	if network.isKnownAddress(remoteAddress) {
		return
	}
	if connErr != nil {
		panic("Address given by connection failed parsing")
	}
	enc := gob.NewEncoder(conn)
	encoders := <-network.encoders
	encoders[remoteAddress] = enc
	network.encoders <- encoders
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
		//Todo Change sender to actual sender such that no peer can cheat by assigning another sender
		network.incomingMessages <- utils.MakeDeepCopyOfMessage(*newMsg)
	}
}

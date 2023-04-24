package network

import (
	"encoding/gob"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"io"
	"net"
)

type Network struct {
	ownAddress       Address
	incomingMessages chan blockchain.Message
}

func (network *Network) GetAddress() Address {
	return network.ownAddress
}

func (network *Network) StartNetwork(add Address) (incomingMessages chan blockchain.Message) {
	network.ownAddress = add
	connectionChannel := make(chan net.Conn, 4)
	network.incomingMessages = make(chan blockchain.Message, 20)
	go network.listenerLoop(connectionChannel)
	go network.connDelegationLoop(connectionChannel)

	//Maybe more here?
	return incomingMessages
}

func (network *Network) establishConnectionWith(ip string, port int) bool {
	return false //TODO
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

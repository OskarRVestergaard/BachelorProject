package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBasicNetwork(t *testing.T) {
	//Setup
	net1 := network.Network{}
	addr1 := network.Address{
		Ip:   "127.0.0.1",
		Port: 65065,
	}

	net2 := network.Network{}
	addr2 := network.Address{
		Ip:   "127.0.0.1",
		Port: 65066,
	}

	dummyMessage := blockchain.Message{
		MessageType:       "1",
		MessageSender:     "2",
		SignedTransaction: blockchain.SignedTransaction{},
		MessageBlocks:     nil,
		PeerMap:           nil,
	}

	//Actions
	incomingMessages1 := net1.StartNetwork(addr1)
	incomingMessages2 := net2.StartNetwork(addr2)
	time.Sleep(2000 * time.Millisecond)
	err := net1.SendMessageTo(dummyMessage, addr2)
	if err != nil {
		panic(err.Error())
	}
	err2 := net2.SendMessageTo(dummyMessage, addr1)
	if err2 != nil {
		panic(err2.Error())
	}

	msg1 := <-incomingMessages1
	msg2 := <-incomingMessages2

	//Asserts
	assert.True(t, msg1.MessageType == "1")
	assert.True(t, msg2.MessageType == "1")
}

package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/sha256"
	"github.com/google/uuid"
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

func TestNetworkWithNilElements(t *testing.T) {
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

	dummyMessage1 := blockchain.Message{
		MessageType:       "1",
		MessageSender:     "2",
		SignedTransaction: blockchain.SignedTransaction{},
		MessageBlocks:     nil,
		PeerMap:           nil,
	}

	dummyMessage2 := blockchain.Message{
		MessageType:       "",
		MessageSender:     "",
		SignedTransaction: blockchain.SignedTransaction{},
		MessageBlocks:     nil,
		PeerMap:           nil,
	}

	//Actions
	incomingMessages1 := net1.StartNetwork(addr1)
	incomingMessages2 := net2.StartNetwork(addr2)
	time.Sleep(2000 * time.Millisecond)
	err := net1.SendMessageTo(dummyMessage1, addr2)
	if err != nil {
		panic(err.Error())
	}
	err2 := net2.SendMessageTo(dummyMessage1, addr1)
	if err2 != nil {
		panic(err2.Error())
	}

	msg1 := <-incomingMessages1
	msg2 := <-incomingMessages2

	err3 := net1.SendMessageTo(dummyMessage2, addr2)
	if err3 != nil {
		panic(err3.Error())
	}
	err4 := net2.SendMessageTo(dummyMessage2, addr1)
	if err4 != nil {
		panic(err4.Error())
	}

	msg3 := <-incomingMessages1
	msg4 := <-incomingMessages2

	//Asserts
	assert.True(t, msg1.MessageType == "1")
	assert.True(t, msg2.MessageType == "1")
	assert.True(t, msg3.MessageType == "")
	assert.True(t, msg4.MessageType == "")
}

func TestBiggerNetworkWithFlooding(t *testing.T) {
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

	net3 := network.Network{}
	addr3 := network.Address{
		Ip:   "127.0.0.1",
		Port: 65067,
	}

	net4 := network.Network{}
	addr4 := network.Address{
		Ip:   "127.0.0.1",
		Port: 65068,
	}

	randomId := uuid.New()
	msgBlocks := []blockchain.Block{{
		IsGenesis: false,
		Vk:        "123134",
		Slot:      4,
		Draw: lottery_strategy.WinningLotteryParams{
			Vk:         "9999",
			ParentHash: sha256.HashByteArray([]byte{byte(32), byte(66)}),
			Counter:    43,
		},
		BlockData:  blockchain.BlockData{},
		ParentHash: sha256.HashValue{},
		Signature:  []byte{byte(32), byte(2)}},
	}
	peerMap := make(map[string]models.Void)
	peerMap["1"] = struct{}{}
	peerMap["5"] = struct{}{}

	dummyMessage1 := blockchain.Message{
		MessageType:       "1",
		MessageSender:     "2",
		SignedTransaction: blockchain.SignedTransaction{},
		MessageBlocks:     nil,
		PeerMap:           nil,
	}

	dummyMessage2 := blockchain.Message{
		MessageType:   "9",
		MessageSender: "5",
		SignedTransaction: blockchain.SignedTransaction{
			Id:        randomId,
			From:      "431",
			To:        "21",
			Amount:    43,
			Signature: []byte{byte(4), byte(6)},
		},
		MessageBlocks: msgBlocks,
		PeerMap:       peerMap,
	}

	dummyMessage3 := blockchain.Message{
		MessageType:       "",
		MessageSender:     "",
		SignedTransaction: blockchain.SignedTransaction{},
		MessageBlocks:     nil,
		PeerMap:           nil,
	}

	//Actions
	incomingMessages1 := net1.StartNetwork(addr1)
	incomingMessages2 := net2.StartNetwork(addr2)
	incomingMessages3 := net3.StartNetwork(addr3)
	incomingMessages4 := net4.StartNetwork(addr4)
	time.Sleep(2000 * time.Millisecond)

	_ = net1.EstablishConnectionTo(addr2)
	_ = net1.EstablishConnectionTo(addr3)
	_ = net1.EstablishConnectionTo(addr4)
	_ = net2.EstablishConnectionTo(addr1)
	_ = net2.EstablishConnectionTo(addr3)
	_ = net2.EstablishConnectionTo(addr4)
	_ = net3.EstablishConnectionTo(addr2)
	_ = net3.EstablishConnectionTo(addr1)
	_ = net3.EstablishConnectionTo(addr4)
	_ = net4.EstablishConnectionTo(addr2)
	_ = net4.EstablishConnectionTo(addr3)
	_ = net4.EstablishConnectionTo(addr1)

	net1.FloodMessageToAllKnown(dummyMessage1)
	net4.FloodMessageToAllKnown(dummyMessage2)
	net3.FloodMessageToAllKnown(dummyMessage3)
	net1.FloodMessageToAllKnown(dummyMessage2)
	net2.FloodMessageToAllKnown(dummyMessage1)
	net1.FloodMessageToAllKnown(dummyMessage3)
	net4.FloodMessageToAllKnown(dummyMessage3)

	time.Sleep(150000 * time.Millisecond)
	//Asserts
	assert.Equal(t, len(incomingMessages1), len(incomingMessages2))
	assert.Equal(t, len(incomingMessages2), len(incomingMessages3))
	assert.Equal(t, len(incomingMessages3), len(incomingMessages4))
}

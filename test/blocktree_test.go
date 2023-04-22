package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/test/networkservice"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlockDelivery(t *testing.T) {
	noOfPeers := 2
	noOfMsgs := 3
	noOfNames := 2
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)      //send msg
	debugDraw := lottery_strategy.WinningLotteryParams{
		Vk:         "DEBUG",
		ParentHash: nil,
		Counter:    0,
	}

	time.Sleep(1000 * time.Millisecond)

	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)

	time.Sleep(1000 * time.Millisecond)

	listOfPeers[0].SendBlockWithTransactions(8, debugDraw)

	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)

	time.Sleep(1000 * time.Millisecond)

	listOfPeers[0].SendBlockWithTransactions(6, debugDraw)

	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)

	time.Sleep(1000 * time.Millisecond)

	listOfPeers[0].SendBlockWithTransactions(11, debugDraw)
	time.Sleep(1000 * time.Millisecond)
	time.Sleep(1000 * time.Millisecond)
	assert.True(t, true)
}

func TestPOWNetwork2Peers(t *testing.T) {
	noOfPeers := 2
	noOfMsgs := 3
	noOfNames := 2
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)      //send msg
	for _, peer := range listOfPeers {
		peer.StartMining()
	}
	for i := 0; i < 10; i++ {
		networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(4000 * time.Millisecond)
	}
	time.Sleep(20000 * time.Millisecond)
	//TODO Add way to stop mining, such that the network will stabilize
	tree1 := listOfPeers[0].GetBlockTree()
	tree2 := listOfPeers[1].GetBlockTree()
	assert.Equal(t, tree1, tree2)
	print("test")
}

func TestPOWNetwork4Peers(t *testing.T) {
	noOfPeers := 4
	noOfMsgs := 2
	noOfNames := 4
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)      //send msg
	for _, peer := range listOfPeers {
		peer.StartMining()
	}
	for i := 0; i < 4; i++ {
		networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(10000 * time.Millisecond)
	}
	time.Sleep(100000 * time.Millisecond)

	res1 := listOfPeers[0].GetBlockTree()
	res2 := listOfPeers[1].GetBlockTree()
	res3 := listOfPeers[2].GetBlockTree()
	res4 := listOfPeers[3].GetBlockTree()
	print(res1.GetHead().Slot)
	print(res2.GetHead().Slot)
	print(res3.GetHead().Slot)
	print(res4.GetHead().Slot)
	assert.True(t, true)
}

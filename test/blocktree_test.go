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

func TestPOWNetwork(t *testing.T) {
	noOfPeers := 2
	noOfMsgs := 3
	noOfNames := 2
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)      //send msg
	for _, peer := range listOfPeers {
		peer.StartMining()
	}
	//for i := 0; i < 10; i++ {
	//	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
	//	time.Sleep(5000 * time.Millisecond)
	//}
	time.Sleep(10000 * time.Millisecond)

	assert.True(t, true)
}

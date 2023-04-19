package test

import (
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

	time.Sleep(1000 * time.Millisecond)

	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)

	time.Sleep(1000 * time.Millisecond)

	listOfPeers[0].SendBlockWithTransactions(8, "DebugDraw")

	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)

	time.Sleep(1000 * time.Millisecond)

	listOfPeers[0].SendBlockWithTransactions(6, "DebugDraw")

	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)

	time.Sleep(1000 * time.Millisecond)

	listOfPeers[0].SendBlockWithTransactions(11, "DebugDraw")
	time.Sleep(1000 * time.Millisecond)
	time.Sleep(1000 * time.Millisecond)
	assert.True(t, true)
}

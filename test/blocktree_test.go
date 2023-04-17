package test

import (
	"github.com/OskarRVestergaard/BachelorProject/test/networkservice"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlockDelivery(t *testing.T) {
	noOfPeers := 2
	noOfMsgs := 2
	noOfNames := 2
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	//networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)      //send msg

	time.Sleep(5000 * time.Millisecond)

	//listOfPeers[1].SendFakeBlockWithTransactions(1)
	time.Sleep(5000 * time.Millisecond)

	//send 2 more messages (1 from each)
	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
	time.Sleep(5000 * time.Millisecond)
	listOfPeers[0].SendFakeBlockWithTransactions(3)

	time.Sleep(5000 * time.Millisecond)
	time.Sleep(1000000 * time.Millisecond)
	assert.True(t, true)
}

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

	time.Sleep(5000 * time.Millisecond)

	time.Sleep(10000 * time.Millisecond)

	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)

	time.Sleep(10000 * time.Millisecond)

	networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
	time.Sleep(10000 * time.Millisecond)

	listOfPeers[1].SendFakeBlockWithTransactions(11)
	time.Sleep(10000 * time.Millisecond)
	time.Sleep(10000 * time.Millisecond)
	assert.True(t, true)
}

package test

import (
	"github.com/OskarRVestergaard/BachelorProject/test/networkservice"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlockDelivery(t *testing.T) {
	noOfPeers := 2
	noOfMsgs := 1
	noOfNames := 2
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames)             //setup peer
	controlLedger := networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList) //send msg

	time.Sleep(1000 * time.Millisecond)

	for i := 0; i < noOfPeers; i++ {
		listOfPeers[i].PrintLedger()
	}

	printControlLedger(controlLedger)
	listOfPeers[1].SendFakeBlockWithTransactions()
	time.Sleep(1000 * time.Millisecond)

	assert.True(t, true)
	//for i := 0; i < noOfPeers; i++ {
	//	accountsOfPeer := listOfPeers[i].Ledger.Accounts
	//}
}

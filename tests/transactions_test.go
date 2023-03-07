package test

import (
	"example.com/packages/service"
	"testing"
	"time"
)

func TestTransactionsAppearInList(t *testing.T) {

	noOfPeers := 2
	noOfMsgs := 3
	noOfNames := 2
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames)             //setup peer
	controlLedger := service.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList) //send msg
	print(controlLedger)
	time.Sleep(1000 * time.Millisecond)

	time.Sleep(1000 * time.Millisecond)

	//p1Genesis := listOfPeers[0].GenesisBlock[0].SlotNumber
	p1Genesis := len(listOfPeers[0].UncontrolledTransactions)
	print(p1Genesis)
	//p2Genesis := listOfPeers[1].GenesisBlock[0].SlotNumber
	//assert.Equal(t, 0, p1Genesis, "genesisblock should have slotnumber 0")
	//assert.Equal(t, 0, p2Genesis, "genesisblock should have slotnumber 0")

}

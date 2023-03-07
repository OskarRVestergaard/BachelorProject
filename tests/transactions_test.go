package test

import (
	"example.com/packages/service"
	"github.com/stretchr/testify/assert"
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

	p1TransactionLog := listOfPeers[0].UncontrolledTransactions
	p2TransactionLog := listOfPeers[1].UncontrolledTransactions

	assert.NotEmpty(t, p1TransactionLog)
	assert.NotEmpty(t, p2TransactionLog)

}

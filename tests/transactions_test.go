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
func TestTransactionsAppearAcrossnetwork(t *testing.T) {
	//Act
	noOfPeers := 2
	noOfNames := 2
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	pk0 := pkList[0]
	pk1 := pkList[1]

	//Action
	listOfPeers[1].FloodSignedTransaction(pk1, pk0, 50)
	time.Sleep(200 * time.Millisecond)
	// Assert
	p0Uncontrolledtransactionssize := len(listOfPeers[0].UncontrolledTransactions)
	p1Uncontrolledtransactionssize := len(listOfPeers[1].UncontrolledTransactions)
	assert.Equal(t, p0Uncontrolledtransactionssize, p1Uncontrolledtransactionssize)

}

func TestFloodBlockOnNetworkWithTransactions(t *testing.T) {
	noOfPeers := 2
	noOfNames := 2
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	pk0 := pkList[0]
	pk1 := pkList[1]

	//Action
	listOfPeers[1].FloodSignedTransaction(pk1, pk0, 50)
	time.Sleep(200 * time.Millisecond)
	//Action
	listOfPeers[1].FloodBlocks(50)
	time.Sleep(200 * time.Millisecond)
	print("DET VIRKER", pkList)
}

package test

import (
	"fmt"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/networkservice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactionsAppearInList(t *testing.T) {

	noOfPeers := 2
	noOfMsgs := 3
	noOfNames := 2
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames)             //setup peer
	controlLedger := networkservice.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList) //send msg
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
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
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
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	pk0 := pkList[0]
	pk1 := pkList[1]

	//Action
	listOfPeers[1].FloodSignedTransaction(pk1, pk0, 50)
	time.Sleep(200 * time.Millisecond)
	//Action
	listOfPeers[1].FloodBlocks(50)
	time.Sleep(250 * time.Millisecond)
	//Assert
	assert.Equal(t, len(listOfPeers[0].GenesisBlock), 2)
}

func TestPoW(t *testing.T) {
	noOfPeers := 2
	noOfNames := 2
	//listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	_, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	pk0 := pkList[0]
	//pk1 := pkList[1]
	//listOfPeers[0].
	b, h := lottery_strategy.PoW{}.Mine(pk0, "asdasd")
	fmt.Println(b)
	fmt.Println(h)

}

func TestPeerMinePoW(t *testing.T) {
	noOfPeers := 2
	noOfNames := 2
	listOfPeers, _ := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	listOfPeers[0].Mine()
}

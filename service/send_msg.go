package service

import (
	"example.com/packages/ledger"
	"example.com/packages/peer"
	"math/rand"
)

func SendMsgs(noOfMsgs int, noOfPeers int, listOfPeers []*peer.Peer, pkList []string) *ledger.Ledger {
	noOfNames := len(pkList)
	controlLedger := new(ledger.Ledger)
	controlLedger.TA = 0
	controlLedger.Accounts = make(map[string]int)
	for j := 0; j < noOfMsgs; j++ {
		for i := 0; i < noOfPeers; i++ {

			p1 := pkList[rand.Intn(noOfNames/noOfPeers)*noOfPeers+i]
			p2 := pkList[rand.Intn(noOfNames)]
			value := rand.Intn(100) + 1
			controlLedger.UpdateLedger(p1, p2, value)

			go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
		}
	}
	return controlLedger

}

func FloodTransactionOnNetwork(noOfMsgs int, noOfPeers int, listOfPeers []*peer.Peer, pkList []string) {

}

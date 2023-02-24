package test

import (
	"example.com/packages/ledger"
	"example.com/packages/service"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"
)

func TestSignedAllValid(t *testing.T) {
	noOfPeers := 10
	noOfMsgs := 5000
	noOfNames := 10
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames)             //setup peer
	controlLedger := service.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList) //send msg

	time.Sleep(30000 * time.Millisecond)

	for i := 0; i < noOfPeers; i++ {
		listOfPeers[i].PrintLedger()
	}

	printControlLedger(controlLedger)

	for i := 0; i < noOfPeers; i++ {
		accountsOfPeer := listOfPeers[i].Ledger.Accounts
		assert.True(t, reflect.DeepEqual(accountsOfPeer, controlLedger.Accounts))
	}
}

func TestSignedOneNotValid(t *testing.T) {
	noOfPeers := 5
	noOfMsgs := 1
	noOfNames := 5
	AccountBalance := 1000
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < noOfNames; i++ {
		pkList[i] = listOfPeers[i%noOfPeers].CreateAccount()
		listOfPeers[i].CreateBalanceOnLedger(pkList[i], AccountBalance)
	}

	controlLedger := new(ledger.Ledger)
	controlLedger.TA = 0
	controlLedger.Accounts = make(map[string]int)
	rand.Seed(time.Now().Unix())

	for j := 1; j < noOfMsgs; j++ {
		for i := 0; i < noOfPeers; i++ {
			p1 := pkList[rand.Intn(noOfNames/noOfPeers)*noOfPeers+i]
			p2 := pkList[rand.Intn(noOfNames)]
			value := rand.Intn(100) + 1
			controlLedger.UpdateLedger(p1, p2, value)
			go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
		}
	}

	p1 := pkList[rand.Intn(noOfNames/noOfPeers)*noOfPeers+1]
	p2 := pkList[rand.Intn(noOfNames)]
	value := rand.Intn(100) + 1
	controlLedger.UpdateLedger(p1, p2, value)
	go listOfPeers[0].FloodSignedTransaction(p1, p2, value)

	for i := 1; i < noOfPeers; i++ {
		p1 = pkList[rand.Intn(noOfNames/noOfPeers)*noOfPeers+i]
		p2 = pkList[rand.Intn(noOfNames)]
		value = rand.Intn(100) + 1
		controlLedger.UpdateLedger(p1, p2, value)
		go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
	}

	time.Sleep(250 * time.Millisecond)
	for i := 0; i < noOfPeers; i++ {
		listOfPeers[i].PrintLedger()
	}

	printControlLedger(controlLedger)

	for i := 0; i < noOfPeers; i++ {
		signedTransactionsOfPeer := listOfPeers[i].Ledger.TA

		assert.Equal(t, signedTransactionsOfPeer, noOfPeers*noOfMsgs-1, "One msg was not signed but still validated")

	}

}

func TestSignedAllRandom(t *testing.T) {

	noOfPeers := 5
	noOfMsgs := 15
	noOfNames := 5
	AccountBalance := 100
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < noOfNames; i++ {
		pkList[i] = listOfPeers[i%noOfPeers].CreateAccount()
		listOfPeers[i].CreateBalanceOnLedger(pkList[i], AccountBalance)
	}

	controlLedger := new(ledger.Ledger)
	controlLedger.TA = 0
	controlLedger.Accounts = make(map[string]int)
	rand.Seed(time.Now().Unix())

	for j := 0; j < noOfMsgs; j++ {
		for i := 0; i < noOfPeers; i++ {
			p1 := pkList[rand.Intn(noOfNames)]
			p2 := pkList[rand.Intn(noOfNames)]
			value := rand.Intn(100) + 1
			controlLedger.UpdateLedger(p1, p2, value)
			go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
		}
	}

	time.Sleep(250 * time.Millisecond)
	for i := 0; i < noOfPeers; i++ {
		listOfPeers[i].PrintLedger()
	}

	printControlLedger(controlLedger)

	for i := 1; i < noOfPeers; i++ {
		accountsOfPeer := listOfPeers[i].Ledger.Accounts
		accountsOfPrevPeer := listOfPeers[i-1].Ledger.Accounts
		assert.True(t, reflect.DeepEqual(accountsOfPeer, accountsOfPrevPeer))

	}
}
func TestNoTransactions(t *testing.T) {

	noOfPeers := 5
	noOfNames := 5
	AccountBalance := 100
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < noOfNames; i++ {
		pkList[i] = listOfPeers[i%noOfPeers].CreateAccount()
		listOfPeers[i].CreateBalanceOnLedger(pkList[i], AccountBalance)
	}

	controlLedger := new(ledger.Ledger)
	controlLedger.TA = 0
	controlLedger.Accounts = make(map[string]int)
	rand.Seed(time.Now().Unix())

	time.Sleep(250 * time.Millisecond)
	for i := 0; i < noOfPeers; i++ {
		listOfPeers[i].PrintLedger()
	}

	printControlLedger(controlLedger)

	for i := 1; i < noOfPeers; i++ {
		accountsOfPeer := listOfPeers[i].Ledger.Accounts
		accountsOfPrevPeer := listOfPeers[i-1].Ledger.Accounts
		assert.True(t, reflect.DeepEqual(accountsOfPeer, accountsOfPrevPeer))

	}
}
func Test10AccountsHoldsMoney(t *testing.T) {

	noOfPeers := 10
	noOfNames := 10
	AccountBalance := 100
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer

	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < noOfNames; i++ {
		pkList[i] = listOfPeers[i%noOfPeers].CreateAccount()
		listOfPeers[i].CreateBalanceOnLedger(pkList[i], AccountBalance)
	}

	for i := 0; i < noOfPeers; i++ {
		accountsOfPeer := listOfPeers[i].Ledger.Accounts
		accountBalance := accountsOfPeer[pkList[i]]
		assert.Equal(t, AccountBalance, accountBalance, "Balance should match")

	}
}
func TestShouldNotBeAbleToHaveNegativeBalance(t *testing.T) {

	noOfPeers := 2
	noOfNames := 2

	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer

	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < noOfNames; i++ {
		pkList[i] = listOfPeers[i%noOfPeers].CreateAccount()
	}
	time.Sleep(250 * time.Millisecond)
	p1 := pkList[0]
	p2 := pkList[1]

	go listOfPeers[0].FloodSignedTransaction(p1, p2, 100)
	time.Sleep(250 * time.Millisecond)

	for i := 0; i < noOfPeers; i++ {
		accountsOfPeer := listOfPeers[i].Ledger.Accounts
		accountBalance := accountsOfPeer[pkList[i]]

		//println(accountBalance)
		assert.True(t, accountBalance <= 0, "Balance should be positive")

	}
}

func printControlLedger(controlLedger *ledger.Ledger) {
	l := controlLedger.Accounts
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	println()
	println("----------------Should be--------------")
	print("                            ")
	for _, value := range keys {
		print("[" + strconv.Itoa(l[value]) + "]")
	}
	print(" with the amount of SignedTransactions: " + strconv.Itoa(controlLedger.TA) + " ")
	println()
}

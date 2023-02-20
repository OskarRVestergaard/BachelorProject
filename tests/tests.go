package tests

import (
	"example.com/packages/ledger"
	"example.com/packages/peer"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"
)

func TestSignedAllValid(t *testing.T) {
	peersQt := 5
	tau := 10
	names := 5
	listOfPeers := make([]*peer.Peer, peersQt)

	var connectedPeers []string

	pkList := make([]string, names)

	for i := 0; i < peersQt; i++ {
		var p peer.Peer
		port := strconv.Itoa(18080 + i)
		listOfPeers[i] = &p
		p.RunPeer("127.0.0.1:" + port)

	}
	listOfPeers[0].Connect("Piplup is best water pokemon", 18079)
	connectedPeers = append(connectedPeers, listOfPeers[0].IpPort)
	time.Sleep(250 * time.Millisecond)
	for i := 1; i < peersQt; i++ {

		ipPort := connectedPeers[rand.Intn(len(connectedPeers))]
		ip := ipPort[0:(len(ipPort) - 6)]
		port := ipPort[len(ipPort)-5:]

		port2, _ := strconv.Atoi(port)
		listOfPeers[i].Connect(ip, port2)

		time.Sleep(250 * time.Millisecond)
	}
	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < names; i++ {
		pkList[i] = listOfPeers[i%peersQt].CreateAccount()
	}

	Led := new(ledger.Ledger)
	Led.TA = 0
	Led.Accounts = make(map[string]int)
	rand.Seed(time.Now().Unix())
	p1 := ""
	p2 := ""
	value := 0

	for j := 0; j < tau; j++ {
		for i := 0; i < peersQt; i++ {
			p1 = pkList[rand.Intn(names/peersQt)*peersQt+i]
			p2 = pkList[rand.Intn(names)]
			value = rand.Intn(100) + 1
			Led.UpdateLedger(p1, p2, value)

			go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
		}
	}
	time.Sleep(5000 * time.Millisecond)
	for i := 0; i < peersQt; i++ {
		listOfPeers[i].PrintLedger()
	}

	l := Led.Accounts
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
	print(" with the amount of SignedTransactions: " + strconv.Itoa(Led.TA) + " ")
	println()

	for i := 0; i < peersQt; i++ {
		accountsOfPeer := listOfPeers[i].Ledger.Accounts
		assert.True(t, reflect.DeepEqual(accountsOfPeer, l))
	}
}

func TestSignedOneNotValid(t *testing.T) {
	peersQt := 5
	tau := 10
	names := 5
	listOfPeers := make([]*peer.Peer, peersQt)

	var connectedPeers []string
	pkList := make([]string, names)

	for i := 0; i < peersQt; i++ {
		var p peer.Peer
		port := strconv.Itoa(18080 + i)
		listOfPeers[i] = &p
		p.RunPeer("127.0.0.1:" + port)

	}
	listOfPeers[0].Connect("Piplup is best water pokemon", 18079)
	connectedPeers = append(connectedPeers, listOfPeers[0].IpPort)
	time.Sleep(250 * time.Millisecond)
	for i := 1; i < peersQt; i++ {

		ipPort := connectedPeers[rand.Intn(len(connectedPeers))]
		ip := ipPort[0:(len(ipPort) - 6)]
		port := ipPort[len(ipPort)-5:]

		port2, _ := strconv.Atoi(port)
		listOfPeers[i].Connect(ip, port2)

		time.Sleep(250 * time.Millisecond)
	}
	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < names; i++ {
		pkList[i] = listOfPeers[i%peersQt].CreateAccount()
	}

	Led := new(ledger.Ledger)
	Led.TA = 0
	Led.Accounts = make(map[string]int)
	rand.Seed(time.Now().Unix())
	p1 := ""
	p2 := ""
	value := 0

	for j := 1; j < tau; j++ {
		for i := 0; i < peersQt; i++ {
			p1 = pkList[rand.Intn(names/peersQt)*peersQt+i]
			p2 = pkList[rand.Intn(names)]
			value = rand.Intn(100) + 1
			Led.UpdateLedger(p1, p2, value)
			go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
		}
	}

	p1 = pkList[rand.Intn(names/peersQt)*peersQt+1]
	p2 = pkList[rand.Intn(names)]
	value = rand.Intn(100) + 1
	Led.UpdateLedger(p1, p2, value)
	go listOfPeers[0].FloodSignedTransaction(p1, p2, value)

	for i := 1; i < peersQt; i++ {
		p1 = pkList[rand.Intn(names/peersQt)*peersQt+i]
		p2 = pkList[rand.Intn(names)]
		value = rand.Intn(100) + 1
		Led.UpdateLedger(p1, p2, value)
		go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
	}

	time.Sleep(5000 * time.Millisecond)
	for i := 0; i < peersQt; i++ {
		listOfPeers[i].PrintLedger()
	}

	l := Led.Accounts
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	println()
	println("----------------Unvalidated be--------------")
	print("                            ")
	for _, value := range keys {
		print("[" + strconv.Itoa(l[value]) + "]")
	}
	print(" with the amount of unvalidated Transactions: " + strconv.Itoa(Led.TA) + " ")
	println()

	for i := 0; i < peersQt; i++ {
		signedTransactionsOfPeer := listOfPeers[i].Ledger.TA

		assert.Equal(t, signedTransactionsOfPeer, peersQt*tau-1, "One msg was not signed but still validated")

	}

}

func TestSignedAllRandom(t *testing.T) {
	peersQt := 5
	tau := 15
	names := 5
	listOfPeers := make([]*peer.Peer, peersQt)

	var connectedPeers []string
	pkList := make([]string, names)

	for i := 0; i < peersQt; i++ {
		var p peer.Peer
		port := strconv.Itoa(18080 + i)
		listOfPeers[i] = &p
		p.RunPeer("127.0.0.1:" + port)

	}
	listOfPeers[0].Connect("Piplup is best water pokemon", 18079)
	connectedPeers = append(connectedPeers, listOfPeers[0].IpPort)
	time.Sleep(250 * time.Millisecond)
	for i := 1; i < peersQt; i++ {

		ipPort := connectedPeers[rand.Intn(len(connectedPeers))]
		ip := ipPort[0:(len(ipPort) - 6)]
		port := ipPort[len(ipPort)-5:]

		port2, _ := strconv.Atoi(port)
		listOfPeers[i].Connect(ip, port2)

		time.Sleep(250 * time.Millisecond)
	}
	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < names; i++ {
		pkList[i] = listOfPeers[i%peersQt].CreateAccount()
	}

	Led := new(ledger.Ledger)
	Led.TA = 0
	Led.Accounts = make(map[string]int)
	rand.Seed(time.Now().Unix())
	p1 := ""
	p2 := ""
	value := 0

	for j := 0; j < tau; j++ {
		for i := 0; i < peersQt; i++ {
			p1 = pkList[rand.Intn(names)]
			p2 = pkList[rand.Intn(names)]
			value = rand.Intn(100) + 1
			Led.UpdateLedger(p1, p2, value)
			go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
		}
	}

	time.Sleep(5000 * time.Millisecond)
	for i := 0; i < peersQt; i++ {
		listOfPeers[i].PrintLedger()
	}

	l := Led.Accounts
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	println()
	println("----------------Unvalidated be--------------")
	print("                            ")
	for _, value := range keys {
		print("[" + strconv.Itoa(l[value]) + "]")
	}
	print(" with the amount of unvalidated Transactions: " + strconv.Itoa(Led.TA) + " ")
	println()

	for i := 1; i < peersQt; i++ {
		accountsOfPeer := listOfPeers[i].Ledger.Accounts
		accountsOfPrevPeer := listOfPeers[i-1].Ledger.Accounts
		assert.True(t, reflect.DeepEqual(accountsOfPeer, accountsOfPrevPeer))

	}

}

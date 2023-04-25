package Peer

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"sort"
	"strconv"
)

func (p *Peer) CreateAccount() string {
	secretKey, publicKey := p.signatureStrategy.KeyGen()
	keys := <-p.publicToSecret
	keys[publicKey] = secretKey
	p.publicToSecret <- keys
	return publicKey
}

// CreateBalanceOnLedger for testing only
// TODO Remove
func (p *Peer) CreateBalanceOnLedger(pk string, amount int) {

	p.Ledger.Mutex.Lock()
	p.Ledger.Accounts[pk] += amount
	p.Ledger.Mutex.Unlock()
}

func (p *Peer) UpdateLedger(transactions []*blockchain.SignedTransaction) {

	p.Ledger.Mutex.Lock()
	for _, trans := range transactions {
		p.Ledger.TA = p.Ledger.TA + 1
		p.Ledger.Accounts[trans.From] -= trans.Amount
		p.Ledger.Accounts[trans.To] += trans.Amount
	}
	p.Ledger.Mutex.Unlock()
}

func (p *Peer) PrintLedger() {
	p.Ledger.Mutex.Lock()
	l := p.Ledger.Accounts
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	println()
	print("Ledger of: " + p.network.GetAddress().ToString() + ": ")
	for _, value := range keys {
		print("[" + strconv.Itoa(l[value]) + "]")
	}
	println(" with the amount of SignedTransactions: " + strconv.Itoa(p.Ledger.TA) + " and deniedTransactions: " + strconv.Itoa(p.Ledger.UTA))
	p.Ledger.Mutex.Unlock()
}

func MakeLedger() *models.Ledger {
	ledger := new(models.Ledger)
	ledger.Accounts = make(map[string]int)
	return ledger
}

// TODO Eventually this needs to also be upheld
func (p *Peer) validTransaction(from string, amount int) bool {
	if amount == 0 {
		println("Invalid SignedTransaction with the amount 0")
		return false
	} else if p.Ledger.Accounts[from] < amount {
		println("Account should hold the transaction amount")
		return false
	}
	return true
}

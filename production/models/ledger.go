package models

import (
	"sync"
)

type Ledger struct {
	Accounts map[string]int
	Mutex    sync.Mutex
	TA       int
	UTA      int
}

func (l *Ledger) UpdateLedger(from string, to string, value int) {
	l.Mutex.Lock()

	l.Accounts[from] -= value
	l.Accounts[to] += value

	l.TA = l.TA + 1
	l.Mutex.Unlock()
}

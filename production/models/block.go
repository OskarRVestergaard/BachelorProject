package models

import (
	"crypto/sha256"
	"fmt"
	"github.com/OskarRVestergaard/BachelorProject/production/models/messages"
)

type Block struct {
	SlotNumber   int
	Hash         string
	PreviousHash string
	Transactions []*messages.SignedTransaction

	//TransactionsLog map[int]string
}
type Blockchain struct {
	blockchain map[string]*Block
}

var slot = 0

func MakeBlock(transactions []*messages.SignedTransaction, prevHash string) Block {
	//TODO add maximum blockSize
	var b Block
	//b.slotNumber = slot
	b.PreviousHash = prevHash
	//b.TransactionsLog = transactions
	b.Transactions = transactions
	b.Hash = calculateHash(b.PreviousHash, b.Transactions)
	slot += 1
	return b

}

func calculateHash(PreviousHash string, transactions []*messages.SignedTransaction) string {
	h := sha256.New()

	transactionsString := ConvertToString(transactions)
	h.Write([]byte((PreviousHash + transactionsString)))
	return string(h.Sum(nil))
}

func isValid(block Block, previousBlock Block) bool {
	//if previousBlock.slotNumber >= block.slotNumber {
	//	return false
	//}
	if previousBlock.Hash != block.PreviousHash {
		return false
	}
	if block.Hash != calculateHash(block.PreviousHash, block.Transactions) {
		return false
	}
	return true
}

func ConvertToString(transactions []*messages.SignedTransaction) string {
	var s string
	for key, value := range transactions {
		s = fmt.Sprintf("%v:%s", key, value)

	}
	return s

	//TODO should consider doing it wothg join
}
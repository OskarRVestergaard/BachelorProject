package block

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	slotNumber   int
	Hash         string
	PreviousHash string
	//TransactionsLog map[int]*models.SignedTransaction
	TransactionsLog map[int]string
}
type Blockchain struct {
	blockchain []Block
}

var Blockchain2 []Block

var slot = 1

func MakeBlock(transactions map[int]string, prevHash string) Block {
	//TODO add maximum blockSize
	var b Block
	b.slotNumber = slot
	b.PreviousHash = prevHash
	b.TransactionsLog = transactions
	b.Hash = calculateHash(b.PreviousHash, b.TransactionsLog)
	slot += 1
	return b

}

func calculateHash(PreviousHash string, TransactionsLog map[int]string) string {
	h := sha256.New()

	transactionsString := ConvertToString(TransactionsLog)
	h.Write([]byte((PreviousHash + transactionsString)))
	println(transactionsString)
	return string(h.Sum(nil))
}

func isValid(block Block, previousBlock Block) bool {
	if previousBlock.slotNumber >= block.slotNumber {
		return false
	}
	if previousBlock.Hash != block.PreviousHash {
		return false
	}
	if block.Hash != calculateHash(block.PreviousHash, block.TransactionsLog) {
		return false
	}
	return true
}

func ConvertToString(transactions map[int]string) string {
	var s string
	for key, value := range transactions {
		s = fmt.Sprintf("%v:%s", key, value)

	}
	return s
}

//func main() {
//	now := time.Now()
//	println(now)
//}

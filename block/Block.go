package block

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Hash         string
	PreviousHash string
	//TransactionsLog map[int]*models.SignedTransaction
	TransactionsLog map[int]string
}

func (b *Block) MakeBlock(transactions map[int]string, prevHash string) {
	//TODO add maximum blockSize
	b.PreviousHash = prevHash
	b.TransactionsLog = transactions
	b.Hash = b.calculateHash()

}

func (b *Block) calculateHash() string {
	h := sha256.New()

	transactionsString := ConvertToString(b.TransactionsLog)
	h.Write([]byte((b.PreviousHash + transactionsString)))
	println(transactionsString)
	return string(h.Sum(nil))
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

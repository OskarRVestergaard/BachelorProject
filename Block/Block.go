package Block

import (
	"crypto/sha256"
	"example.com/packages/models"
	"time"
)

type Block struct {
	Hash            string
	PreviousHash    string
	TransactionsLog map[int]*models.SignedTransaction
	//TimeStamp      time.Time
}

func (b *Block) MakeBlock(transactions []*models.SignedTransaction, prevHash string) {
	b.PreviousHash = prevHash
	b.TransactionsLog = make(map[int]*models.SignedTransaction)
	//b.TransactionsLog = transactions
	b.Hash = b.calculateHash()
	//b.calculateHash()
	//b.TimeStamp = time.Now()

}

func (b *Block) calculateHash() string {
	h := sha256.New()
	h.Write([]byte("hello world\n"))
	h.Write([]byte(b.PreviousHash))
	//calculatedHash := crypto.SHA256(b.PreviousHash)
	//b.Hash = "asd"
	//return h.Sum(nil)
	return string(h.Sum(nil))
}
func main() {
	now := time.Now()
	println(now)
}

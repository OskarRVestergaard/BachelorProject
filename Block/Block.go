package Block

import (
	"example.com/packages/models"
	"time"
)

type Block struct {
	Hash           string
	PreviousHash   string
	TrasactionsLog map[models.SignedTransaction]string
	TimeStamp      time.Time
}

func (b *Block) MakeBlock() {
	b.Hash = calculateHash()
	b.TimeStamp = time.Now()

}

func calculateHash() string {
	return cry
}
func main() {
	now := time.Now()
	println(now)
}

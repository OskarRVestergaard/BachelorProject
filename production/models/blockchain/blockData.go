package blockchain

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/messages"
	"strconv"
)

type BlockData struct {
	Hardness     int
	Transactions []messages.SignedTransaction
}

func (blockData *BlockData) ToString() string {
	return strconv.Itoa(blockData.Hardness)
}

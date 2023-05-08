package models

import (
	"strconv"
)

type BlockData struct {
	Hardness     int
	Transactions []SignedTransaction
}

func (blockData *BlockData) ToString() string {
	return strconv.Itoa(blockData.Hardness)
}

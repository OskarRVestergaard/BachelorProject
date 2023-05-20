package PoWblockchain

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"strconv"
)

type BlockData struct {
	Hardness     int
	Transactions []models.SignedPaymentTransaction
}

func (blockData *BlockData) ToString() string {
	return strconv.Itoa(blockData.Hardness)
}

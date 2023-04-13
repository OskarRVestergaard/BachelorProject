package blockchain

import "strconv"

type BlockData struct {
	Hardness int
}

func (blockData *BlockData) ToString() string {
	return strconv.Itoa(blockData.Hardness)
}

package blockchain

import (
	"bytes"
	"math/big"
	"strconv"
)

type Block struct {
	IsGenesis bool      //True only if the block is the genesis block
	Vk        big.Int   //verification key
	Slot      int       //slot number
	Draw      string    //winner ticket
	U         BlockData //Block data
	H         []byte    //block hash of some previous hash
	Sigma     string    //signature
}

/*
GetVal

returns the val of the block to be used for PathWeight calculations,
and also true if it is genesis (to be treated as infinite)
*/
func (block *Block) GetVal() (val string, isGenesis bool) {
	if block.IsGenesis {
		return "Genesis", true
	}
	return block.Vk.String() + strconv.Itoa(block.Slot) + block.Draw, false
}

/*
ToByteArray

returns a byte array representation of the block to be used for hashing
*/
func (block *Block) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(block.Vk.String())
	buffer.WriteString(strconv.Itoa(block.Slot))
	buffer.WriteString(block.Draw)
	buffer.WriteString(block.U.ToString())
	buffer.WriteString(string(block.H))
	buffer.WriteString(block.Sigma)

	return buffer.Bytes()
}

/*
CreateGenesisBlock

Creates the default Genesis-block to be used in a blocktree
*/
func CreateGenesisBlock() Block {
	return Block{
		IsGenesis: true,
		Vk:        big.Int{},
		Slot:      0,
		Draw:      "",
		U: BlockData{
			Hardness: 8,
		},
		H:     nil,
		Sigma: "",
	}
}

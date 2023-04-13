package blockchain

import (
	"bytes"
	"math/big"
	"strconv"
)

type Block struct {
	Vk    big.Int //verification key
	Slot  int     //slot number
	Draw  string  //winner ticket
	U     string  //Block data
	H     []byte  //block hash of some previous hash
	Sigma string  //signature
}

func (block *Block) GetVal() string {
	return block.Vk.String() + strconv.Itoa(block.Slot) + block.Draw
}

func (block *Block) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(block.Vk.String())
	buffer.WriteString(strconv.Itoa(block.Slot))
	buffer.WriteString(block.Draw)
	buffer.WriteString(block.U)
	buffer.WriteString(string(block.H))
	buffer.WriteString(block.Sigma)

	return buffer.Bytes()
}

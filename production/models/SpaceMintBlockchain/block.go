package PoWblockchain

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/peer_strategy/SpacemintPeer"
	"strconv"
)

type Block struct {
	IsGenesis           bool //True only if the block is the genesis block
	ParentHash          sha256.HashValue
	HashSubBlock        HashSubBlock
	TransactionSubBlock TransactionSubBlock
	SignatureSubBlock   SignatureSubBlock
}

type HashSubBlock struct {
	Slot                      int                              //index or slot number
	SignatureOnParentSubBlock []byte                           //Signature linking this block to its parent
	Draw                      SpacemintPeer.PoSpaceLotteryDraw //The Proof of space associated with the block
}

func (subBlock *HashSubBlock) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(subBlock.Slot))
	buffer.WriteString(";_;")
	buffer.Write(subBlock.SignatureOnParentSubBlock)
	buffer.WriteString(";_;")
	buffer.Write(subBlock.Draw.ToByteArray())
	return buffer.Bytes()
}

type TransactionSubBlock struct {
	Slot         int                   //index or slot number
	Transactions SpacemintTransactions //List of transactions
}

func (subBlock *TransactionSubBlock) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(subBlock.Slot))
	buffer.WriteString(";_;")
	buffer.Write(subBlock.Transactions.ToByteArray())
	return buffer.Bytes()
}

type SignatureSubBlock struct {
	Slot                                  int //index or slot number
	SignatureOnCurrentTransactionSubBlock []byte
	SignatureOnParentSubBlock             []byte
}

func (subBlock *SignatureSubBlock) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(subBlock.Slot))
	buffer.WriteString(";_;")
	buffer.Write(subBlock.SignatureOnCurrentTransactionSubBlock)
	buffer.WriteString(";_;")
	buffer.Write(subBlock.SignatureOnParentSubBlock)
	return buffer.Bytes()
}

func (block *Block) GetQuality() (value sha256.HashValue, isGenesis bool) {
	if block.IsGenesis {
		return sha256.HashValue{}, true
	}
	return block.HashOfBlock(), false //TODO Return a proper quality, (which can actually only be done with information from parents)
}

/*
HashOfBlock

returns a byte array representation of the block to be used for hashing
*/
func (block *Block) HashOfBlock() sha256.HashValue {
	byteArrayString := block.ToByteArray()
	hash := sha256.HashByteArray(byteArrayString)
	return hash
}

/*
ToByteArray

returns a byte array representation, if you want the hash use HashOfBlock instead
*/
func (block *Block) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.Write(block.ParentHash.ToSlice())
	buffer.WriteString(";_;")
	buffer.Write(block.HashSubBlock.ToByteArray())
	buffer.WriteString(";_;")
	buffer.Write(block.TransactionSubBlock.ToByteArray())
	buffer.WriteString(";_;")
	buffer.Write(block.SignatureSubBlock.ToByteArray())
	return buffer.Bytes()
}

/*
CreateGenesisBlock

Creates the default Genesis-block to be used in a blocktree
*/
func CreateGenesisBlock() Block {
	return Block{
		IsGenesis:  true,
		ParentHash: sha256.HashValue{},
	}
}

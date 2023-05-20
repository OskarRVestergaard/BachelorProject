package SpaceMintBlockchain

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy/PoSpace"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
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
	Slot                          int                 //index or slot number
	SignatureOnParentHashSubBlock []byte              //Signature linking this block to its parent
	Draw                          PoSpace.LotteryDraw //The Proof of space associated with the block
}

func (block *Block) HasConsistentSignaturesAndSlots(parentBlock Block, signatureStrategy signature_strategy.SignatureInterface) bool {
	if block.HashSubBlock.Slot != block.TransactionSubBlock.Slot {
		return false
	}
	if block.SignatureSubBlock.Slot != block.TransactionSubBlock.Slot {
		return false
	}
	vk := block.HashSubBlock.Draw.Vk
	if !signatureStrategy.Verify(vk, block.TransactionSubBlock.ToByteArray(), block.SignatureSubBlock.SignatureOnCurrentTransactionSubBlock) {
		return false
	}
	if !signatureStrategy.Verify(vk, parentBlock.SignatureSubBlock.ToByteArray(), block.SignatureSubBlock.SignatureOnParentSignatureSubBlock) {
		return false
	}
	if !signatureStrategy.Verify(vk, parentBlock.HashSubBlock.ToByteArray(), block.HashSubBlock.SignatureOnParentHashSubBlock) {
		return false
	}
	return true
}

func (block *Block) SignBlock(parentBlock Block, signatureStrategy signature_strategy.SignatureInterface, secretKey string) {
	currentTransactionSubBlockSignature := signatureStrategy.Sign(block.TransactionSubBlock.ToByteArray(), secretKey)
	prevSignatureSubBlockSignature := signatureStrategy.Sign(parentBlock.SignatureSubBlock.ToByteArray(), secretKey)
	prevHashSubBlockSignature := signatureStrategy.Sign(parentBlock.HashSubBlock.ToByteArray(), secretKey)
	block.SignatureSubBlock.SignatureOnCurrentTransactionSubBlock = currentTransactionSubBlockSignature
	block.SignatureSubBlock.SignatureOnParentSignatureSubBlock = prevSignatureSubBlockSignature
	block.HashSubBlock.SignatureOnParentHashSubBlock = prevHashSubBlockSignature
}

func (subBlock *HashSubBlock) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(subBlock.Slot))
	buffer.WriteString(";_;")
	buffer.Write(subBlock.SignatureOnParentHashSubBlock)
	buffer.WriteString(";_;")
	buffer.Write(subBlock.Draw.ToByteArray())
	return buffer.Bytes()
}

func (subBlock *HashSubBlock) HashOfBlock() sha256.HashValue {
	byteArrayString := subBlock.ToByteArray()
	hash := sha256.HashByteArray(byteArrayString)
	return hash
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
	SignatureOnParentSignatureSubBlock    []byte
}

func (subBlock *SignatureSubBlock) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(subBlock.Slot))
	buffer.WriteString(";_;")
	buffer.Write(subBlock.SignatureOnCurrentTransactionSubBlock)
	buffer.WriteString(";_;")
	buffer.Write(subBlock.SignatureOnParentSignatureSubBlock)
	return buffer.Bytes()
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
	buffer.WriteString(strconv.FormatBool(block.IsGenesis))
	buffer.WriteString(";_;")
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

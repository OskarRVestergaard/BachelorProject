package PoWblockchain

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/peer_strategy/SpacemintPeer"
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
	Slot                      int                              //index or slot number
	SignatureOnParentSubBlock []byte                           //Signature linking this block to its parent
	Draw                      SpacemintPeer.PoSpaceLotteryDraw //The Proof of space associated with the block
}

type TransactionSubBlock struct {
	Slot      int              //index or slot number
	BlockData models.BlockData //List of transactions TODO Change to the "spacemint" transactions
}

type SignatureSubBlock struct {
	Slot                                  int //index or slot number
	SignatureOnCurrentTransactionSubBlock []byte
	SignatureOnParentSubBlock             []byte
}

/*
GetVal

returns the val of the block to be used for PathWeight calculations,
and also true if it is genesis (to be treated as infinite)
*/
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
	buffer.WriteString(strconv.Itoa(block.Slot))
	buffer.WriteString(";_;")
	buffer.WriteString(block.Draw.ToString())
	buffer.WriteString(";_;")
	buffer.WriteString(block.BlockData.ToString())
	buffer.WriteString(";_;")
	buffer.Write(block.ParentHash.ToSlice())
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

func (block *Block) SignBlock(signatureStrategy signature_strategy.SignatureInterface, secretSigningKey string) {
	data := block.toByteArrayWithoutSign()
	hashedData := sha256.HashByteArray(data).ToSlice()
	signature := signatureStrategy.Sign(hashedData, secretSigningKey)
	block.Signature = signature
}

func (block *Block) HasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface) bool {
	blockVerificationKey := block.Vk
	blockHashWithoutSign := sha256.HashByteArray(block.toByteArrayWithoutSign()).ToSlice()
	blockSignature := block.Signature
	result := signatureStrategy.Verify(blockVerificationKey, blockHashWithoutSign, blockSignature)
	return result
}

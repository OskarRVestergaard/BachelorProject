package blockchain

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"strconv"
)

type Block struct {
	IsGenesis  bool      //True only if the block is the genesis block
	Vk         string    //verification key
	Slot       int       //slot number
	Draw       string    //winner ticket
	BlockData  BlockData //Block data
	ParentHash []byte    //block hash of some previous hash
	Signature  []byte    //signature
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
	return block.Vk + strconv.Itoa(block.Slot) + block.Draw, false
}

/*
HashOfBlock

returns a byte array representation of the block to be used for hashing'
//TODO Make into a hashing function instead? maybe pass hashing strategy, when this is interfaced
*/
func (block *Block) HashOfBlock() []byte {
	byteArrayString := block.ToByteArray()
	hash := hash_strategy.HashByteArray(byteArrayString)
	return hash
}

/*
ToByteArray

returns a byte array representation, if you want the hash use HashOfBlock instead
*/
//TODO ADD SOME SORT OF SEPERATOR BETWEEN THEM, SINCE THIS IS ONLY ONE WAY, AND CAN BE EXPLOITED
func (block *Block) ToByteArray() []byte {
	var firstBytes = block.toByteArrayWithoutSign()
	firstBytes = append(firstBytes, block.Signature...)
	return firstBytes
}

/*
ToByteArrayWithoutSign

returns a byte array representation, to be used for signature calculation
*/
//TODO ADD SOME SORT OF SEPERATOR BETWEEN THEM, SINCE THIS IS ONLY ONE WAY, AND CAN BE EXPLOITED
func (block *Block) toByteArrayWithoutSign() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(block.Vk)
	buffer.WriteString(strconv.Itoa(block.Slot))
	buffer.WriteString(block.Draw)
	buffer.WriteString(block.BlockData.ToString())
	buffer.WriteString(string(block.ParentHash))

	return buffer.Bytes()
}

/*
CreateGenesisBlock

Creates the default Genesis-block to be used in a blocktree
*/
func CreateGenesisBlock() Block {
	return Block{
		IsGenesis: true,
		Vk:        "",
		Slot:      0,
		Draw:      "",
		BlockData: BlockData{
			Hardness: 8,
		},
		ParentHash: nil,
		Signature:  nil,
	}
}

func (block *Block) CalculateSignature(signatureStrategy signature_strategy.SignatureInterface, secretSigningKey string) int {
	data := block.toByteArrayWithoutSign()
	hashedData := hash_strategy.HashByteArray(data)
	signature := signatureStrategy.Sign(hashedData, secretSigningKey)
	block.Signature = signature
	return 1
}

func (block *Block) HasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface) bool {
	blockVerificationKey := block.Vk
	blockHashWithoutSign := hash_strategy.HashByteArray(block.toByteArrayWithoutSign())
	blockSignature := block.Signature
	result := signatureStrategy.Verify(blockVerificationKey, blockHashWithoutSign, blockSignature)
	return result
}

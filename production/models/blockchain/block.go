package blockchain

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"math/big"
	"strconv"
)

type Block struct {
	IsGenesis  bool      //True only if the block is the genesis block
	Vk         big.Int   //verification key
	Slot       int       //slot number
	Draw       string    //winner ticket
	BlockData  BlockData //Block data
	ParentHash []byte    //block hash of some previous hash
	Signature  big.Int   //signature
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
func (block *Block) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(block.Vk.String())
	buffer.WriteString(strconv.Itoa(block.Slot))
	buffer.WriteString(block.Draw)
	buffer.WriteString(block.BlockData.ToString())
	buffer.WriteString(string(block.ParentHash))
	buffer.WriteString(block.Signature.String())

	return buffer.Bytes()
}

/*
ToByteArrayWithoutSign

returns a byte array representation, to be used for signature calculation
*/
func (block *Block) toByteArrayWithoutSign() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(block.Vk.String())
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
		Vk:        big.Int{},
		Slot:      0,
		Draw:      "",
		BlockData: BlockData{
			Hardness: 8,
		},
		ParentHash: nil,
		Signature:  big.Int{},
	}
}

func (block *Block) CalculateSignature(signatureStrategy signature_strategy.SignatureInterface, secretSigningKey string) big.Int {
	data := block.toByteArrayWithoutSign()
	signature := signatureStrategy.Sign(data, secretSigningKey)
	return *signature
}

func (block *Block) HasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface) bool {
	blockVerificationKey := block.Vk.String()
	blockHashWithoutSign := hash_strategy.HashByteArray(block.toByteArrayWithoutSign())
	blockSignature := block.Signature
	result := signatureStrategy.Verify(blockVerificationKey, blockHashWithoutSign, &blockSignature)
	return result
}

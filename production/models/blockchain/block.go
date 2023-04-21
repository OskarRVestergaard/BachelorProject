package blockchain

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"strconv"
)

type Block struct {
	IsGenesis  bool                                  //True only if the block is the genesis block
	Vk         string                                //verification key
	Slot       int                                   //slot number
	Draw       lottery_strategy.WinningLotteryParams //winner ticket
	BlockData  BlockData                             //Block data
	ParentHash []byte                                //block hash of some previous hash
	Signature  []byte                                //signature
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
	return block.Vk + strconv.Itoa(block.Slot) + block.Draw.ToString(), false
}

/*
HashOfBlock

returns a byte array representation of the block to be used for hashing
*/
func (block *Block) HashOfBlock() []byte {
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
	buffer.Write(block.toByteArrayWithoutSign())
	buffer.WriteString(";_;")
	buffer.Write(block.Signature)
	return buffer.Bytes()
}

/*
ToByteArrayWithoutSign

returns a byte array representation, to be used for signature calculation
*/
func (block *Block) toByteArrayWithoutSign() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(block.Vk)
	buffer.WriteString(";_;")
	buffer.WriteString(strconv.Itoa(block.Slot))
	buffer.WriteString(";_;")
	buffer.WriteString(block.Draw.ToString())
	buffer.WriteString(";_;")
	buffer.WriteString(block.BlockData.ToString())
	buffer.WriteString(";_;")
	buffer.Write(block.ParentHash)
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
		Draw: lottery_strategy.WinningLotteryParams{
			Vk:         "",
			ParentHash: nil,
			Counter:    0,
		},
		BlockData: BlockData{
			Hardness: 24,
		},
		ParentHash: nil,
		Signature:  nil,
	}
}

func (block *Block) SignBlock(signatureStrategy signature_strategy.SignatureInterface, secretSigningKey string) {
	data := block.toByteArrayWithoutSign()
	hashedData := sha256.HashByteArray(data)
	signature := signatureStrategy.Sign(hashedData, secretSigningKey)
	block.Signature = signature
}

func (block *Block) HasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface) bool {
	blockVerificationKey := block.Vk
	blockHashWithoutSign := sha256.HashByteArray(block.toByteArrayWithoutSign())
	blockSignature := block.Signature
	result := signatureStrategy.Verify(blockVerificationKey, blockHashWithoutSign, blockSignature)
	return result
}

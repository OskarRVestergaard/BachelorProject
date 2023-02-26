package test

import (
	"crypto/sha256"
	"example.com/packages/block"
	"example.com/packages/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockHash(t *testing.T) {
	//var b block.Block
	transactionsLog := make(map[int]string)
	transactionsLog[1] = "a->b:100"
	prevHash := "hejhej"
	var b = block.MakeBlock(transactionsLog, prevHash)

	h := sha256.New()

	h.Write([]byte((prevHash + block.ConvertToString(transactionsLog))))
	assert.Equal(t, b.Hash, string(h.Sum(nil)), "hashes match")
}
func TestBlockchainLengthOfGenesis(t *testing.T) {
	var blockChain = makeGenesisBlockchain()
	assert.Equal(t, 1, len(blockChain), "blockchain length should be 1")
}

func TestPeerCanAppendBlockToBlockchain(t *testing.T) {

	//var blockChain = makeGenesisBlockchain()
	//
	//println(blockChain)
	//blockChain = append(blockChain, genesisBlock)

	println(utils.SignedTransaction)
}
func makeGenesisBlockchain() []block.Block {
	genesisBlock := block.Block{
		Hash:            "GenesisBlock",
		PreviousHash:    "GenesisBlock",
		TransactionsLog: nil,
	}
	var blockChain = block.Blockchain2
	blockChain = append(blockChain, genesisBlock)
	return blockChain
}

//func Test

package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockTransactionsOnBlocksAndOwn(t *testing.T) {
	var Unhandled []blockchain.SignedTransaction

	var blockTree = blockchain.NewBlocktree(blockchain.CreateGenesisBlock())
	for i := 0; i < 100; i++ {
		var b = blockchain.SignedTransaction{
			Id:        uuid.New(),
			From:      "",  //not relevant for test
			To:        "",  //not relevant for test
			Amount:    0,   //not relevant for test
			Signature: nil, //not relevant for test
		}
		Unhandled = append(Unhandled, b)
		parent := blockTree.GetHead()
		if i%3 == 0 {
			var trans []blockchain.SignedTransaction
			trans = append(trans, b)
			blockTree.AddBlock(blockchain.Block{
				IsGenesis: false,
				Vk:        "",
				Slot:      0,
				Draw:      "",
				BlockData: blockchain.BlockData{
					Hardness:     0,
					Transactions: trans,
				},
				ParentHash: parent.HashOfBlock(),
				Signature:  nil,
			})
		}
	}
	dif := blockTree.GetTransactionsNotInTree(Unhandled)
	assert.Equal(t, 66, len(dif))

}

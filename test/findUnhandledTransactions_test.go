package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockTransactionsOnBlocksAndOwn(t *testing.T) {
	var Unhandled []blockchain.SignedTransaction

	var blockTree, blockTreeCreationWentWell = blockchain.NewBlocktree(blockchain.CreateGenesisBlock())
	assert.True(t, blockTreeCreationWentWell)
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
				Draw: lottery_strategy.WinningLotteryParams{
					Vk:         "DEBUG",
					ParentHash: sha256.HashValue{},
					Counter:    0,
				},
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

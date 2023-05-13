package temp2

import (
	"encoding/binary"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy/PoW"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMine(t *testing.T) {
	miner := PoW.PoW{}
	blockToExtend := make([]byte, 4)
	binary.LittleEndian.PutUint32(blockToExtend, 31415926)
	blocksChannel := make(chan []byte, 0)
	resultChannel := make(chan lottery_strategy.WinningLotteryParams, 0)
	vk := "4"
	hardness := 20

	miner.StartNewMiner(vk, hardness, blockToExtend, blocksChannel, resultChannel)

	blocksChannel <- blockToExtend
	result := <-resultChannel

	isVerified := miner.Verify(vk, result.ParentHash, hardness, result.Counter)

	assert.True(t, isVerified)
}

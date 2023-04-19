package temp2

import (
	"encoding/binary"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"testing"
)

func TestMine(t *testing.T) {
	miner := lottery_strategy.PoW{}
	blockToExtend := make([]byte, 4)
	binary.LittleEndian.PutUint32(blockToExtend, 31415926)

	found, big := miner.Mine("2", blockToExtend)
	print(found)
	print(big)
}

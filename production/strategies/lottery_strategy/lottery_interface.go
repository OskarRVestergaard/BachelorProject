package lottery_strategy

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/sha256"
	"strconv"
)

type LotteryInterface interface {
	StartNewMiner(vk string, hardness int, initialHash sha256.HashValue, newBlockHashes chan sha256.HashValue, winningDraws chan WinningLotteryParams, stopMinerSignal chan struct{})
	Verify(vk string, parentHash sha256.HashValue, hardness int, counter int) bool
}

type WinningLotteryParams struct {
	Vk         string
	ParentHash sha256.HashValue
	Counter    int
}

func (p WinningLotteryParams) ToByteSlice() []byte {
	result := sha256.ToSlice(p.ParentHash)
	result = append(result, p.Vk...)
	counterBytes := strconv.Itoa(p.Counter)
	result = append(result, counterBytes...)
	return result
}

func (p WinningLotteryParams) ToString() string {
	buf := bytes.NewBuffer(p.ToByteSlice())
	return buf.String()
}

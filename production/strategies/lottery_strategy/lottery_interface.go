package lottery_strategy

import (
	"bytes"
	"strconv"
)

type LotteryInterface interface {
	StartNewMiner(vk string, hardness int, initialHash []byte, newBlockHashes chan []byte, winningDraws chan WinningLotteryParams)
	Verify(vk string, parentHash []byte, hardness int, counter int) bool
}

type WinningLotteryParams struct {
	Vk         string
	ParentHash []byte
	Counter    int
}

func (p WinningLotteryParams) ToByteArray() []byte {
	result := p.ParentHash
	result = append(result, p.Vk...)
	counterBytes := strconv.Itoa(p.Counter)
	result = append(result, counterBytes...)
	return result
}

func (p WinningLotteryParams) ToString() string {
	buf := bytes.NewBuffer(p.ToByteArray())
	return buf.String()
}

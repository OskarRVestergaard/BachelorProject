package lottery_strategy

import (
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
)

type PoW struct {
}

type internalCombination struct {
	minerShouldContinue bool
	parentHash          sha256.HashValue
}

func (lottery *PoW) StartNewMiner(vk string, hardness int, initialHash sha256.HashValue, newBlockHashes chan sha256.HashValue, winningDraws chan WinningLotteryParams, stopMinerSignal chan struct{}) {
	newBlockHashesInternal := make(chan internalCombination)
	lottery.combineChannels(newBlockHashes, stopMinerSignal, newBlockHashesInternal)
	go lottery.startNewMinerInternal(vk, hardness, initialHash, newBlockHashesInternal, winningDraws)
}

func (lottery *PoW) combineChannels(newHashes chan sha256.HashValue, stopMiner chan struct{}, internalStruct chan internalCombination) {
	go func() {
		for {
			newParentHash := <-newHashes
			combination := internalCombination{
				minerShouldContinue: true,
				parentHash:          newParentHash,
			}
			internalStruct <- combination
		}
	}()
	go func() {
		for {
			_ = <-stopMiner
			combination := internalCombination{
				minerShouldContinue: false,
				parentHash:          sha256.HashValue{},
			}
			internalStruct <- combination
		}
	}()
}

func (lottery *PoW) startNewMinerInternal(vk string, hardness int, initialHash sha256.HashValue, newBlockHashesInternal chan internalCombination, winningDraws chan WinningLotteryParams) {
	internalStruct := internalCombination{
		minerShouldContinue: true,
		parentHash:          initialHash,
	}
	for internalStruct.minerShouldContinue {
		done := make(chan struct{})
		go lottery.mineOnSingleBlock(vk, internalStruct.parentHash, hardness, done, winningDraws)
		internalStruct = <-newBlockHashesInternal
		done <- struct{}{}
	}
}

func (lottery *PoW) mineOnSingleBlock(vk string, parentHash sha256.HashValue, hardness int, done chan struct{}, winningDraws chan WinningLotteryParams) {
	c := 0
	for {
		select {
		case <-done:
			return
		default:
			c = c + 1
			draw := WinningLotteryParams{
				Vk:         vk,
				ParentHash: parentHash,
				Counter:    c,
			}
			hashOfTicket := sha256.HashByteArray(draw.ToByteSlice())
			if verifyPoW(hashOfTicket, hardness) {
				winningDraws <- draw
				_ = <-done
				return
			}
		}
	}

}

func verifyPoW(hashedTicket sha256.HashValue, hardness int) bool {
	byteAmount := hardness / 8
	restAmount := hardness - 8*byteAmount
	for i := 0; i < byteAmount; i++ {
		if hashedTicket[i] != 0 {
			return false
		}
	}
	byteToCheck := hashedTicket[byteAmount]
	for i := 0; i < restAmount; i++ {
		if (byteToCheck >> (7 - i)) != 0 {
			return false
		}
	}
	return true
}

func (lottery *PoW) Verify(vk string, parentHash sha256.HashValue, hardness int, counter int) bool {
	draw := WinningLotteryParams{
		Vk:         vk,
		ParentHash: parentHash,
		Counter:    counter,
	}
	hashed := sha256.HashByteArray(draw.ToByteSlice())

	return verifyPoW(hashed, hardness)
}

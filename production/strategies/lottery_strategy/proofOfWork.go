package lottery_strategy

import (
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
)

type PoW struct {
}

func (lottery *PoW) StartNewMiner(vk string, hardness int, initialHash sha256.HashValue, newBlockHashes chan sha256.HashValue, winningDraws chan WinningLotteryParams, stopMinerSignal chan struct{}) {
	newBlockHashesInternal := make(chan ChannelCombinationStruct)
	lottery.combineChannels(newBlockHashes, stopMinerSignal, newBlockHashesInternal)
	go lottery.startNewMinerInternal(vk, hardness, initialHash, newBlockHashesInternal, winningDraws)
}

func (lottery *PoW) combineChannels(newHashes chan sha256.HashValue, stopMiner chan struct{}, internalStruct chan ChannelCombinationStruct) {
	go func() {
		for {
			newParentHash := <-newHashes
			combination := ChannelCombinationStruct{
				minerShouldContinue: true,
				parentHash:          newParentHash,
			}
			internalStruct <- combination
		}
	}()
	go func() {
		for {
			_ = <-stopMiner
			combination := ChannelCombinationStruct{
				minerShouldContinue: false,
				parentHash:          sha256.HashValue{},
			}
			internalStruct <- combination
		}
	}()
}

func (lottery *PoW) startNewMinerInternal(vk string, hardness int, initialHash sha256.HashValue, newBlockHashesInternal chan ChannelCombinationStruct, winningDraws chan WinningLotteryParams) {
	internalStruct := ChannelCombinationStruct{
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
	return (byteToCheck >> (8 - restAmount)) == 0
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

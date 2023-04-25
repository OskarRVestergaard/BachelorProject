package lottery_strategy

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/sha256"
)

type PoW struct {
}

func (lottery *PoW) StartNewMiner(vk string, hardness int, initialHash sha256.HashValue, newBlockHashes chan sha256.HashValue, winningDraws chan WinningLotteryParams) {
	go lottery.startNewMinerInternal(vk, hardness, initialHash, newBlockHashes, winningDraws)
}

func (lottery *PoW) startNewMinerInternal(vk string, hardness int, initialHash sha256.HashValue, newBlockHashes chan sha256.HashValue, winningDraws chan WinningLotteryParams) {
	parentHash := initialHash
	for i := 0; i < 10; i++ { //TODO Make it stop mining when given a command to do so (maybe for loop running through both receiving hashes and receiving stop mining, this would be active thou)
		done := make(chan struct{})
		go lottery.mine(vk, parentHash, hardness, done, winningDraws)
		parentHash = <-newBlockHashes
		done <- struct{}{}
	}
}

func (lottery *PoW) mine(vk string, parentHash sha256.HashValue, hardness int, done chan struct{}, winningDraws chan WinningLotteryParams) {
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
			if verify(hashOfTicket, hardness) {
				winningDraws <- draw
				_ = <-done
				return
			}
		}
	}

}

func verify(hashedTicket sha256.HashValue, hardness int) bool {
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

	return verify(hashed, hardness)
}

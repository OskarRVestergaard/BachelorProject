package lottery_strategy

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/sha256"
)

type PoW struct {
}

func (lottery *PoW) StartNewMiner(vk string, hardness int, newBlockHashes chan []byte, winningDraws chan WinningLotteryParams) {
	go lottery.startNewMinerInternal(vk, hardness, newBlockHashes, winningDraws)
}

func (lottery *PoW) startNewMinerInternal(vk string, hardness int, newBlockHashes chan []byte, winningDraws chan WinningLotteryParams) {
	parentHash := <-newBlockHashes
	for {
		done := make(chan struct{})
		go lottery.mine(vk, parentHash, hardness, done, winningDraws)
		parentHash = <-newBlockHashes
		done <- struct{}{}
	}
}

func (lottery *PoW) mine(vk string, parentHash []byte, hardness int, done chan struct{}, winningDraws chan WinningLotteryParams) {
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
			hashOfTicket := sha256.HashByteArray(draw.toByteArray())
			if verify(hashOfTicket, hardness) {
				winningDraws <- draw
				return
			}
		}
	}

}

func verify(hashedTicket []byte, hardness int) bool {
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

func (lottery *PoW) Verify(vk string, parentHash []byte, hardness int, counter int) bool {
	draw := WinningLotteryParams{
		Vk:         vk,
		ParentHash: parentHash,
		Counter:    counter,
	}
	hashed := sha256.HashByteArray(draw.toByteArray())

	return verify(hashed, hardness)
}

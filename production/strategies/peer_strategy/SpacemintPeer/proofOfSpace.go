package SpacemintPeer

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Parties"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"strconv"
)

type PoSpace struct {
}

type internalCombination struct {
	minerShouldContinue bool
	parentHash          sha256.HashValue
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

func (lottery *PoSpace) StartNewMiner(PoSpacePrm Models.Parameters, vk string, hardness int, initialHash sha256.HashValue, newBlockHashes chan sha256.HashValue, winningDraws chan WinningLotteryParams, stopMinerSignal chan struct{}) (commitment []byte) {
	newBlockHashesInternal := make(chan internalCombination)
	prover := Parties.Prover{}
	prover.InitializationPhase1(PoSpacePrm)
	result := prover.GetCommitment()
	proverSingleton := make(chan Parties.Prover, 1)
	proverSingleton <- prover
	lottery.combineChannels(newBlockHashes, stopMinerSignal, newBlockHashesInternal)
	go lottery.startNewMinerInternal(proverSingleton, vk, hardness, initialHash, newBlockHashesInternal, winningDraws)
	return result
}

func (lottery *PoSpace) combineChannels(newHashes chan sha256.HashValue, stopMiner chan struct{}, internalStruct chan internalCombination) {
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

func (lottery *PoSpace) startNewMinerInternal(proverSingleton chan Parties.Prover, vk string, hardness int, initialHash sha256.HashValue, newBlockHashesInternal chan internalCombination, winningDraws chan WinningLotteryParams) {
	internalStruct := internalCombination{
		minerShouldContinue: true,
		parentHash:          initialHash,
	}
	for internalStruct.minerShouldContinue {
		done := make(chan struct{})
		go lottery.mine(proverSingleton, vk, internalStruct.parentHash, hardness, done, winningDraws)
		internalStruct = <-newBlockHashesInternal
		done <- struct{}{}
	}
}

func (lottery *PoSpace) mine(proverSingleton chan Parties.Prover, vk string, parentHash sha256.HashValue, hardness int, done chan struct{}, winningDraws chan WinningLotteryParams) {
	//TODO Fake it challenges:
	//challenges := []int{0, 1}

	//Use challenges for:
	//Calculating execution and checking quality
	//If Quality good, then do phase 2 of init (proof merkle tree and graph fit together)
	//Send winning draw
	for {
		select {
		case <-done:
			return
		default:
			//c = c + 1
			draw := WinningLotteryParams{
				Vk:         vk,
				ParentHash: parentHash,
				Counter:    0, //c,
			}
			hashOfTicket := sha256.HashByteArray(draw.ToByteSlice())
			if verifyPoSpace(hashOfTicket, hardness) {
				winningDraws <- draw
				_ = <-done
				return
			}
		}
	}

}

func verifyPoSpace(hashedTicket sha256.HashValue, hardness int) bool {
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

func (lottery *PoSpace) Verify(vk string, parentHash sha256.HashValue, hardness int, counter int) bool {
	draw := WinningLotteryParams{
		Vk:         vk,
		ParentHash: parentHash,
		Counter:    counter,
	}
	hashed := sha256.HashByteArray(draw.ToByteSlice())

	return verifyPoSpace(hashed, hardness)
}

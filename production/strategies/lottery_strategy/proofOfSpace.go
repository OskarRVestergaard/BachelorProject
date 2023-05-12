package lottery_strategy

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Parties"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"strconv"
)

type PoSpace struct {
}

type PoSpaceLotteryDraw struct {
	Vk                        string
	ParentHash                sha256.HashValue
	ProofOfSpaceA             []Models.OpeningTriple
	ProofOfCorrectCommitmentB []Models.OpeningTriple
}

func (draw PoSpaceLotteryDraw) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(draw.Vk)
	buffer.WriteString(";_;")
	buffer.Write(draw.ParentHash.ToSlice())
	buffer.WriteString(";_;")
	for _, triple := range draw.ProofOfSpaceA {
		buffer.Write(triple.Value.ToSlice())
		buffer.WriteString(";_;")
		buffer.WriteString(strconv.Itoa(triple.Index))
		buffer.WriteString(";_;")
		for _, openValue := range triple.OpenValues {
			buffer.Write(openValue.ToSlice())
			buffer.WriteString(";_;")
		}
	}
	for _, triple := range draw.ProofOfCorrectCommitmentB {
		buffer.Write(triple.Value.ToSlice())
		buffer.WriteString(";_;")
		buffer.WriteString(strconv.Itoa(triple.Index))
		buffer.WriteString(";_;")
		for _, openValue := range triple.OpenValues {
			buffer.Write(openValue.ToSlice())
			buffer.WriteString(";_;")
		}
	}
	return buffer.Bytes()
}

func (lottery *PoSpace) StartNewMiner(PoSpacePrm Models.Parameters, vk string, hardness int, initialHash sha256.HashValue, newBlockHashes chan sha256.HashValue, potentiallyWinningDraws chan PoSpaceLotteryDraw, stopMinerSignal chan struct{}) (commitment sha256.HashValue) {
	newBlockHashesInternal := make(chan ChannelCombinationStruct)
	prover := Parties.Prover{}
	prover.InitializationPhase1(PoSpacePrm)
	result := prover.GetCommitment()
	proverSingleton := make(chan Parties.Prover, 1)
	proverSingleton <- prover
	lottery.combineChannels(newBlockHashes, stopMinerSignal, newBlockHashesInternal)
	go lottery.startNewMinerInternal(proverSingleton, vk, hardness, initialHash, newBlockHashesInternal, potentiallyWinningDraws)
	return result
}

func (lottery *PoSpace) combineChannels(newHashes chan sha256.HashValue, stopMiner chan struct{}, internalStruct chan ChannelCombinationStruct) {
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

func (lottery *PoSpace) startNewMinerInternal(proverSingleton chan Parties.Prover, vk string, hardness int, initialHash sha256.HashValue, newBlockHashesInternal chan ChannelCombinationStruct, winningDraws chan PoSpaceLotteryDraw) {
	internalStruct := ChannelCombinationStruct{
		minerShouldContinue: true,
		parentHash:          initialHash,
	}
	for internalStruct.minerShouldContinue {
		go lottery.mineOnSingleBlock(proverSingleton, vk, internalStruct.parentHash, hardness, winningDraws)
		internalStruct = <-newBlockHashesInternal
	}
}

func (lottery *PoSpace) mineOnSingleBlock(proverSingleton chan Parties.Prover, vk string, parentHash sha256.HashValue, hardness int, winningDraws chan PoSpaceLotteryDraw) {
	//TODO Fake it challenges:
	challengesSetP := []int{0, 1}
	challengesSetV := []int{0, 1, 2}

	prover := <-proverSingleton
	proofOfSpaceExecution := prover.AnswerChallenges(challengesSetP, false)

	//TODO Add true quality function (currently implicily done by pathweight)
	qualityIsGoodEnough := true
	if qualityIsGoodEnough {
		proofOfCorrectCommitment := prover.AnswerChallenges(challengesSetV, true)
		draw := PoSpaceLotteryDraw{
			Vk:                        vk,
			ParentHash:                parentHash,
			ProofOfSpaceA:             proofOfSpaceExecution,
			ProofOfCorrectCommitmentB: proofOfCorrectCommitment,
		}
		winningDraws <- draw
	}
	proverSingleton <- prover
}

func (lottery *PoSpace) Verify(draw PoSpaceLotteryDraw, hardness int, commitment sha256.HashValue) bool {
	//TODO Fake it challenges:
	challengesSetP := []int{0, 1}
	challengesSetV := []int{0, 1, 2}

	verifier := Parties.Verifier{}
	verifier.SaveCommitment(commitment)
	if verifier.VerifyChallenges(challengesSetP, draw.ProofOfSpaceA, false) {
		return verifier.VerifyChallenges(challengesSetV, draw.ProofOfCorrectCommitmentB, true)
	}
	return false
}

package PoSpace

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Parties"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"math/rand"
	"strconv"
)

type PoSpace struct {
}

type MiningLocation struct {
	Slot          int
	ParentHash    sha256.HashValue
	ChallengeSetP []int
	ChallengeSetV []int
}

type channelCombination struct {
	minerShouldContinue bool
	miningLocation      MiningLocation
}

type LotteryDraw struct {
	Vk                        string
	ParentHash                sha256.HashValue
	ProofOfSpaceA             []Models.OpeningTriple
	ProofOfCorrectCommitmentB []Models.OpeningTriple
}

func (draw LotteryDraw) ToByteArray() []byte {
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

func (lottery *PoSpace) StartNewMiner(PoSpacePrm Models.Parameters, vk string, hardness int, initialMiningLocation MiningLocation, newMiningLocation chan MiningLocation, potentiallyWinningDraws chan LotteryDraw, stopMinerSignal chan struct{}) (commitment sha256.HashValue) {
	newMiningLocationsCombination := make(chan channelCombination)
	prover := Parties.Prover{}
	prover.InitializationPhase1(PoSpacePrm)
	result := prover.GetCommitment()
	proverSingleton := make(chan Parties.Prover, 1)
	proverSingleton <- prover
	lottery.combineChannels(newMiningLocation, stopMinerSignal, newMiningLocationsCombination)
	go lottery.startNewMinerInternal(proverSingleton, vk, initialMiningLocation, newMiningLocationsCombination, potentiallyWinningDraws)
	return result
}

func (lottery *PoSpace) combineChannels(newMiningLocations chan MiningLocation, stopMiner chan struct{}, miningLocationCombination chan channelCombination) {
	go func() {
		for {
			newMiningLocation := <-newMiningLocations
			combination := channelCombination{
				minerShouldContinue: true,
				miningLocation:      newMiningLocation,
			}
			miningLocationCombination <- combination
		}
	}()
	go func() {
		for {
			_ = <-stopMiner
			combination := channelCombination{
				minerShouldContinue: false,
			}
			miningLocationCombination <- combination
		}
	}()
}

func (lottery *PoSpace) startNewMinerInternal(proverSingleton chan Parties.Prover, vk string, initialMiningLocation MiningLocation, newBlockHashesInternal chan channelCombination, winningDraws chan LotteryDraw) {
	miningLocationCombination := channelCombination{
		minerShouldContinue: true,
		miningLocation:      initialMiningLocation,
	}
	for miningLocationCombination.minerShouldContinue {
		go lottery.mineOnSingleBlock(proverSingleton, vk, miningLocationCombination.miningLocation, winningDraws)
		miningLocationCombination = <-newBlockHashesInternal
	}
}

func (lottery *PoSpace) mineOnSingleBlock(proverSingleton chan Parties.Prover, vk string, miningLocation MiningLocation, winningDraws chan LotteryDraw) {
	prover := <-proverSingleton
	proofOfSpaceExecution := prover.AnswerChallenges(miningLocation.ChallengeSetP, false)

	//TODO Add true quality function
	qualityIsGoodEnough := 0.95 < rand.Float64() //TODO, THIS IS FAKE CODE TO TEST ONLY SENDING "GOOD QUALITY"
	if qualityIsGoodEnough {
		proofOfCorrectCommitment := prover.AnswerChallenges(miningLocation.ChallengeSetV, true)
		draw := LotteryDraw{
			Vk:                        vk,
			ParentHash:                miningLocation.ParentHash,
			ProofOfSpaceA:             proofOfSpaceExecution,
			ProofOfCorrectCommitmentB: proofOfCorrectCommitment,
		}
		winningDraws <- draw
	}
	proverSingleton <- prover
}

func (lottery *PoSpace) Verify(draw LotteryDraw, miningLocation MiningLocation, commitment sha256.HashValue) bool {
	verifier := Parties.Verifier{}
	verifier.SaveCommitment(commitment)
	if verifier.VerifyChallenges(miningLocation.ChallengeSetP, draw.ProofOfSpaceA, false) {
		return verifier.VerifyChallenges(miningLocation.ChallengeSetV, draw.ProofOfCorrectCommitmentB, true)
	}
	return false
}

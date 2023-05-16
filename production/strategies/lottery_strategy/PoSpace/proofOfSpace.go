package PoSpace

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Parties"
	"github.com/OskarRVestergaard/BachelorProject/Task1/PoSpaceModels"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
)

type PoSpace struct {
	n int
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
	ProofOfSpaceA             []PoSpaceModels.OpeningTriple
	ProofOfCorrectCommitmentB []PoSpaceModels.OpeningTriple
}

func (draw LotteryDraw) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(draw.Vk)
	buffer.WriteString(";_;")
	buffer.Write(draw.ParentHash.ToSlice())
	buffer.WriteString(";_;")
	buffer.Write(PoSpaceModels.ListOfTripleToByteArray(draw.ProofOfSpaceA))
	buffer.WriteString(";_;")
	buffer.Write(PoSpaceModels.ListOfTripleToByteArray(draw.ProofOfCorrectCommitmentB))
	return buffer.Bytes()
}

func (lottery *PoSpace) StartNewMiner(PoSpacePrm PoSpaceModels.Parameters, vk string, hardness int, initialMiningLocation MiningLocation, newMiningLocation chan MiningLocation, potentiallyWinningDraws chan LotteryDraw, stopMinerSignal chan struct{}) (commitment sha256.HashValue) {
	newMiningLocationsCombination := make(chan channelCombination)
	prover := Parties.Prover{}
	lottery.n = PoSpacePrm.StorageBound / 2
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
	lastSlot := -1
	lastParent := sha256.HashValue{}
	miningLocationCombination := channelCombination{
		minerShouldContinue: true,
		miningLocation:      initialMiningLocation,
	}
	for miningLocationCombination.minerShouldContinue {
		//This conditional is needed because of the peers sometimes sending duplicates, should be investigated further
		if miningLocationCombination.miningLocation.Slot != lastSlot || !miningLocationCombination.miningLocation.ParentHash.Equals(lastParent) {
			lastSlot = miningLocationCombination.miningLocation.Slot
			lastParent = miningLocationCombination.miningLocation.ParentHash
			go lottery.mineOnSingleBlock(proverSingleton, vk, miningLocationCombination.miningLocation, winningDraws)
		}
		miningLocationCombination = <-newBlockHashesInternal
	}
}

func (lottery *PoSpace) mineOnSingleBlock(proverSingleton chan Parties.Prover, vk string, miningLocation MiningLocation, winningDraws chan LotteryDraw) {
	prover := <-proverSingleton
	proofOfSpaceExecution := prover.AnswerChallenges(miningLocation.ChallengeSetP, false)

	HashOfAnswers := sha256.HashByteArray(PoSpaceModels.ListOfTripleToByteArray(proofOfSpaceExecution))
	quality := models.CalculateQuality(HashOfAnswers, int64(lottery.n))
	qualityIsGoodEnough := 0.97 < quality //TODO, One could use the best known total sum of N in network to calculate what quality would give a realistic chance to win
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

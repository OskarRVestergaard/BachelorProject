package SpacemintPeer

import (
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Parties"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
)

type PoSpace struct {
}

type internalCombination struct {
	minerShouldContinue bool
	parentHash          sha256.HashValue
}

type PoSpaceLotteryDraw struct {
	Vk                        string
	ParentHash                sha256.HashValue
	ProofOfSpaceA             []Models.OpeningTriple
	ProofOfCorrectCommitmentB []Models.OpeningTriple
}

func (lottery *PoSpace) StartNewMiner(PoSpacePrm Models.Parameters, vk string, hardness int, initialHash sha256.HashValue, newBlockHashes chan sha256.HashValue, potentiallyWinningDraws chan PoSpaceLotteryDraw, stopMinerSignal chan struct{}) (commitment sha256.HashValue) {
	newBlockHashesInternal := make(chan internalCombination)
	prover := Parties.Prover{}
	prover.InitializationPhase1(PoSpacePrm)
	result := prover.GetCommitment()
	proverSingleton := make(chan Parties.Prover, 1)
	proverSingleton <- prover
	lottery.combineChannels(newBlockHashes, stopMinerSignal, newBlockHashesInternal)
	go lottery.startNewMinerInternal(proverSingleton, vk, hardness, initialHash, newBlockHashesInternal, potentiallyWinningDraws)
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

func (lottery *PoSpace) startNewMinerInternal(proverSingleton chan Parties.Prover, vk string, hardness int, initialHash sha256.HashValue, newBlockHashesInternal chan internalCombination, winningDraws chan PoSpaceLotteryDraw) {
	internalStruct := internalCombination{
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

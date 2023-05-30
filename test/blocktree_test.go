package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/peer_strategy"
	"sync"
	"testing"
	"time"

	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/models/SpaceMintBlockchain"
	"github.com/OskarRVestergaard/BachelorProject/test/test_utils"
	"github.com/stretchr/testify/assert"
)

type TestingConstants struct {
	peerConstants      peer_strategy.PeerConstants
	noOfPeers          int
	noOfMsgs           int
	WaitBetweenMessage time.Duration
	Iterations         int
	useProofOfSpace    bool
	ProofSizeN         int
}

func testBlockChain(t *testing.T, constants TestingConstants) {
	noOfPeers := constants.noOfPeers
	noOfMsgs := constants.noOfMsgs
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfPeers, constants.useProofOfSpace, constants.peerConstants) //setup peer
	var wg sync.WaitGroup
	wg.Add(len(listOfPeers))
	for _, peer := range listOfPeers {
		go func(p peer_strategy.PeerInterface) {
			err := p.StartMining(constants.ProofSizeN)
			if err != nil {
				wg.Done()
				panic(err.Error())
			}
			wg.Done()
		}(peer)
	}
	wg.Wait()
	println("Miners finished setting up")
	starTime := time.Now()
	for _, peer := range listOfPeers {
		peer.ActivatePeer(starTime, constants.peerConstants.SlotLength)
	}
	for i := 0; i < constants.Iterations; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(constants.WaitBetweenMessage)
		return
	}
	wg.Add(len(listOfPeers))
	for _, peer := range listOfPeers {
		go func(p peer_strategy.PeerInterface) {
			err := p.StopMining()
			if err != nil {
				wg.Done()
				panic(err.Error())
			}
			wg.Done()
		}(peer)
	}
	wg.Wait()
	println("All mining activity has been stopped")
	time.Sleep(2 * constants.peerConstants.SlotLength)
	if constants.useProofOfSpace {
		for i, _ := range listOfPeers {
			if i != 0 {
				tree1 := listOfPeers[i-1].GetBlockTree().(SpaceMintBlockchain.Blocktree)
				tree2 := listOfPeers[i].GetBlockTree().(SpaceMintBlockchain.Blocktree)
				test := tree1.Equals(tree2)
				if !test {
					assert.True(t, test) //This conditional is for debugging purposes
				}
			}
		}
		tree := listOfPeers[0].GetBlockTree().(SpaceMintBlockchain.Blocktree)
		visualTree := tree.RootToVisualNode()
		chainFromHead := tree.HeadToChain()
		print(&chainFromHead)
		print(&visualTree)
	} else {
		for i, _ := range listOfPeers {
			if i != 0 {
				tree1 := listOfPeers[i-1].GetBlockTree().(PoWblockchain.Blocktree)
				tree2 := listOfPeers[i].GetBlockTree().(PoWblockchain.Blocktree)
				test := tree1.Equals(tree2)
				if !test {
					assert.True(t, test) //This conditional is for debugging purposes
				}
			}
		}
		tree := listOfPeers[0].GetBlockTree().(PoWblockchain.Blocktree)
		visualTree := tree.RootToVisualNode()
		chainFromHead := tree.HeadToChain()
		print(&chainFromHead)
		print(&visualTree)
	}
}

func TestSlow8PeerPoW(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     40,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    28,
		SlotLength:                  40000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0,
		FixedGraph:                  false,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          8,
		noOfMsgs:           1,
		WaitBetweenMessage: 10000 * time.Millisecond,
		Iterations:         60,
		useProofOfSpace:    false,
		ProofSizeN:         0,
	}
	testBlockChain(t, testingConstants)
}

func TestSlow8PeerPoS(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     40,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    26,
		SlotLength:                  40000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0.99995,
		FixedGraph:                  false,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          8,
		noOfMsgs:           1,
		WaitBetweenMessage: 10000 * time.Millisecond,
		Iterations:         60,
		useProofOfSpace:    true,
		ProofSizeN:         65536,
	}
	testBlockChain(t, testingConstants)
}

func TestFast8PeerPoW(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     40,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    22,
		SlotLength:                  6000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0,
		FixedGraph:                  false,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          8,
		noOfMsgs:           1,
		WaitBetweenMessage: 3000 * time.Millisecond,
		Iterations:         10,
		useProofOfSpace:    false,
		ProofSizeN:         0,
	}
	testBlockChain(t, testingConstants)
}

func TestFast8PeerPoS(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     40,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    28,
		SlotLength:                  10000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  true,
		ForcedD:                     20,
		QualityThreshold:            0.99995,
		FixedGraph:                  false,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          8,
		noOfMsgs:           1,
		WaitBetweenMessage: 5000 * time.Millisecond,
		Iterations:         12,
		useProofOfSpace:    true,
		ProofSizeN:         4096,
	}
	testBlockChain(t, testingConstants)
}

func TestSlow4PeerPoS(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     40,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    26,
		SlotLength:                  40000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0.99995,
		FixedGraph:                  false,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          4,
		noOfMsgs:           1,
		WaitBetweenMessage: 10000 * time.Millisecond,
		Iterations:         60,
		useProofOfSpace:    true,
		ProofSizeN:         65536,
	}
	testBlockChain(t, testingConstants)
}

func TestBiggerSlow8PeerPoS(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     40,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    26,
		SlotLength:                  40000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0.99995,
		FixedGraph:                  false,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          8,
		noOfMsgs:           1,
		WaitBetweenMessage: 10000 * time.Millisecond,
		Iterations:         60,
		useProofOfSpace:    true,
		ProofSizeN:         65536 * 2,
	}
	testBlockChain(t, testingConstants)
}

func TestSinglePeerPoS(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     40,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    26,
		SlotLength:                  40000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0.99995,
		FixedGraph:                  false,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          1,
		noOfMsgs:           1,
		WaitBetweenMessage: 10000 * time.Millisecond,
		Iterations:         8,
		useProofOfSpace:    true,
		ProofSizeN:         65536,
	}
	testBlockChain(t, testingConstants)
}

func TestSlow8PeerPoSFixedGraph(t *testing.T) {
	n := 65536
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     40,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    26,
		SlotLength:                  40000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0.99995,
		FixedGraph:                  true,
		FixedN:                      n,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          8,
		noOfMsgs:           1,
		WaitBetweenMessage: 10000 * time.Millisecond,
		Iterations:         60,
		useProofOfSpace:    true,
		ProofSizeN:         n,
	}
	testBlockChain(t, testingConstants)
}

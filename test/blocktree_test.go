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
	time.Sleep(constants.peerConstants.SlotLength)
	starTime := time.Now()
	for _, peer := range listOfPeers {
		peer.ActivatePeer(starTime, constants.peerConstants.SlotLength)
	}
	for i := 0; i < constants.Iterations; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(constants.WaitBetweenMessage)
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
		Hardness:                    30,
		SlotLength:                  40000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0,
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
		Hardness:                    28,
		SlotLength:                  40000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.0625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0.99995,
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

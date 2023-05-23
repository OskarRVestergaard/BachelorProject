package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/peer_strategy"
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
	for _, peer := range listOfPeers {
		err := peer.StartMining(constants.ProofSizeN)
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(constants.peerConstants.SlotLength)
	starTime := time.Now()
	for _, peer := range listOfPeers {
		peer.ActivatePeer(starTime, constants.peerConstants.SlotLength)
	}
	for i := 0; i < constants.Iterations; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(constants.WaitBetweenMessage)
	}
	for _, peer := range listOfPeers {
		err := peer.StopMining()
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(constants.peerConstants.SlotLength)
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

func TestPOWNetwork2Peers(t *testing.T) {
	noOfPeers := 2
	noOfMsgs := 3
	noOfNames := 2
	constants := peer_strategy.GetStandardConstants()
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, false, constants) //setup peer
	test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)                        //send msg
	for _, peer := range listOfPeers {
		err := peer.StartMining(0)
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(5000 * time.Millisecond)
	starTime := time.Now()
	for _, peer := range listOfPeers {
		peer.ActivatePeer(starTime, constants.SlotLength)
	}
	for i := 0; i < 10; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(4000 * time.Millisecond)
	}
	for _, peer := range listOfPeers {
		err := peer.StopMining()
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(7000 * time.Millisecond)
	tree1 := (listOfPeers[0].GetBlockTree()).(PoWblockchain.Blocktree)
	tree2 := listOfPeers[1].GetBlockTree().(PoWblockchain.Blocktree)
	assert.True(t, tree1.Equals(tree2))
}

func TestPOWNetwork16Peers(t *testing.T) {
	noOfPeers := 16
	noOfMsgs := 2
	noOfNames := 16
	constants := peer_strategy.GetStandardConstants()
	constants.Hardness = 30
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, false, constants) //setup peer
	test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)                        //send msg
	for _, peer := range listOfPeers {
		err := peer.StartMining(0)
		if err != nil {
			print(err.Error())
		}
	}
	for i := 0; i < 8; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(6000 * time.Millisecond)
	}
	for _, peer := range listOfPeers {
		err := peer.StopMining()
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(10000 * time.Millisecond)
	for i, _ := range listOfPeers {
		if i != 0 {
			tree1 := listOfPeers[i-1].GetBlockTree().(PoWblockchain.Blocktree)
			tree2 := listOfPeers[i].GetBlockTree().(PoWblockchain.Blocktree)
			test := tree1.Equals(tree2)
			assert.True(t, test)
		}
	}
	assert.True(t, true)
	tree1 := listOfPeers[0].GetBlockTree().(PoWblockchain.Blocktree)
	visualTree := tree1.RootToVisualNode()
	chainFromHead := tree1.HeadToChain()
	print(&chainFromHead)
	print(&visualTree)
}

func TestPoSpaceNetwork4Peers(t *testing.T) {
	noOfPeers := 4
	noOfMsgs := 4
	noOfNames := 4
	sizeOfProofsN := 32
	constants := peer_strategy.GetStandardConstants()
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, true, constants) //setup peer
	for _, peer := range listOfPeers {
		err := peer.StartMining(sizeOfProofsN)
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(5000 * time.Millisecond)
	starTime := time.Now()
	for _, peer := range listOfPeers {
		peer.ActivatePeer(starTime, constants.SlotLength)
	}
	time.Sleep(7000 * time.Millisecond) //Wait such that the peers are aware of who is mining
	for i := 0; i < 4; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(6000 * time.Millisecond)
	}
	for _, peer := range listOfPeers {
		err := peer.StopMining()
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(6000 * time.Millisecond)
	for i, _ := range listOfPeers {
		if i != 0 {
			tree1 := listOfPeers[i-1].GetBlockTree().(SpaceMintBlockchain.Blocktree)
			tree2 := listOfPeers[i].GetBlockTree().(SpaceMintBlockchain.Blocktree)
			test := tree1.Equals(tree2)
			assert.True(t, test)
		}
	}
	assert.True(t, true)
	tree1 := listOfPeers[0].GetBlockTree().(SpaceMintBlockchain.Blocktree)
	visualTree := tree1.RootToVisualNode()
	chainFromHead := tree1.HeadToChain()
	print(&chainFromHead)
	print(&visualTree)
}

func TestPoSpaceNetwork16Peers(t *testing.T) {
	noOfPeers := 16
	noOfMsgs := 2
	noOfNames := 16
	sizeOfProofsN := 512
	constants := peer_strategy.GetStandardConstants()
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, true, constants) //setup peer
	for _, peer := range listOfPeers {
		err := peer.StartMining(sizeOfProofsN)
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(20000 * time.Millisecond)
	starTime := time.Now()
	for _, peer := range listOfPeers {
		peer.ActivatePeer(starTime, constants.SlotLength)
	}
	time.Sleep(10000 * time.Millisecond) //Wait such that the peers are aware of who is mining
	for i := 0; i < 7; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(8000 * time.Millisecond)
	}
	for _, peer := range listOfPeers {
		err := peer.StopMining()
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(25000 * time.Millisecond)
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
	assert.True(t, true)
	tree1 := listOfPeers[0].GetBlockTree().(SpaceMintBlockchain.Blocktree)
	visualTree := tree1.RootToVisualNode()
	chainFromHead := tree1.HeadToChain()
	print(&chainFromHead)
	print(&visualTree)
}

func TestSlowOver15MinBig8PeerTestPoW(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     20,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    28,
		SlotLength:                  45000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          8,
		noOfMsgs:           4,
		WaitBetweenMessage: 22500 * time.Millisecond,
		Iterations:         5,
		useProofOfSpace:    false,
		ProofSizeN:         0,
	}
	testBlockChain(t, testingConstants)
}

func TestSlowOver20MinBig8PeerTestAbout1GBprPeerPoS(t *testing.T) {
	peerConstants := peer_strategy.PeerConstants{
		BlockPaymentAmountLimit:     20,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    28,
		SlotLength:                  45000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0.9999,
	}
	testingConstants := TestingConstants{
		peerConstants:      peerConstants,
		noOfPeers:          8,
		noOfMsgs:           4,
		WaitBetweenMessage: 22500 * time.Millisecond,
		Iterations:         5,
		useProofOfSpace:    true,
		ProofSizeN:         65536,
	}
	testBlockChain(t, testingConstants)
}

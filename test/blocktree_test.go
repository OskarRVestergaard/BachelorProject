package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/peer_strategy"
	"math/rand"
	"testing"
	"time"

	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/models/SpaceMintBlockchain"
	"github.com/OskarRVestergaard/BachelorProject/test/test_utils"
	"github.com/stretchr/testify/assert"
)

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
	for i := 0; i < 10; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(2000 * time.Millisecond)
	}
	for _, peer := range listOfPeers {
		err := peer.StopMining()
		if err != nil {
			print(err.Error())
		}
	}
	time.Sleep(3000 * time.Millisecond)
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
	time.Sleep(7000 * time.Millisecond) //Wait such that the peers are aware of who is mining
	for i := 0; i < 4; i++ {
		test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)
		time.Sleep(3000 * time.Millisecond)
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
	time.Sleep(15000 * time.Millisecond)
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

func TestRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(3))
	rnd2 := rand.New(rand.NewSource(3))
	println(rnd.Int())
	println(rnd2.Int())
	println(rnd.Int())
	println(rnd2.Int())
	println(rnd.Int())
	println(rnd2.Int())
	assert.True(t, true)
}

package test

import (
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
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, false) //setup peer
	test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)             //send msg
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
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, false) //setup peer
	test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)             //send msg
	for _, peer := range listOfPeers {
		err := peer.StartMining(0)
		if err != nil {
			print(err.Error())
		}
	}
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
}

func TestPoSpaceNetwork4Peers(t *testing.T) {
	noOfPeers := 4
	noOfMsgs := 4
	noOfNames := 4
	sizeOfProofsN := 8
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, true) //setup peer
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
	time.Sleep(15000 * time.Millisecond)
	for i, _ := range listOfPeers {
		if i != 0 {
			tree1 := listOfPeers[i-1].GetBlockTree().(SpaceMintBlockchain.Blocktree)
			tree2 := listOfPeers[i].GetBlockTree().(SpaceMintBlockchain.Blocktree)
			test := tree1.Equals(tree2)
			assert.True(t, test)
		}
	}
	assert.True(t, true)
}

func TestPoSpaceNetwork16Peers(t *testing.T) {
	noOfPeers := 16
	noOfMsgs := 2
	noOfNames := 16
	sizeOfProofsN := 8
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, true) //setup peer
	for _, peer := range listOfPeers {
		err := peer.StartMining(sizeOfProofsN)
		if err != nil {
			print(err.Error())
		}
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
	time.Sleep(30000 * time.Millisecond)
	for i, _ := range listOfPeers {
		if i != 0 {
			tree1 := listOfPeers[i-1].GetBlockTree().(SpaceMintBlockchain.Blocktree)
			tree2 := listOfPeers[i].GetBlockTree().(SpaceMintBlockchain.Blocktree)
			test := tree1.Equals(tree2)
			assert.True(t, test)
		}
	}
	assert.True(t, true)
}

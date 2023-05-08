package test

import (
	"github.com/OskarRVestergaard/BachelorProject/test/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPOWNetwork2Peers(t *testing.T) {
	noOfPeers := 2
	noOfMsgs := 3
	noOfNames := 2
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, false) //setup peer
	test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)             //send msg
	for _, peer := range listOfPeers {
		err := peer.StartMining()
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
	tree1 := listOfPeers[0].GetBlockTree()
	tree2 := listOfPeers[1].GetBlockTree()
	assert.True(t, tree1.Equals(tree2))
}

func TestPOWNetwork16Peers(t *testing.T) {
	noOfPeers := 16
	noOfMsgs := 2
	noOfNames := 16
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, false) //setup peer
	test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)             //send msg
	for _, peer := range listOfPeers {
		err := peer.StartMining()
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
			tree1 := listOfPeers[i-1].GetBlockTree()
			tree2 := listOfPeers[i].GetBlockTree()
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
	listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames, true) //setup peer
	test_utils.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)            //send msg
	for _, peer := range listOfPeers {
		err := peer.StartMining()
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
	time.Sleep(30000 * time.Millisecond)
	for i, _ := range listOfPeers {
		if i != 0 {
			tree1 := listOfPeers[i-1].GetBlockTree()
			tree2 := listOfPeers[i].GetBlockTree()
			test := tree1.Equals(tree2)
			assert.True(t, test)
		}
	}
	assert.True(t, true)
}

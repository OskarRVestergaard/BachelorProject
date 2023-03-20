package test

import (
	"example.com/packages/block"
	"example.com/packages/peer"
	"example.com/packages/service"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestBlockHash(t *testing.T) {
	//var transactionsLog = []*models.SignedTransaction
	//transactionsLog[1] = "a->b:100"
	//prevHash := "hejhej"
	//var b = block.MakeBlock(transactionsLog, prevHash)
	//
	//h := sha256.New()
	//
	//h.Write([]byte((prevHash + block.ConvertToString(transactionsLog))))
	//assert.Equal(t, b.Hash, string(h.Sum(nil)), "hashes match")
}

func TestBlockchainLengthOfGenesis(t *testing.T) {
	var blockChain = makeGenesisBlockchain()
	assert.Equal(t, 1, len(blockChain), "blockchain length should be 1")
}

func TestFirstConnectedPeerHasGenesisBlockslot0(t *testing.T) {
	noOfPeers := 1

	listOfPeers := make([]*peer.Peer, noOfPeers)

	var connectedPeers []string

	for i := 0; i < noOfPeers; i++ {
		var p peer.Peer
		freePort, _ := service.GetFreePort()
		port := strconv.Itoa(freePort)
		listOfPeers[i] = &p
		p.RunPeer("127.0.0.1:" + port)

	}
	listOfPeers[0].Connect("Piplup is best water pokemon", 18079)
	connectedPeers = append(connectedPeers, listOfPeers[0].IpPort)
	time.Sleep(250 * time.Millisecond)
	slotNumberGenesis := (listOfPeers[0].GenesisBlock[0].SlotNumber)
	assert.Equal(t, 0, slotNumberGenesis, "genesisblock should have slotnumber 0")

}
func Test2PeersHaveSameGenesisBlock(t *testing.T) {

	noOfPeers := 2
	noOfMsgs := 1
	noOfNames := 2
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames)             //setup peer
	controlLedger := service.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList) //send msg
	print(controlLedger)
	time.Sleep(1000 * time.Millisecond)

	time.Sleep(1000 * time.Millisecond)
	p1Genesis := listOfPeers[0].GenesisBlock[0].SlotNumber
	p2Genesis := listOfPeers[1].GenesisBlock[0].SlotNumber
	assert.Equal(t, 0, p1Genesis, "genesisblock should have slotnumber 0")
	assert.Equal(t, 0, p2Genesis, "genesisblock should have slotnumber 0")

}

func TestPeer1WinsLottery(t *testing.T) {
	noOfPeers := 2
	noOfNames := 2
	listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	pk0 := pkList[0]
	pk1 := pkList[1]
	//Action
	listOfPeers[1].FloodSignedTransaction(pk1, pk0, 50)
	time.Sleep(200 * time.Millisecond)
	// wins slot action
	
}

func makeGenesisBlockchain() map[int]*block.Block {
	genesisBlock := &block.Block{
		SlotNumber:   0,
		Hash:         "GenesisBlock",
		PreviousHash: "GenesisBlock",
		//TransactionsLog: nil,
	}
	var blockChain = make(map[int]*block.Block)
	blockChain[genesisBlock.SlotNumber] = genesisBlock
	return blockChain
}

//func Test

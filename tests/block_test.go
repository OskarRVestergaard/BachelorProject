package test

import (
	"crypto/sha256"
	"example.com/packages/block"
	"example.com/packages/peer"
	"example.com/packages/service"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestBlockHash(t *testing.T) {
	//var b block.Block
	transactionsLog := make(map[int]string)
	transactionsLog[1] = "a->b:100"
	prevHash := "hejhej"
	var b = block.MakeBlock(transactionsLog, prevHash)

	h := sha256.New()

	h.Write([]byte((prevHash + block.ConvertToString(transactionsLog))))
	assert.Equal(t, b.Hash, string(h.Sum(nil)), "hashes match")
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
	slotNumberGenesis := (listOfPeers[0].VisibleBlockchain[0].SlotNumber)
	assert.Equal(t, 0, slotNumberGenesis, "genesisblock should have slotnumber 0")

}
func Test2PeersHaveSameGenesisBlock(t *testing.T) {

}

func makeGenesisBlockchain() map[int]*block.Block {
	genesisBlock := &block.Block{
		SlotNumber:      0,
		Hash:            "GenesisBlock",
		PreviousHash:    "GenesisBlock",
		TransactionsLog: nil,
	}
	var blockChain = make(map[int]*block.Block)
	blockChain[genesisBlock.SlotNumber] = genesisBlock
	//blockChain = append(blockChain, (genesisBlock))
	return blockChain
}

//func Test

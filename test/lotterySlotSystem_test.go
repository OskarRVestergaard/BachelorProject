package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/networkservice"
	"github.com/stretchr/testify/assert"
	"math"
	"math/big"
	"testing"
	"time"
)

var G = models.Block{
	SlotNumber:   0,
	Hash:         "",
	PreviousHash: "",
	Transactions: nil,
}

func tree(ReceivedBlocks []models.Block) {
	tree := GetExampleTreeFig16()
	print(tree[0])
}
func TestTree(t *testing.T) {
	tree(nil)
}

func TestDElguess2Peers(T *testing.T) {
	var t = calculateSlotLength()
	print(t)
}

func TestHashStructure(t *testing.T) {
	g := TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	}
	node1 := TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     hashStructure(g),
		sigma: "",
	}
	assert.NotEqual(t, hashStructure(g), hashStructure(node1))
	assert.Equal(t, hashStructure(g), hashStructure(g))
}

func TestValidSlotPositive1(t *testing.T) {
	var slotLength = 3
	var slotNumber = 1

	startTime := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	oneSecAfter := startTime.Add(1 * time.Second)
	twoSecAfter := startTime.Add(2 * time.Second)

	t1 := validSlot(startTime, oneSecAfter, int64(slotNumber), int64(slotLength))
	t2 := validSlot(startTime, twoSecAfter, int64(slotNumber), int64(slotLength))
	assert.True(t, t1)
	assert.True(t, t2)
	print(oneSecAfter.Format("h"))
	//var startTime = time.Now().(time.Millisecond(1000))

}
func TestValidSlotPositive2(t *testing.T) {
	var slotLength = 3
	var slotNumber = 2

	startTime := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	fourSecAfter := startTime.Add(4 * time.Second)
	fiveSecAfter := startTime.Add(5 * time.Second)

	t1 := validSlot(startTime, fourSecAfter, int64(slotNumber), int64(slotLength))
	t2 := validSlot(startTime, fiveSecAfter, int64(slotNumber), int64(slotLength))
	assert.True(t, t1)
	assert.True(t, t2)
}
func TestValidSlotNegative(t *testing.T) {
	var slotLength = 6
	var slotNumber = 2

	startTime := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	fourSecAfter := startTime.Add(4 * time.Second)
	fiveSecAfter := startTime.Add(5 * time.Second)

	t1 := validSlot(startTime, fourSecAfter, int64(slotNumber), int64(slotLength))
	t2 := validSlot(startTime, fiveSecAfter, int64(slotNumber), int64(slotLength))
	assert.True(t, !t1)
	assert.True(t, !t2)
}

func validSlot(start time.Time, current time.Time, CurrentSlot int64, slotLength int64) bool {
	var diff = int64(current.Sub(start).Seconds()) / slotLength
	if CurrentSlot-1 <= diff && diff < CurrentSlot {
		return true
	}
	return false
}

func calculateSlotLength() int {
	noOfPeers := 2
	noOfNames := 2
	listOfPeers, pkList := networkservice.SetupPeers(noOfPeers, noOfNames) //setup peer
	pk0 := pkList[0]
	pk1 := pkList[1]

	startTime := time.Now()
	listOfPeers[1].FloodSignedTransaction(pk1, pk0, 50) //TODO should be validating block and send
	t := time.Now()

	var MaxDrift = int(math.Pow(2, -20))
	var MaxTrans = int(t.Sub(startTime).Milliseconds())
	var MaxComp = int(t.Sub(startTime).Milliseconds())
	var slotLength = 2*MaxDrift + MaxTrans + MaxComp
	return slotLength * 250 //evcery 2 seconds new slot
}

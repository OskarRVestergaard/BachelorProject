package temp2

import (
	"bytes"
	"crypto/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/networkservice"
	"math"
	"math/big"
	"sort"
	"strconv"
	"time"
)

// Definition 16.1
type TreeNode struct {
	Block string  //block
	vk    big.Int //verification kety
	slot  int     // slot number
	Draw  string  //winner ticket
	U     string  //Block data
	h     string  //block hash of some previous hash
	sigma string  //signature
}

var G = models.Block{
	SlotNumber:   0,
	Hash:         "",
	PreviousHash: "",
	Transactions: nil,
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

func PathWeight(treenodeList []TreeNode) {
	//sortedList := sortTreeNodeListBySlotNumber(treenodeList)

}

func sortTreeNodeListBySlotNumber(treenodeList []TreeNode) []TreeNode {
	sorted := treenodeList
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].slot < sorted[j].slot
	})

	return sorted
}

func HashStructure(node TreeNode) string {
	h := sha256.New()
	buffer := convertToString(node)
	h.Write(buffer.Bytes())
	return string(h.Sum(nil))
}
func convertToString(node TreeNode) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString(node.Block)
	buffer.WriteString(node.vk.String())
	buffer.WriteString(strconv.Itoa(node.slot))
	buffer.WriteString(node.Draw)
	buffer.WriteString(node.U)
	buffer.WriteString(node.h)
	buffer.WriteString(node.sigma)

	return buffer
}

package temp2

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestSortTreeNodeListBySlotNumber(t *testing.T) {
	var nodeList []TreeNode
	nodeList = append(nodeList, TreeNode{
		Block: "",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	nodeList = append(nodeList, TreeNode{
		Block: "",
		vk:    big.Int{},
		slot:  2,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	nodeList = append(nodeList, TreeNode{
		Block: "",
		vk:    big.Int{},
		slot:  1,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	nodeList = append(nodeList, TreeNode{
		Block: "",
		vk:    big.Int{},
		slot:  5,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	nodeList = append(nodeList, TreeNode{
		Block: "",
		vk:    big.Int{},
		slot:  4,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})

	var sorted = sortTreeNodeListBySlotNumber(nodeList)
	//fmt.Println(sorted)
	for i := 1; i < len(sorted); i++ {
		assert.True(t, sorted[i-1].slot < sorted[i].slot)

	}

}

func TestTreeIsTheSameAsFig16_1(t *testing.T) {
	tree := GetExampleTreeFig16()

	//tree 1
	Tree1 := tree[1]
	assert.True(t, Tree1[0].h == "")

	//tree 2
	Tree2 := tree[2]
	assert.True(t, Tree2[0].h == "")
	assert.True(t, Tree2[1].h == HashStructure(Tree2[0]))

	//tree 3
	Tree3 := tree[3]
	assert.True(t, Tree3[0].h == "")
	assert.True(t, Tree3[1].h == HashStructure(Tree3[0]))

	Tree4 := tree[4]
	assert.True(t, Tree4[0].h == "")
	assert.True(t, Tree4[1].h == HashStructure(Tree4[0]))
	assert.True(t, Tree4[2].h == HashStructure(Tree4[1]))
	assert.True(t, Tree4[3].h == HashStructure(Tree4[1]))

	Tree5 := tree[5]
	assert.True(t, Tree5[0].h == "")
	assert.True(t, Tree5[1].h == HashStructure(Tree5[0]))
	assert.True(t, Tree5[2].h == HashStructure(Tree5[1]))
	assert.True(t, Tree5[3].h == HashStructure(Tree5[1]))

	Tree6 := tree[6]
	assert.True(t, Tree6[0].h == "")
	assert.True(t, Tree6[1].h == HashStructure(Tree6[0]))
	assert.True(t, Tree6[2].h == HashStructure(Tree6[1]))
	assert.True(t, Tree6[3].h == HashStructure(Tree6[1]))
	assert.True(t, Tree6[4].h == HashStructure(Tree6[2]))
	assert.True(t, Tree6[5].h == HashStructure(Tree6[3]))

	Tree7 := tree[7]
	assert.True(t, Tree7[0].h == "")
	assert.True(t, Tree7[1].h == HashStructure(Tree7[0]))
	assert.True(t, Tree7[2].h == HashStructure(Tree7[1]))
	assert.True(t, Tree7[3].h == HashStructure(Tree7[1]))
	assert.True(t, Tree7[4].h == HashStructure(Tree7[2]))
	assert.True(t, Tree7[5].h == HashStructure(Tree7[3]))
	assert.True(t, Tree7[6].h == HashStructure(Tree7[4]))

	Tree8 := tree[8]
	assert.True(t, Tree8[0].h == "")
	assert.True(t, Tree8[1].h == HashStructure(Tree8[0]))
	assert.True(t, Tree8[2].h == HashStructure(Tree8[1]))
	assert.True(t, Tree8[3].h == HashStructure(Tree8[1]))
	assert.True(t, Tree8[4].h == HashStructure(Tree8[2]))
	assert.True(t, Tree8[5].h == HashStructure(Tree8[3]))
	assert.True(t, Tree8[6].h == HashStructure(Tree8[4]))

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
		h:     HashStructure(g),
		sigma: "",
	}
	assert.NotEqual(t, HashStructure(g), HashStructure(node1))
	assert.Equal(t, HashStructure(g), HashStructure(g))
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

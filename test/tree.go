package test

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"strconv"
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

func hashStructure(node TreeNode) string {
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

func GetExampleTreeFig16() map[int][]TreeNode {
	a := make(map[int][]TreeNode)
	a[1] = append(a[1], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	//a[1] = append(a[1], TreeNode{
	//	Block: "BLOCK",
	//	vk:    big.Int{},
	//	slot:  0,
	//	Draw:  "",
	//	U:     "",
	//	h:     "",
	//	sigma: "",
	//})
	a[2] = append(a[2], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[2] = append(a[2], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  8,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[2][0]),
		sigma: "",
	})
	a[3] = append(a[3], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[3] = append(a[3], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  8,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[3][0]),
		sigma: "",
	})
	a[4] = append(a[4], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[4] = append(a[4], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  8,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[4][0]),
		sigma: "",
	})
	a[4] = append(a[4], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[4][1]),
		sigma: "",
	})
	a[4] = append(a[4], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[4][1]),
		sigma: "",
	})

	a[5] = append(a[5], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[5] = append(a[5], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  8,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[5][0]),
		sigma: "",
	})
	a[5] = append(a[5], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[5][1]),
		sigma: "",
	})
	a[5] = append(a[5], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[5][1]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  8,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[6][0]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[6][1]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[6][1]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  14,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[6][2]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  15,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[6][3]),
		sigma: "",
	})

	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  8,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[7][0]),
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[7][1]),
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[7][2]),
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  14,
		Draw:  "",
		U:     "",
		h:     hashStructure(a[2][0]),
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  15,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  18,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})

	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  0,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  8,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  14,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  15,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  18,
		Draw:  "",
		U:     "",
		h:     "",
		sigma: "",
	})
	//print (a[0])
	return a
}

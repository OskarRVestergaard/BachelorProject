package temp2

import (
	"math/big"
)

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
		h:     HashStructure(a[2][0]),
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
		h:     HashStructure(a[3][0]),
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
		h:     HashStructure(a[4][0]),
		sigma: "",
	})
	a[4] = append(a[4], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[4][1]),
		sigma: "",
	})
	a[4] = append(a[4], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[4][1]),
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
		h:     HashStructure(a[5][0]),
		sigma: "",
	})
	a[5] = append(a[5], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[5][1]),
		sigma: "",
	})
	a[5] = append(a[5], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[5][1]),
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
		h:     HashStructure(a[6][0]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[6][1]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[6][1]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  14,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[6][2]),
		sigma: "",
	})
	a[6] = append(a[6], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  15,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[6][3]),
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
		h:     HashStructure(a[7][0]),
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[7][1]),
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[7][1]),
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  14,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[7][2]),
		sigma: "",
	})

	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  15,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[7][3]),
		sigma: "",
	})
	a[7] = append(a[7], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  18,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[7][4]),
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
		h:     HashStructure(a[8][0]),
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  11,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[8][1]),
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  12,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[8][1]),
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  14,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[8][2]),
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  15,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[8][3]),
		sigma: "",
	})
	a[8] = append(a[8], TreeNode{
		Block: "BLOCK",
		vk:    big.Int{},
		slot:  18,
		Draw:  "",
		U:     "",
		h:     HashStructure(a[8][4]),
		sigma: "",
	})

	return a
}

package main

import "example.com/packages/models"

type MerkleTree struct {
	//Binary tree, children are at (index + 1) * 2 and (index + 1) * 2 + 1
	nodes [][]byte
}

func CreateMerkleTree(graph models.Graph) {

}

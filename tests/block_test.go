package test

import (
	"crypto/sha256"
	"example.com/packages/block"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockHash(t *testing.T) {
	var b block.Block
	transactionsLog := make(map[int]string)
	transactionsLog[1] = "a->b:100"
	prevHash := "hejhej"
	b.MakeBlock(transactionsLog, prevHash)

	h := sha256.New()

	h.Write([]byte((prevHash + block.ConvertToString(transactionsLog))))
	assert.Equal(t, b.Hash, string(h.Sum(nil)), "hashes match")
}

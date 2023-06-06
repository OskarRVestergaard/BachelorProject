package PoSpaceModels

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"sort"
	"strconv"
)

// OpeningTriple All the information needed to check a specific opening.
type OpeningTriple struct {
	Index      int
	Value      sha256.HashValue
	OpenValues []sha256.HashValue //The order of the values is the order in which they are hashed during verification
}

func (triple OpeningTriple) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.Write(triple.Value.ToSlice())
	buffer.WriteString(";_;")
	buffer.WriteString(strconv.Itoa(triple.Index))
	buffer.WriteString(";_;")
	for _, openValue := range triple.OpenValues {
		buffer.Write(openValue.ToSlice())
		buffer.WriteString(";_;")
	}
	return buffer.Bytes()
}

func SortOpeningTriples(triples []OpeningTriple) []OpeningTriple {
	sort.Slice(triples, func(i, j int) bool {
		return triples[i].Index < triples[j].Index
	})
	return triples
}

func ListOfTripleToByteArray(triples []OpeningTriple) []byte {
	var buffer bytes.Buffer
	for _, triple := range triples {
		buffer.Write(triple.ToByteArray())
	}
	return buffer.Bytes()
}

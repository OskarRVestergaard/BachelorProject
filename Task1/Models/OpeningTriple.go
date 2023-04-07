package Models

// OpeningTriple All the information needed to check a specific opening.
type OpeningTriple struct {
	Index      int
	Value      []byte
	OpenValues [][]byte //The order of the values is the order in which they are hashed during verification
}

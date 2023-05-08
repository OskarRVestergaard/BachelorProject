package Models

import "github.com/OskarRVestergaard/BachelorProject/production/sha256"

// OpeningTriple All the information needed to check a specific opening.
type OpeningTriple struct {
	Index      int
	Value      sha256.HashValue
	OpenValues []sha256.HashValue //The order of the values is the order in which they are hashed during verification
}

package models

import "math/big"

type SignedTransaction struct {
	From      string
	To        string
	Amount    int
	Signature *big.Int
}

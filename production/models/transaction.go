package models

import (
	"github.com/google/uuid"
	"math/big"
)

type SignedTransaction struct {
	Id        uuid.UUID
	From      string
	To        string
	Amount    int
	Signature *big.Int
}

package blockchain

import (
	"github.com/google/uuid"
)

type SignedTransaction struct {
	Id        uuid.UUID
	From      string
	To        string
	Amount    int
	Signature []byte
}

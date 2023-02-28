package signature_strategy

import (
	"example.com/packages/models"
	"math/big"
)

type SignatureInterface interface {
	KeyGen() (*big.Int, *big.Int, *big.Int)
	Sign(models.SignedTransaction, string) *big.Int
	Verify(models.SignedTransaction) bool
}

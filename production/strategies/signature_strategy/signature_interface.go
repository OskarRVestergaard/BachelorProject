package signature_strategy

import (
	"math/big"
)

type SignatureInterface interface {
	KeyGen() (string, string)
	Sign([]byte, string) *big.Int
	Verify(string, []byte, *big.Int) bool
}
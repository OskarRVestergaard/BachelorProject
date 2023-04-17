package hash_strategy

import (
	"crypto/sha256"
	"math/big"
)

func Hash_SHA256(msg string) *big.Int {

	h := sha256.New()
	h.Write([]byte(msg))

	hm := new(big.Int).SetBytes(h.Sum(nil))

	return hm
}

func HashByteArray(toBeHashed []byte) []byte {
	h := sha256.New()
	h.Write(toBeHashed)

	return h.Sum(nil)
}

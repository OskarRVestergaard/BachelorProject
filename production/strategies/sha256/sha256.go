package sha256

import (
	"crypto/sha256"
)

func HashByteArray(toBeHashed []byte) []byte {
	h := sha256.New()
	h.Write(toBeHashed)

	return h.Sum(nil)
}

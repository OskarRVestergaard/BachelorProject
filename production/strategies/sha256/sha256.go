package sha256

import (
	"crypto/sha256"
)

type HashValue [32]byte

func HashByteArrayToByteArray(toBeHashed []byte) []byte {
	h := sha256.New()
	h.Write(toBeHashed)

	return h.Sum(nil)
}

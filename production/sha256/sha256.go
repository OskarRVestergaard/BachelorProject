package sha256

import (
	"bytes"
	"crypto/sha256"
	"log"
)

type HashValue [32]byte

func HashByteArrayToByteArray(toBeHashed []byte) []byte {
	h := sha256.New()
	h.Write(toBeHashed)
	return h.Sum(nil)
}

func HashByteArray(toBeHashed []byte) HashValue {
	h := sha256.New()
	h.Write(toBeHashed)

	return sliceToHash(h.Sum(nil))
}

func ToString(hash HashValue) string {
	return string(hash[:])
}

func ToSlice(hash HashValue) []byte {
	slice := hash[:]
	return slice
}

func sliceToHash(bytes []byte) HashValue {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Tried to convert a byte slice to a hash-value but failed, probably because the slice had the wrong size!")
		}
	}()
	s4 := (*HashValue)(bytes)
	return *s4
}

func HashesEqual(hash1 HashValue, hash2 HashValue) bool {
	return bytes.Equal(ToSlice(hash1), ToSlice(hash2))
}

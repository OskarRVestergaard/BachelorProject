package sha256

import (
	"bytes"
	"crypto/sha256"
	"log"
)

type HashValue [32]byte

func HashByteArray(toBeHashed []byte) HashValue {
	h := sha256.New()
	h.Write(toBeHashed)

	return sliceToHash(h.Sum(nil))
}

func (hash HashValue) ToString() string {
	return string(hash[:])
}

func (hash HashValue) ToSlice() []byte {
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

func (hash HashValue) Equals(comparisonHash HashValue) bool {
	return bytes.Equal(hash.ToSlice(), comparisonHash.ToSlice())
}

package hash_strategy

import (
	"crypto/sha256"
	"example.com/packages/models"
	"math/big"
	"strconv"
	"strings"
)

func Hash_SHA256(msg string) *big.Int {

	h := sha256.New()
	h.Write([]byte(msg))

	hm := new(big.Int).SetBytes(h.Sum(nil))

	return hm
}

func HashSignedTransactionToByteArrayWowSoCool(transaction models.SignedTransaction) []byte {
	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
	t = strings.Replace(t, ";", "", -1)

	h := sha256.New()
	h.Write([]byte(t))

	hm := h.Sum(nil)

	//hash := sha256.Sum256([]byte(t))

	return hm

}

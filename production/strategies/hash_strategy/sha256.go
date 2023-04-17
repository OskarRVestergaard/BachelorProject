package hash_strategy

import (
	"crypto/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/models/messages"
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

func HashByteArray(toBeHashed []byte) []byte {
	h := sha256.New()
	h.Write(toBeHashed)

	return h.Sum(nil)
}

func HashSignedTransactionToByteArrayWowSoCool(transaction messages.SignedTransaction) []byte {
	t := transaction.From + transaction.To + strconv.Itoa(transaction.Amount)
	t = strings.Replace(t, ";", "", -1)
	//Todo This is NOT!!!!! safe, different predictable signed transactions use the same signature! - Daniel
	//TODO Also the ID should be used right?
	h := sha256.New()
	h.Write([]byte(t))

	hm := h.Sum(nil)

	//hash := sha256.Sum256([]byte(t))

	return hm

}

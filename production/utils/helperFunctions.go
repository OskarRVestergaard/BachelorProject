package utils

import (
	"crypto/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"math/big"
	"strconv"
	"strings"
)

func GetSomeKey[t comparable](m map[t]t) t {
	for k := range m {
		return k
	}
	panic("Cant get key from an empty map!")
}

func ConvertStringToBigInt(str string) *big.Int {
	result := big.NewInt(0)
	result, wasSuccessful := result.SetString(str, 10)
	if wasSuccessful {
		return result
	}
	panic("Unable to convert string to bigint: " + str)
}

func TransactionHasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface, signedTrans blockchain.SignedTransaction) bool {
	hashedMessage := HashSignedTransactionToByteArrayWowSoCool(signedTrans)
	publicKey := signedTrans.From
	signature := signedTrans.Signature
	return signatureStrategy.Verify(publicKey, hashedMessage, &signature)
}

func HashSignedTransactionToByteArrayWowSoCool(transaction blockchain.SignedTransaction) []byte {
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

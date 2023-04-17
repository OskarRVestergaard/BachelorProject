package utils

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/messages"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"math/big"
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

func TransactionHasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface, signedTrans messages.SignedTransaction) bool {
	hashedMessage := hash_strategy.HashSignedTransactionToByteArrayWowSoCool(signedTrans)
	publicKey := signedTrans.From
	signature := signedTrans.Signature
	return signatureStrategy.Verify(publicKey, hashedMessage, signature)
}

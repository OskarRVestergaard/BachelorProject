package utils

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
)

func GetSomeKey[t comparable](m map[t]t) t {
	for k := range m {
		return k
	}
	panic("Cant get key from an empty map!")
}

func TransactionHasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface, signedTrans blockchain.SignedTransaction) bool {
	transByteArray := signedTrans.ToByteArrayWithoutSign()
	hashedMessage := hash_strategy.HashByteArray(transByteArray)
	publicKey := signedTrans.From
	signature := signedTrans.Signature
	return signatureStrategy.Verify(publicKey, hashedMessage, signature)
}

func MakeDeepCopyOfTransaction(transaction blockchain.SignedTransaction) (copyOfTransaction blockchain.SignedTransaction) {
	oldSign := transaction.Signature
	signatureCopy := make([]byte, len(oldSign))
	copy(signatureCopy, oldSign)
	deepCopyTransaction := blockchain.SignedTransaction{
		Id:        transaction.Id,
		From:      transaction.From,
		To:        transaction.To,
		Amount:    transaction.Amount,
		Signature: signatureCopy,
	}
	return deepCopyTransaction
}

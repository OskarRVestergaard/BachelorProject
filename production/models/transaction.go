package models

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"github.com/google/uuid"
	"strconv"
)

type SignedTransaction struct {
	Id        uuid.UUID
	From      string
	To        string
	Amount    int
	Signature []byte
}

func (signedTransaction *SignedTransaction) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.Write(signedTransaction.ToByteArrayWithoutSign())
	buffer.WriteString(";_;")
	buffer.Write(signedTransaction.Signature)
	return buffer.Bytes()
}

func (signedTransaction *SignedTransaction) ToByteArrayWithoutSign() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(signedTransaction.Id.String())
	buffer.WriteString(";_;")
	buffer.WriteString(signedTransaction.From)
	buffer.WriteString(";_;")
	buffer.WriteString(signedTransaction.To)
	buffer.WriteString(";_;")
	buffer.WriteString(strconv.Itoa(signedTransaction.Amount))
	return buffer.Bytes()
}

func (signedTransaction *SignedTransaction) SignTransaction(signatureStrategy signature_strategy.SignatureInterface, secretSigningKey string) {
	byteArrayTransaction := signedTransaction.ToByteArrayWithoutSign()
	hashedTransaction := sha256.HashByteArray(byteArrayTransaction).ToSlice()
	signature := signatureStrategy.Sign(hashedTransaction, secretSigningKey)
	signedTransaction.Signature = signature
}

func GetTransactionsInList1ButNotList2(list1 []SignedTransaction, list2 []SignedTransaction) []SignedTransaction {
	//Currently, since the lists are unsorted the algorithm just loops over all nm combinations, could be sorted first and then i would run in nlogn+mlogm
	var difference []SignedTransaction
	found := false
	for _, val1 := range list1 {
		found = false
		for _, val2 := range list2 {
			if val1.Id == val2.Id {
				found = true
			}
		}
		if !found {
			difference = append(difference, val1)
		}
	}

	return difference
}

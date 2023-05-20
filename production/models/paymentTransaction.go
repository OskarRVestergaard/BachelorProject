package models

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"github.com/google/uuid"
	"strconv"
)

type SignedPaymentTransaction struct {
	Id        uuid.UUID
	From      string
	To        string
	Amount    int
	Signature []byte
}

func (signedPaymentTransaction *SignedPaymentTransaction) ToByteArray() []byte {
	var buffer bytes.Buffer
	buffer.Write(signedPaymentTransaction.ToByteArrayWithoutSign())
	buffer.WriteString(";_;")
	buffer.Write(signedPaymentTransaction.Signature)
	return buffer.Bytes()
}

func (signedPaymentTransaction *SignedPaymentTransaction) ToByteArrayWithoutSign() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(signedPaymentTransaction.Id.String())
	buffer.WriteString(";_;")
	buffer.WriteString(signedPaymentTransaction.From)
	buffer.WriteString(";_;")
	buffer.WriteString(signedPaymentTransaction.To)
	buffer.WriteString(";_;")
	buffer.WriteString(strconv.Itoa(signedPaymentTransaction.Amount))
	return buffer.Bytes()
}

func (signedPaymentTransaction *SignedPaymentTransaction) SignTransaction(signatureStrategy signature_strategy.SignatureInterface, secretSigningKey string) {
	byteArrayTransaction := signedPaymentTransaction.ToByteArrayWithoutSign()
	hashedTransaction := sha256.HashByteArray(byteArrayTransaction).ToSlice()
	signature := signatureStrategy.Sign(hashedTransaction, secretSigningKey)
	signedPaymentTransaction.Signature = signature
}

func GetTransactionsInList1ButNotList2(list1 []SignedPaymentTransaction, list2 []SignedPaymentTransaction) []SignedPaymentTransaction {
	//Currently, since the lists are unsorted the algorithm just loops over all nm combinations, could be sorted first and then i would run in nlogn+mlogm
	var difference []SignedPaymentTransaction
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

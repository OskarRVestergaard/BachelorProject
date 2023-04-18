package blockchain

import (
	"bytes"
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

// TODO ADD SOME SORT OF SEPERATOR BETWEEN THEM, SINCE THIS IS ONLY ONE WAY, AND CAN BE EXPLOITED
func (signedTransaction *SignedTransaction) ToByteArray() []byte {
	var firstBytes = signedTransaction.ToByteArrayWithoutSign()
	firstBytes = append(firstBytes, signedTransaction.Signature...)
	return firstBytes
}

// TODO ADD SOME SORT OF SEPERATOR BETWEEN THEM, SINCE THIS IS ONLY ONE WAY, AND CAN BE EXPLOITED
func (signedTransaction *SignedTransaction) ToByteArrayWithoutSign() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(signedTransaction.Id.String())
	buffer.WriteString(signedTransaction.From)
	buffer.WriteString(signedTransaction.To)
	buffer.WriteString(strconv.Itoa(signedTransaction.Amount))

	return buffer.Bytes()
}

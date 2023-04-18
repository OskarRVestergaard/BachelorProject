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

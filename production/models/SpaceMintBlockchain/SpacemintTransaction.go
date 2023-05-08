package PoWblockchain

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/google/uuid"
)

type SpacemintTransactions struct {
	payments         []models.SignedTransaction
	SpaceCommitments []SpaceCommitment
	Penalties        []Penalty
}

func (transactions *SpacemintTransactions) ToByteArray() []byte {
	var buffer bytes.Buffer
	for _, payment := range transactions.payments {
		buffer.Write(payment.ToByteArray())
		buffer.WriteString(";_;")
	}
	for _, spaceCommitment := range transactions.SpaceCommitments {
		buffer.WriteString(spaceCommitment.Id.String())
		buffer.WriteString(";_;")
		buffer.WriteString(spaceCommitment.PublicKey)
		buffer.WriteString(";_;")
		buffer.Write(spaceCommitment.Commitment.ToSlice())
		buffer.WriteString(";_;")
	}
	for _, penalty := range transactions.Penalties {
		buffer.WriteString(penalty.Id.String())
		buffer.WriteString(";_;")
		buffer.WriteString(penalty.PublicKey)
		buffer.WriteString(";_;")
		buffer.WriteString(penalty.Proof)
		buffer.WriteString(";_;")
	}
	return buffer.Bytes()
}

type SpaceCommitment struct {
	Id         uuid.UUID
	PublicKey  string
	Commitment sha256.HashValue
}

type Penalty struct {
	Id        uuid.UUID
	PublicKey string
	Proof     string //TODO Change to other type, to reflect these proofs of dishonest behaviour
}

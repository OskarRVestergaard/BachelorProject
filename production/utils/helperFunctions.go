package utils

import (
	"time"

	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
)

func GetSomeKey[t comparable](m map[t]t) t {
	for k := range m {
		return k
	}
	panic("Cant get key from an empty map!")
}

func TransactionHasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface, signedTrans models.SignedTransaction) bool {
	transByteArray := signedTrans.ToByteArrayWithoutSign()
	hashedMessage := sha256.HashByteArray(transByteArray).ToSlice()
	publicKey := signedTrans.From
	signature := signedTrans.Signature
	result := signatureStrategy.Verify(publicKey, hashedMessage, signature)
	return result
}

func CalculateSlot(startTime time.Time) int {
	timeDifference := time.Now().Sub(startTime)
	slot := timeDifference.Milliseconds() / constants.SlotLength.Milliseconds()
	return int(slot)
}

/*
startTimeSlotUpdater returns a channel that reports when a new time slot has started and what the time slot is
*/
func StartTimeSlotUpdater(startTime time.Time) chan int {
	updater := make(chan int)
	prevSlot := 0
	go func() {
		for {
			currentSlot := CalculateSlot(startTime)
			if currentSlot > prevSlot {
				updater <- currentSlot
			}
			time.Sleep(constants.SlotLength / 10)
		}
	}()
	return updater
}

func PowerOfTwo(n int) bool {
	i := 1
	for i < n {
		i = i * 2
	}
	return i == n
}

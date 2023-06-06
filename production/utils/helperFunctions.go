package utils

import (
	"time"

	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
)

func GetSomeKey[t comparable](m map[t]t) t {
	for k := range m {
		return k
	}
	panic("Cant get key from an empty map!")
}

func TransactionHasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface, signedTrans models.SignedPaymentTransaction) bool {
	transByteArray := signedTrans.ToByteArrayWithoutSign()
	hashedMessage := sha256.HashByteArray(transByteArray).ToSlice()
	publicKey := signedTrans.From
	signature := signedTrans.Signature
	result := signatureStrategy.Verify(publicKey, hashedMessage, signature)
	return result
}

func CalculateSlot(startTime time.Time, slotLength time.Duration) int {
	timeDifference := time.Now().Sub(startTime)
	slot := timeDifference.Milliseconds() / slotLength.Milliseconds()
	return int(slot)
}

/*
startTimeSlotUpdater returns a channel that reports when a new time slot has started and what the time slot is
*/
func StartTimeSlotUpdater(startTime time.Time, slotLength time.Duration, slotNotifier chan int) {
	prevSlot := 0
	go func() {
		for {
			currentSlot := CalculateSlot(startTime, slotLength)
			if currentSlot > prevSlot {
				slotNotifier <- currentSlot
				prevSlot = currentSlot
			}
			time.Sleep(slotLength / 10)
		}
	}()
}

func PowerOfTwo(n int) bool {
	i := 1
	for i < n {
		i = i * 2
	}
	return i == n
}

package peer_strategy

import "time"

type PeerConstants struct {
	BlockPaymentAmountLimit     int
	BlockSpaceCommitAmountLimit int
	BlockPenaltyAmountLimit     int
	Hardness                    int
	SlotLength                  time.Duration
	GraphK                      int
	Alpha                       float64
	Beta                        float64
	UseForcedD                  bool
	ForcedD                     int
	QualityThreshold            float64
}

func GetStandardConstants() PeerConstants {
	return PeerConstants{
		BlockPaymentAmountLimit:     20,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    23,
		SlotLength:                  30000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.625,
		Beta:                        0.925,
		UseForcedD:                  false,
		ForcedD:                     0,
		QualityThreshold:            0.999,
	}
}

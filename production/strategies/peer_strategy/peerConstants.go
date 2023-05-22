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
}

func GetStandardConstants() PeerConstants {
	return PeerConstants{
		BlockPaymentAmountLimit:     20,
		BlockSpaceCommitAmountLimit: 32,
		BlockPenaltyAmountLimit:     32,
		Hardness:                    23,
		SlotLength:                  5000 * time.Millisecond,
		GraphK:                      128,
		Alpha:                       0.25,
		Beta:                        0.5,
		UseForcedD:                  false,
		ForcedD:                     0,
	}
}

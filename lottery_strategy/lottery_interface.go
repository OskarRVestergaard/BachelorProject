package lottery_strategy

import "math/big"

type LotteryInterface interface {
	Mine(string, string) (bool, *big.Int)
	Verify(string) bool
}

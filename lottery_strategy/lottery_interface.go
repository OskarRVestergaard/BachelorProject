package lottery_strategy

type LotteryInterface interface {
	Mine(string) string
	Verify(string) bool
}

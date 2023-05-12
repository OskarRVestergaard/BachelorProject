package lottery_strategy

import "github.com/OskarRVestergaard/BachelorProject/production/sha256"

type ChannelCombinationStruct struct {
	minerShouldContinue bool
	parentHash          sha256.HashValue
}

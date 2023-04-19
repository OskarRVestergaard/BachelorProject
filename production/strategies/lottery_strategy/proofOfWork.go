package lottery_strategy

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/sha256"
	"math/big"
	"strconv"
)

type PoW struct {
}

var hardness = 2

func (lottery PoW) Mine(vk string, aBlockToExtend []byte) (bool, *big.Int) {

	vkByte := []byte(vk)
	aAndVk := append(vkByte, aBlockToExtend...)
	loopCon := 1000

	for c := 1; c < loopCon; c++ {
		tempToHash := append(aAndVk, strconv.Itoa(c)...) //TODO Change so aAndVk not a problem, and not string counting (prob not important)
		hash := sha256.HashByteArray(tempToHash)
		hm := new(big.Int).SetBytes(hash)
		originalHm := big.NewInt(0)
		originalHm.Set(hm)
		hm = hm.Rsh(hm, uint(256-hardness))

		zero := big.NewInt(0)

		if zero.Cmp(hm) == 0 {
			println("I won :)")
			return true, originalHm
		}
	}

	return false, big.NewInt(0)
}

func (lottery PoW) Verify(block string) bool {
	loopCon := 100
	for testHash := 0; testHash < loopCon; testHash++ {
		println(hardness)
	}

	return true
}

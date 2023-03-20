package lottery_strategy

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"strconv"
)

type PoW struct {
}

var hardness = 2

func (lottery PoW) Mine(vk string, a string) (bool, *big.Int) {
	//c := 1
	aByte := []byte(a)
	vkByte := []byte(vk)
	H := sha256.New()
	loopCon := 1000

	for c := 1; c < loopCon; c++ {

		temp := [][]byte{aByte, []byte(strconv.Itoa(c)), vkByte}
		H.Write(bytes.Join(temp, []byte("")))
		hm := new(big.Int).SetBytes(H.Sum(nil))
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

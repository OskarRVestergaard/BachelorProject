package lottery_strategy

import (
	"bytes"
	"crypto/sha256"
	"strconv"
)

type PoW struct {
}

var hardness = 1

func (lottery PoW) Mine(vk string, a string) (bool, []byte) {
	//c := 1
	aByte := []byte(a)
	vkByte := []byte(vk)
	H := sha256.New()
	loopCon := 1000

	zeroHard := make([]byte, hardness)
	for c := 1; c < loopCon; c++ {
		temp := [][]byte{aByte, []byte(strconv.Itoa(c)), vkByte}
		H.Write(bytes.Join(temp, []byte("")))
		hm := H.Sum(nil)
		if bytes.HasPrefix(hm, zeroHard) {
			return true, hm
		}

	}

	return false, []byte("1")
}

func (lottery PoW) Verify(block string) bool {
	hardness := 10
	loopCon := 100
	for testHash := 0; testHash < loopCon; testHash++ {
		println(hardness)
	}

	return true
}

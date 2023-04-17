package utils

import "math/big"

func GetSomeKey[t comparable](m map[t]t) t {
	for k := range m {
		return k
	}
	panic("Cant get key from an empty map!")
}

func ConvertStringToBigInt(str string) *big.Int {
	result := big.NewInt(0)
	result, wasSuccessful := result.SetString(str, 10)
	if wasSuccessful {
		return result
	}
	panic("Unable to convert string to bigint: " + str)
}

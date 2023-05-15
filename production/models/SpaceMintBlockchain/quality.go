package SpaceMintBlockchain

import (
	"math"
	"math/big"
)

func CalculateQuality(block Block) (quality float64) {
	//Floating point quality might not be optimal, but should be good enough for our purposes

	//All is done under the assumption that a hash has size 32 bytes and there is lots of unneeded allocation
	hashValue := &big.Int{}
	hashValue = hashValue.SetBytes(block.HashOfBlock().ToSlice()) //TODO Should not be the hash of the block, but the hash of the challenges a_i
	hashValueFloat := &big.Float{}
	hashValueFloat = hashValueFloat.SetInt(hashValue)

	//Max Uint64
	maxInt := &big.Int{}
	maxInt = maxInt.SetUint64(math.MaxUint64)

	//2^64
	maxIntPlusOne := &big.Int{}
	maxIntPlusOne = maxIntPlusOne.Add(maxInt, big.NewInt(1))

	//2^256
	maxHashValueSize := &big.Int{}
	maxHashValueSize = maxHashValueSize.Exp(maxIntPlusOne, big.NewInt(4), nil)
	maxHashValueSizeFloat := &big.Float{}
	maxHashValueSizeFloat = maxHashValueSizeFloat.SetInt(maxHashValueSize)

	//Normalized hash
	normalizedHash := &big.Float{}
	normalizedHash = normalizedHash.Quo(hashValueFloat, maxHashValueSizeFloat)

	//TODO REMOVE
	var n int64
	n = 8

	//Exponent fraction
	numerator := &big.Float{}
	numerator = numerator.SetInt(big.NewInt(1))
	denominator := &big.Float{}
	denominator = denominator.SetInt(big.NewInt(n))

	exponent := &big.Float{}
	exponent = exponent.Quo(numerator, denominator)

	//Final quality
	finalQuality := pow(normalizedHash, exponent, big.NewFloat(.0000000000000000000000000000001))

	//Conversion to float64
	qualityResult, _ := finalQuality.Float64()
	if qualityResult < 0 || qualityResult > 1 {
		panic("Problem during quality calculation")
	}
	return qualityResult
}

func CalculateChainQuality(singleBlockQualitiesFromHeadToGenesis []float64) (chainQuality float64) {
	//TODO Implement, current just fake it code
	return singleBlockQualitiesFromHeadToGenesis[0]
}

func mySqr(a *big.Float) *big.Float {
	return big.NewFloat(0).Mul(a, a)
}

func pow(a *big.Float, b *big.Float, precision *big.Float) *big.Float {
	if b.Cmp(big.NewFloat(0)) == -1 {
		return big.NewFloat(1).Quo(big.NewFloat(1), pow(a, big.NewFloat(0).Neg(b), precision))
	}
	if b.Cmp(big.NewFloat(10)) != -1 {
		return mySqr(pow(a, big.NewFloat(0).Quo(b, big.NewFloat(2)), big.NewFloat(0).Quo(precision, big.NewFloat(2))))
	}
	if b.Cmp(big.NewFloat(1)) != -1 {
		return big.NewFloat(0).Mul(a, pow(a, big.NewFloat(0).Sub(b, big.NewFloat(1)), precision))
	}
	if precision.Cmp(big.NewFloat(1)) != -1 {
		return big.NewFloat(0).Sqrt(a)
	}
	return big.NewFloat(0).Sqrt(pow(a, big.NewFloat(0).Mul(b, big.NewFloat(2)), big.NewFloat(0).Mul(precision, big.NewFloat(2))))
}

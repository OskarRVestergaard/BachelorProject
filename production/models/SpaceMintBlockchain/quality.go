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

	//Exponent fraction
	//TODO Should have access to N, which can only be done after space commit transactions have been implemented fully, here we just assume n = 8 BUT should be changed!
	n := int64(8)
	numerator := &big.Float{}
	numerator = numerator.SetInt(big.NewInt(1))
	denominator := &big.Float{}
	denominator = denominator.SetInt(big.NewInt(n))

	exponent := &big.Float{}
	exponent = exponent.Quo(numerator, denominator)

	//Final quality
	//TODO Figure out how to do power opeartion with big floats

	//Conversion to float64
	qualityResult, _ := normalizedHash.Float64() //Todo should be final quality instead ofc
	if qualityResult < 0 || qualityResult > 1 {
		panic("Problem during quality calculation")
	}
	return qualityResult
}

func CalculateChainQuality(singleBlockQualitiesFromHeadToGenesis []float64) (chainQuality float64) {
	//TODO Implement, current just fake it code
	return singleBlockQualitiesFromHeadToGenesis[0]
}

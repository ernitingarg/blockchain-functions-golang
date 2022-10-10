package helpers

import "math/big"

// FromSatoshiToBtc convert a value in satoshi (int) to a value in btc (float)
func FromSatoshiToBtc(i *big.Int) (f float64) {
	flVal := new(big.Float).SetInt(i)
	f, _ = flVal.Mul(flVal, big.NewFloat(10e-9)).Float64()
	return
}

// FromBtcToSatoshi convert a value in btc (float) to a value in satoshi (int)
func FromBtcToSatoshi(f *big.Float) *big.Int {
	i, _ := f.Mul(f, big.NewFloat(10e9)).Int64()
	return big.NewInt(i)
}

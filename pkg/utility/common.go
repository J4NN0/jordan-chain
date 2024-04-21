package utility

import (
	"math"
	"math/big"
	"regexp"
)

func WeiToETh(wei *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(wei.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue
}

func IsAddressValid(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}

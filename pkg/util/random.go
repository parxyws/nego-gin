package util

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func GenerateRandomInteger() (string, error) {

	randomNum, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", randomNum.Int64()), nil
}

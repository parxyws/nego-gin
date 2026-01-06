package util

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func GenerateRandomInteger() (string, error) {

	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", n.Int64()), nil
}

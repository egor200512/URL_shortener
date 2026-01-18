package links

import (
	"crypto/rand"
	"math/big"
)

const (
	CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	LENGTH  = 6
)

func GenerateShortLink() (string, error) {
	result := make([]byte, LENGTH)

	for i := range LENGTH {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(CHARSET))))
		if err != nil {
			return "", err
		}
		result[i] = CHARSET[num.Int64()]
	}

	return string(result), nil
}

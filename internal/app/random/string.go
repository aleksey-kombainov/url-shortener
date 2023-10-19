package random

import (
	"math/rand"
)

const alpha = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"

func GenString(length int) string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(buf)
}

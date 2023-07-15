package util

import (
	"math/rand"
	"strings"
)

const numbers = "123456789"

func RandomStringNumber(lengthOfString uint8) string {
	var sb strings.Builder

	len := len(numbers)

	var i uint8 = 0

	for ; i < lengthOfString; i++ {
		c := numbers[uint8(rand.Int31n(int32(len)))]
		sb.WriteByte(c)
	}

	return sb.String()
}

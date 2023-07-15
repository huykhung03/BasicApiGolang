package util

import (
	"math/rand"
	"strings"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnm"

func RandomString(lengthOfString uint8) string {
	var sb strings.Builder

	len := len(alphabet)

	var i uint8 = 0

	for ; i < lengthOfString; i++ {
		c := alphabet[uint8(rand.Int31n(int32(len)))]
		sb.WriteByte(c)
	}

	return sb.String()
}

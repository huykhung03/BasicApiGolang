package util

import (
	"math/rand"
	"strings"
)

// RandIntNumber generates an unsinged integer number between min -> max
func RandomIntNumber(min, max uint32) uint32 {
	return min + uint32(rand.Int31n(int32(max-min+1)))
}

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

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "YEN", "VND"}
	len := len(currencies)
	return currencies[rand.Int31n(int32(len))]
}

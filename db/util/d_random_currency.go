package util

import "math/rand"

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "YEN", "VND"}
	len := len(currencies)
	return currencies[rand.Int31n(int32(len))]
}

package util

import "math/rand"

// RandIntNumber generates an unsinged integer number between min -> max
func RandomIntNumber(min, max uint32) uint32 {
	return min + uint32(rand.Int31n(int32(max-min+1)))
}

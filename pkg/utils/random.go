package utils

import (
	"math/rand"
	"time"
)

func RandomIntN(n int) int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	return random.Intn(n)
}

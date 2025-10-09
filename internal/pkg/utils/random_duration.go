package utils

import (
	"math/rand"
	"time"
)

func RandomDuration(min, max int) time.Duration {
	return time.Duration(min+rand.Intn(max-min+1)) * time.Second
}

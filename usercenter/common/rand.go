package common

import (
	"math/rand"
	"time"
)

func RandString(length int) string {
	if length == 0 {
		length = 16
	}

	sources := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"!@#$%^&*()_+"
	rand.NewSource(time.Now().UnixNano())

	buf := make([]byte, length)

	for i := range buf {
		buf[i] = sources[rand.Intn(len(sources))]
	}
	return string(buf)
}

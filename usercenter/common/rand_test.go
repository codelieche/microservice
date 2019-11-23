package common

import (
	"log"
	"testing"
)

func TestRandString(t *testing.T) {
	i := 0
	for i < 10 {
		i++
		value := RandString(32)
		log.Println(value)
	}
}

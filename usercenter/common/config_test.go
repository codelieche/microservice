package common

import (
	"log"
	"testing"
)

func TestGetConfig(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := ParseConfig(); err != nil {
		t.Error(err.Error())
	}
	c := GetConfig()
	log.Println(c)
}

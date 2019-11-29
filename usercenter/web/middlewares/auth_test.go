package middlewares

import (
	"log"
	"testing"
)

func TestCheckSessionIsOK(t *testing.T) {

	sessionID := "117ec402-c475-4a70-a3f6-2b9eeea0add1"
	userID := 1
	log.Println(CheckSessionIsOK(sessionID, userID))
}

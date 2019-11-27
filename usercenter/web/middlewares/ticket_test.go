package middlewares

import "testing"

func TestCheckTicketFromSSOServer(t *testing.T) {
	ticket := "0bc109722b81f09434ef9f53b6e082d9"
	CheckTicketFromSSOServer(ticket)
}

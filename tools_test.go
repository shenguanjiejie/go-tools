package tools

import "testing"

func TestInternet(t *testing.T) {
	pass := InternetCheck()
	Logln(pass)
}

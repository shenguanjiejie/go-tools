package tools

import (
	"testing"
	"time"
)

func TestInternet(t *testing.T) {

	go func() {
		for range time.Tick(time.Second * 2) {
			pass := InternetCheck(LogResponse, func(online bool) {
				Logln("switchAction", online)
			})
			Logln(pass)
		}
	}()

	time.Sleep(time.Minute * 5)
}

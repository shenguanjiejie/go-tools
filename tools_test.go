package tools

import (
	"testing"
	"time"
)

func TestInternet(t *testing.T) {

	go func() {
		for range time.Tick(time.Second * 2) {
			pass := InternetCheck(NetLogNone, func(online bool) {
				Logln("switchAction", online)
			})
			Logln(pass)
		}
	}()

	time.Sleep(time.Minute * 5)
}

func TestWaitHandle(t *testing.T) {
	waitChan := make(chan interface{}, 0)
	WaitHandle(waitChan, 10, func() {
		for i := 0; i < 1000; i++ {
			waitChan <- i
		}
	}, func(channelObj interface{}) {
		time.Sleep(time.Second)
		Logln(channelObj)
	})
}

package tools

import (
	"testing"
)

func TestLog(t *testing.T) {
	Slogln()
	Slogln("test")
	Slogln(Condition(true), "true")
	Slogln(Condition(false), "false") // RJ 2022-10-17 10:24:04 不打印
	Slogln(CallerLevel(0), "true")
	Slogln(CallerLevel(1), "true")
	Slogln(Condition(true), CallerLevel(1), "true 1")
	Slogln(Condition(false), CallerLevel(1), "false 1") // RJ 2022-10-17 10:24:09 不打印
}

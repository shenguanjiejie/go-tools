package main

import "github.com/shenguanjiejie/go-tools"

func main() {
	// tools.HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, &tools.HttpConfig{Log: tools.LogAll})

	testLog()
}

func testLog() {
	// tools.Slogln()
	// tools.Slogln("test")
	// tools.Slogln(tools.Condition(true), "true")
	// tools.Slogln(tools.Condition(false), "false") // RJ 2022-10-17 10:24:04 不打印
	// tools.Slogln(tools.CallerLevel(0), "true")
	// tools.Slogln(tools.CallerLevel(1), "true")
	tools.Slogln(tools.Condition(true), tools.CallerLevel(1), "true 1")
	tools.Slogln(tools.Condition(false), tools.CallerLevel(1), "false 1") // RJ 2022-10-17 10:24:09 不打印
}

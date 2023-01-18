/*
 * @Author: shenguanjiejie 835166018@qq.com
 * @Date: 2022-10-17 11:51:51
 * @LastEditors: shenguanjiejie 835166018@qq.com
 * @LastEditTime: 2022-10-20 12:09:53
 * @FilePath: /go-tools/main/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"github.com/shenguanjiejie/go-tools"
)

func main() {
	// tools.HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, &tools.HttpConfig{Log: tools.LogAll})

	testLog()
}

func testLog() {
	// tools.Logln()
	// tools.Logln("test")
	// tools.Logln(tools.LogCondition(true), "true")
	// tools.Logln(tools.LogCondition(false), "false") // RJ 2022-10-17 10:24:04 不打印
	// tools.Logln(tools.LogCallerLevel(0), "true")
	// tools.Logln(tools.LogCallerLevel(1), "true")
	tools.Logln(tools.LogCondition(true), tools.LogCallerLevel(1), "true 1")
	tools.Logln(tools.LogCondition(false), tools.LogCallerLevel(1), "false 1") // RJ 2022-10-17 10:24:09 不打印
}

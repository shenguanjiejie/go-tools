/*
 * @Author: shenguanjiejie 835166018@qq.com
 * @Date: 2022-10-17 13:05:27
 * @LastEditors: shenguanjiejie 835166018@qq.com
 * @LastEditTime: 2022-10-19 18:49:50
 * @FilePath: /go-tools/logger_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tools

import (
	"testing"
)

func TestLog(t *testing.T) {
	Slogln()
	Slogln("test")
	Slogln(Condition(true), "true")
	Slogln(Condition(false), "false") // RJ 2022-10-17 10:24:04 不打印
	Slogln(CallerLevel(0), 0)
	Slogln(CallerLevel(1), 1)
	Slogln(Condition(true), CallerLevel(1), "true 1")
	Slogln(Condition(false), CallerLevel(1), "false 1") // RJ 2022-10-17 10:24:09 不打印
}

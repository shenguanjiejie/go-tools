/*
 * @Author: shenguanjiejie 835166018@qq.com
 * @Date: 2022-10-17 13:05:27
 * @LastEditors: shenguanjiejie 835166018@qq.com
 * @LastEditTime: 2022-10-20 12:32:31
 * @FilePath: /go-tools/logger_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tools

import (
	"fmt"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	Logln()
	Logln(nil, "test", 1)
	Logln(Condition(true), "true")
	Logln(Condition(false), "false") // RJ 2022-10-17 10:24:04 不打印
	Logln(CallerLevel(0), 0)
	Logln(CallerLevel(1), 1)
	Logln(Condition(true), CallerLevel(0), nil, "true 1")
	Logln(Condition(false), CallerLevel(1), "false 1") // RJ 2022-10-17 10:24:09 不打印
	timeNow := time.Now().Add(1 * time.Hour)
	Logln(timeNow)
	Logln(time.Now())
	fmt.Println(timeNow)
	fmt.Println(time.Now())
}

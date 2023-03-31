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
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	func() {
		Logln()
		Logln(nil, "test", 1)
		Logln(LogCondition(true), "true")
		Logln(LogCondition(false), "false") // 不打印
		Logln(LogCallerLevel(0), 0)
		Logln(LogCallerLevel(1), LogLineLevel(0), 1)
		Logln(nil, LogLineLevel(1), "true 1", LogCallerLevel(0), []int{1, 2, 3}, LogCondition(true)) // All
		Logln(LogCondition(false), LogCallerLevel(1), "false 1")                                     // 不打印
		Logln(time.Now())
		fmt.Println(time.Now())
	}()

}

func TestLogf(t *testing.T) {
	str := "logf"
	i := 100
	obj := struct{ Logf string }{"I'm an object"}
	Logf("%s,%d,%v\n", str, i, obj)
}

func TestLogType(t *testing.T) {
	s := "指针"
	sP := &s

	arr := []int{1, 2, 3}
	arrP := &arr

	obj := struct{ a int }{1}
	objP := &obj

	Logln(nil)
	Logln(s, *sP)
	Logln(arr, arrP)
	Logln(obj, objP)
	Logln(1, int64(1))
	Logln(errors.New("error"))
	Logln(true)
}

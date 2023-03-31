/*
 * @Author: shenguanjiejie 835166018@qq.com
 * @Date: 2022-10-17 11:51:51
 * @LastEditors: shenguanjiejie 835166018@qq.com
 * @LastEditTime: 2022-10-20 12:09:53
 * @FilePath: /go-tools/main/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

// "go.uber.org/zap"

func main() {
	testLog()
}

func testLog() {
	// logger, _ := zap.NewProduction()
	// defer logger.Sync() // flushes buffer, if any
	// sugar := logger.Sugar()
	// tools.SetLogger(sugar)
	// tools.SetBaseFormat(func(timeStr string, funcName string, line int) (format string, args []interface{}) {
	// 	return "%s--%d: ", []interface{}{timeStr, line}
	// })

	// tools.Logln(tools.LogLevelError, "err_msg")
	// num := 100
	// numF := 3.14
	// sugar.Debugf("%d_%f", num, numF)
	// tools.Logf("%d_%f", tools.LogLevelDebug, num, numF)
	// tools.Log("%d_%f", tools.LogLevelWarn, num, numF)
	// tools.Logln(num, nil, time.Now(), "哈哈")
}

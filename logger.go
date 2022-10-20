package tools

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// RJ 2022-10-14 15:55:01 条件打印, 如果要设置该条件, 需放在第一个参数位置
type Condition bool

var conditionType = reflect.TypeOf(Condition(true))

// RJ 2022-10-14 15:23:01 输出的方法名的层级, 默认为0, 代表输出当前方法名, 如果要输出上层方法名(比如闭包内打印), 则第一个参数设置为CallerLevel(1)即可, 以此类推.
type CallerLevel int

var callerLevelType = reflect.TypeOf(CallerLevel(0))

// Log 带行号输出
func Log(a ...interface{}) {
	condition, newA, pc, codeLine, ok := logStackInfo(a...)

	if !condition {
		return
	}

	if !ok {
		fmt.Print(newA...)
		return
	}

	format, slice := formatWithValues(pc, codeLine, slogFormat(newA), newA...)
	fmt.Printf(format, slice...)
}

// Logln 带行号输出
func Logln(a ...interface{}) {
	condition, newA, pc, codeLine, ok := logStackInfo(a...)

	if !condition {
		return
	}

	if !ok {
		fmt.Println(newA...)
		return
	}

	format, slice := formatWithValues(pc, codeLine, slogFormat(newA), newA...)
	fmt.Printf(format+"\n", slice...)
}

// Logf 带行号格式输出
func Logf(format string, a ...interface{}) {
	condition, newA, pc, codeLine, ok := logStackInfo(a...)

	if !condition {
		return
	}

	if !ok {
		fmt.Printf(format, newA...)
		return
	}

	finalFormat, slice := formatWithValues(pc, codeLine, format, newA...)
	fmt.Printf(finalFormat, slice...)
}

func logStackInfo(a ...interface{}) (condition bool, newA []interface{}, pc uintptr, line int, ok bool) {
	newA = a
	currentLevel := 2

	callerHandle := func(data interface{}) uintptr {
		if data == nil {
			return 0
		}

		// RJ 2022-10-17 10:41:05 获取对应CallerLevel的pc
		if reflect.TypeOf(data) == callerLevelType {
			// RJ 2022-10-17 10:52:52 由于在闭包内, 所以需要再+1
			callerInt := int(data.(CallerLevel)) + 1
			newA = newA[1:]
			pc, _, _, ok = runtime.Caller(callerInt + currentLevel)
			if !ok {
				return 0
			}
		}
		return pc
	}

	// RJ 2022-10-14 15:16:56 判断是否传入了Condition或CallerLevel来指定输出的层
	if len(a) > 0 && a[0] != nil {
		data := a[0]

		if reflect.TypeOf(data) == conditionType {
			// RJ 2022-10-14 16:02:08 Condition是不是false, 如果是, 则无需打印
			if !data.(Condition) {
				return
			}

			newA = a[1:]

			if len(newA) > 0 {
				pc = callerHandle(newA[0])
			}
		} else {
			pc = callerHandle(data)
		}
	}

	// RJ 2022-10-14 15:18:50 未指定输出的方法层, 则默认取当前层pc
	if pc == 0 {
		pc, _, line, ok = runtime.Caller(currentLevel)
	} else {
		_, _, line, ok = runtime.Caller(currentLevel)
	}
	return true, newA, pc, line, ok
}

func slogFormat(a ...interface{}) string {
	formatStr := ""
	dataArr, ok := a[0].([]interface{})
	if !ok {
		formatStr += "%v, "
	}

	for _, data := range dataArr {
		switch data.(type) {
		case nil, bool, int, int64, string, error:
			formatStr += "%v, "
		default:
			formatStr += "%#v, "
		}
	}
	return formatStr
}

func formatWithValues(pc uintptr, codeLine int, format string, a ...interface{}) (string, []interface{}) {
	funName := runtime.FuncForPC(pc).Name()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	slice := []interface{}{timeStr, funName, codeLine}
	slice = append(slice, a...)
	return "%s--%s--第%d行--: " + format, slice
}

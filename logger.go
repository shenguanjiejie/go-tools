package tools

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

/**RJ 2023-01-18 14:37:40
LogCondition & LogCallerLevel & LogLineLevel
1. 仅用来配置打印的条件, 这几个参数在log之前会被移除, 不会打印出来.
2. 不限制放在log中的位置, 不过建议Condition放在最前面

eg: tools.Logln(LogCondition(true),LogCallerLevel(1),LogLineLevel(2),"test", nil, []int{1,2,3})
*/

// RJ 2022-10-14 15:55:01 条件打印
type LogCondition bool

var conditionType = reflect.TypeOf(LogCondition(true))

// RJ 2022-10-14 15:23:01 输出的方法名的层级, 默认为0, 代表输出当前方法名, 如果要输出上层方法名(比如闭包内打印), 则该参数设置为CallerLevel(1)即可, 以此类推.
type LogCallerLevel int

var callerLevelType = reflect.TypeOf(LogCallerLevel(0))

// RJ 2023-01-18 13:36:16  输出的行号的层级, 默认为0, 代表输出当前所在代码块的行号, 如果要输出上层代码块的行号(比如闭包内打印), 则该参数设置为LineLevel(1)即可, 以此类推.
type LogLineLevel int

var lineLevelType = reflect.TypeOf(LogLineLevel(0))

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

	format, slice := formatWithValues(pc, codeLine, logFormat(newA), newA...)
	fmt.Printf(format, slice...)
}

// Logln 带行号换行输出
func Logln(a ...interface{}) {
	condition, newA, pc, codeLine, ok := logStackInfo(a...)

	if !condition {
		return
	}

	if !ok {
		fmt.Println(newA...)
		return
	}

	format, slice := formatWithValues(pc, codeLine, logFormat(newA), newA...)
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

	condition = true

	for i := 0; i < len(newA); i++ {
		obj := newA[i]
		// RJ 2023-01-18 14:27:32 不是logger用来做判定的类型
		if obj == nil || (reflect.TypeOf(obj) != conditionType && reflect.TypeOf(obj) != callerLevelType && reflect.TypeOf(obj) != lineLevelType) {
			continue
		}
		if reflect.TypeOf(obj) == conditionType {
			if !obj.(LogCondition) {
				condition = false
			}
		} else if condition && reflect.TypeOf(obj) == callerLevelType {
			callerLevelInt := int(obj.(LogCallerLevel))
			if line == 0 {
				pc, _, line, ok = runtime.Caller(callerLevelInt + currentLevel)
			} else {
				pc, _, _, ok = runtime.Caller(callerLevelInt + currentLevel)
			}

		} else if condition && reflect.TypeOf(obj) == lineLevelType {
			lineLevelInt := int(obj.(LogLineLevel))
			if pc == 0 {
				pc, _, line, _ = runtime.Caller(lineLevelInt + currentLevel)
			} else {
				_, _, line, _ = runtime.Caller(lineLevelInt + currentLevel)
			}
		}

		newA = append(newA[:i], newA[i+1:]...)
		i--
	}

	if condition && pc == 0 && line == 0 {
		pc, _, line, ok = runtime.Caller(currentLevel)
	}
	return
}

func logFormat(a ...interface{}) string {
	formatStr := ""
	dataArr, ok := a[0].([]interface{})
	if !ok {
		formatStr += "%v, "
	}

	for _, data := range dataArr {
		switch data.(type) {
		case nil, bool, int, int32, int64, string, error, float32, float64:
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

package tools

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

/**
LogCondition & LogCallerLevel & LogLineLevel
1. 仅用来配置打印的条件, 这几个参数在log之前会被移除, 不会打印出来.
2. 不限制放在参数中的位置, 不过建议Condition放在最前面
3. go-tools是一个轻量级的工具库, 本身没有实现对LogLevel的支持. 如果有LogLevel的需求, 可以集成其他logger(调用SetLogger方法设置). 前提是要支持下方定义的Logger接口.

eg: tools.Logln(LogCondition(true),LogCallerLevel(1), LogLevelError,"err_msg", nil, []int{1,2,3})
*/

// 条件打印
type LogCondition bool

var conditionType = reflect.TypeOf(LogCondition(true))

// 输出的方法名的层级, 默认为0, 代表输出当前方法名, 如果要输出上层方法名(比如闭包内打印), 则该参数设置为CallerLevel(1)即可, 以此类推.
type LogCallerLevel int

var callerLevelType = reflect.TypeOf(LogCallerLevel(0))

// 输出的行号的层级, 默认为0, 代表输出当前所在代码块的行号, 如果要输出上层代码块的行号(比如闭包内打印), 则该参数设置为LineLevel(1)即可, 以此类推.
type LogLineLevel int

var lineLevelType = reflect.TypeOf(LogLineLevel(0))

type logLevel int

// 如果集成了其他实现了Logger接口的logger, 可以把级别作为参数传入, 默认为LogLevelInfo
const (
	logLevelInfo logLevel = iota
	logLevelDebug
	logLevelWarn
	logLevelError
)

var logLevelType = reflect.TypeOf(logLevelInfo)

// 调用SetLogger方法设置logger
var logger Logger

var baseLogBlock func(timeStr string, funcName string, line int) (format string, args []interface{})

type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})
}

// Log 带行号输出
func Log(a ...interface{}) {
	condition, newA, pc, codeLine, ok, logLevel := logStackInfo(a...)

	if !condition {
		return
	}

	if !ok {
		fmt.Print(newA...)
		return
	}

	format, slice := formatWithValues(pc, codeLine, logFormat(newA), newA...)

	if logger == nil {
		fmt.Printf(format, slice...)
	} else {
		logLevelLog(logLevel, format, slice...)
	}
}

// Logln 带行号换行输出
func Logln(a ...interface{}) {
	condition, newA, pc, codeLine, ok, logLevel := logStackInfo(a...)

	if !condition {
		return
	}

	if !ok {
		fmt.Println(newA...)
		return
	}

	format, slice := formatWithValues(pc, codeLine, logFormat(newA), newA...)

	if logger == nil {
		fmt.Printf(format+"\n", slice...)
	} else {
		logLevelLog(logLevel, format, slice...)
	}
}

// Logf 带行号格式输出
func Logf(format string, a ...interface{}) {
	condition, newA, pc, codeLine, ok, logLevel := logStackInfo(a...)

	if !condition {
		return
	}

	if !ok {
		fmt.Printf(format, newA...)
		return
	}

	finalFormat, slice := formatWithValues(pc, codeLine, format, newA...)
	if logger == nil {
		fmt.Printf(finalFormat, slice...)
	} else {
		logLevelLog(logLevel, finalFormat, slice...)
	}
}

func logStackInfo(a ...interface{}) (condition bool, newA []interface{}, pc uintptr, line int, ok bool, level logLevel) {
	newA = a
	currentLevel := 2
	condition = true
	level = logLevelInfo

	for i := 0; i < len(newA); i++ {
		obj := newA[i]
		// 不是logger用来做判定的类型
		objType := reflect.TypeOf(obj)
		if obj == nil || (objType != conditionType && objType != callerLevelType && objType != lineLevelType && objType != logLevelType) {
			continue
		}
		if objType == conditionType {
			if !obj.(LogCondition) {
				condition = false
			}
		} else if condition && objType == callerLevelType {
			callerLevelInt := int(obj.(LogCallerLevel))
			if line == 0 {
				pc, _, line, ok = runtime.Caller(callerLevelInt + currentLevel)
			} else {
				pc, _, _, ok = runtime.Caller(callerLevelInt + currentLevel)
			}

		} else if condition && objType == lineLevelType {
			lineLevelInt := int(obj.(LogLineLevel))
			if pc == 0 {
				pc, _, line, _ = runtime.Caller(lineLevelInt + currentLevel)
			} else {
				_, _, line, _ = runtime.Caller(lineLevelInt + currentLevel)
			}
		} else if condition && logger != nil && objType == logLevelType {
			level = obj.(logLevel)
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
		case nil, bool, int, int8, int16, int32, int64, string, error, float32, float64:
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
	baseFormat := "%s__%s__第%d行__: "
	slice := []interface{}{timeStr, funName, codeLine}

	if baseLogBlock != nil {
		baseFormat, slice = baseLogBlock(timeStr, funName, codeLine)
	}
	slice = append(slice, a...)
	return baseFormat + format, slice
}

func SetLogger(yourLogger Logger) {
	logger = yourLogger
}

func SetBaseFormat(block func(timeStr string, funcName string, line int) (format string, args []interface{})) {
	baseLogBlock = block
}

func logLevelLog(level logLevel, format string, slice ...interface{}) {
	switch level {
	case logLevelInfo:
		logger.Infof(format, slice...)
	case logLevelDebug:
		logger.Debugf(format, slice...)
	case logLevelWarn:
		logger.Warnf(format, slice...)
	case logLevelError:
		logger.Errorf(format, slice...)
	}
}

func Debug(args ...interface{}) {
	Logln(append(args, logLevelDebug)...)
}

func Info(args ...interface{}) {
	Logln(append(args, logLevelInfo)...)
}

func Warn(args ...interface{}) {
	Logln(append(args, logLevelWarn)...)
}

func Error(args ...interface{}) {
	Logln(append(args, logLevelError)...)
}

func Debugf(template string, args ...interface{}) {
	Logf(template, append(args, logLevelDebug)...)
}

func Infof(template string, args ...interface{}) {
	Logf(template, append(args, logLevelInfo)...)
}

func Warnf(template string, args ...interface{}) {
	Logf(template, append(args, logLevelWarn)...)
}

func Errorf(template string, args ...interface{}) {
	Logf(template, append(args, logLevelError)...)
}

package tools

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"
)

/**
LogCondition & LogCallerLevel & LogLineLevel
1. 仅用来配置打印的条件, 这几个参数在log之前会被移除, 不会打印出来.
2. 不限制放在参数中的位置, 不过建议Condition放在最前面
3. go-tools支持集成其他logger(调用SetLogger方法设置). 前提是要支持下方定义的Logger接口.

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

type logLevel string

const (
	logLevelInfo  logLevel = "info"
	logLevelDebug logLevel = "debug"
	logLevelWarn  logLevel = "warn"
	logLevelError logLevel = "error"
)

var logLevelType = reflect.TypeOf(logLevelInfo)

// 调用SetLogger方法设置logger
var logger Logger

var baseLogBlock func(timeStr string, level string, funcName string, line int) (format string, args []interface{})

type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})
}

// Log 带行号输出
func Log(args ...interface{}) {
	log("", false, func() {
		fmt.Print(args...)
	}, args...)
}

// Logln 带行号换行输出
func Logln(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, args...)
}

// Logf 带行号格式输出
func Logf(format string, args ...interface{}) {
	log(format, false, func() {
		fmt.Printf(format, args...)
	}, args...)
}

func logStackInfo(args ...interface{}) (condition bool, newA []interface{}, pc uintptr, line int, ok bool, level logLevel) {
	newA = args
	currentLevel := 3
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
		} else if condition && objType == logLevelType {
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

func logFormat(args ...interface{}) string {
	formatStr := ""
	dataArr, ok := args[0].([]interface{})
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

func formatWithValues(pc uintptr, level logLevel, codeLine int, format string, args ...interface{}) (string, []interface{}) {
	funName := runtime.FuncForPC(pc).Name()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	baseFormat := "%s__%s__%s__第%d行__: "
	slice := []interface{}{timeStr, level, funName, codeLine}

	if baseLogBlock != nil {
		baseFormat, slice = baseLogBlock(timeStr, string(level), funName, codeLine)
	}
	slice = append(slice, args...)
	return baseFormat + format, slice
}

func SetLogger(yourLogger Logger) {
	logger = yourLogger
}

// 设置基本信息相关format格式, 默认为: "%s__%s__%s__第%d行__: " 和 []interface{}{timeStr, level, funName, codeLine}
func SetBaseFormat(block func(timeStr string, level string, funcName string, line int) (format string, args []interface{})) {
	baseLogBlock = block
}

func log(format string, ln bool, ifErr func(), a ...interface{}) {

	condition, newA, pc, codeLine, ok, level := logStackInfo(a...)

	if !condition {
		return
	}

	if !ok {
		ifErr()
		return
	}

	if format == "" {
		format = logFormat(newA) + format
	}

	if ln {
		format = format + "\n"
	}

	finalFormat, slice := formatWithValues(pc, level, codeLine, format, newA...)

	if logger == nil {
		fmt.Printf(finalFormat, slice...)
		if level == logLevelError {
			debug.PrintStack()
		}
		return
	}

	switch level {
	case logLevelInfo:
		logger.Infof(finalFormat, slice...)
	case logLevelDebug:
		logger.Debugf(finalFormat, slice...)
	case logLevelWarn:
		logger.Warnf(finalFormat, slice...)
	case logLevelError:
		logger.Errorf(finalFormat, slice...)
	}
}

func Debug(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, append(args, logLevelDebug)...)
}

func Info(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, append(args, logLevelInfo)...)
}

func Warn(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, append(args, logLevelWarn)...)
}

func Error(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, append(args, logLevelError)...)
}

func Debugf(template string, args ...interface{}) {
	log(template, false, func() {
		fmt.Printf(template, args...)
	}, append(args, logLevelDebug)...)
}

func Infof(template string, args ...interface{}) {
	log(template, false, func() {
		fmt.Printf(template, args...)
	}, append(args, logLevelInfo)...)
}

func Warnf(template string, args ...interface{}) {
	log(template, false, func() {
		fmt.Printf(template, args...)
	}, append(args, logLevelWarn)...)
}

func Errorf(template string, args ...interface{}) {
	log(template, false, func() {
		fmt.Printf(template, args...)
	}, append(args, logLevelError)...)
}

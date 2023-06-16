package tools

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"
)

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

// Logger go-tools支持集成其他logger(调用SetLogger方法设置). 前提是要支持下方定义的Logger接口.
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
	currentSkip := 3
	condition = true
	level = logLevelInfo
	o := new(logOptions)
	o.LogCondition = &condition

	for i := 0; i < len(newA); i++ {
		obj := newA[i]
		objType := reflect.TypeOf(obj)
		// 不是logger用来做判定的类型
		if obj == nil || (objType != logOptionType && objType != logLevelType) {
			continue
		}

		if objType == logOptionType {
			obj.(LogOptionFunc)(o)
		}

		if *o.LogCondition {
			if o.LogCallerSkip > 0 && pc == 0 {
				if line == 0 {
					pc, _, line, ok = runtime.Caller(o.LogCallerSkip + currentSkip)
				} else {
					pc, _, _, ok = runtime.Caller(o.LogCallerSkip + currentSkip)
				}
			} else if o.LogLineSkip > 0 && line == 0 {
				if pc == 0 {
					pc, _, line, _ = runtime.Caller(o.LogLineSkip + currentSkip)
				} else {
					_, _, line, _ = runtime.Caller(o.LogLineSkip + currentSkip)
				}
			} else if objType == logLevelType {
				level = obj.(logLevel)
			}
		}

		newA = append(newA[:i], newA[i+1:]...)
		i--
	}

	condition = *o.LogCondition
	if condition && pc == 0 && line == 0 {
		pc, _, line, ok = runtime.Caller(currentSkip)
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

// SetLogger 设置日志输出实例
func SetLogger(yourLogger Logger) {
	logger = yourLogger
}

// SetBaseFormat 设置基本信息相关format格式, 默认为: "%s__%s__%s__第%d行__: " 和 []interface{}{timeStr, level, funName, codeLine}
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

// Debug debug
func Debug(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, append(args, logLevelDebug)...)
}

// Info info
func Info(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, append(args, logLevelInfo)...)
}

// Warn warn
func Warn(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, append(args, logLevelWarn)...)
}

// Error error
func Error(args ...interface{}) {
	log("", true, func() {
		fmt.Println(args...)
	}, append(args, logLevelError)...)
}

// Debugf debug with template
func Debugf(template string, args ...interface{}) {
	log(template, false, func() {
		fmt.Printf(template, args...)
	}, append(args, logLevelDebug)...)
}

// Infof info with template
func Infof(template string, args ...interface{}) {
	log(template, false, func() {
		fmt.Printf(template, args...)
	}, append(args, logLevelInfo)...)
}

// Warnf warn with template
func Warnf(template string, args ...interface{}) {
	log(template, false, func() {
		fmt.Printf(template, args...)
	}, append(args, logLevelWarn)...)
}

// Errorf error with template
func Errorf(template string, args ...interface{}) {
	log(template, false, func() {
		fmt.Printf(template, args...)
	}, append(args, logLevelError)...)
}

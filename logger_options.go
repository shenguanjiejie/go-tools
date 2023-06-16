package tools

import "reflect"

// LogOptionFunc for web
type LogOptionFunc func(o *logOptions)

var logOptionType = reflect.TypeOf(LogOptionFunc(nil))

// LogOptions 打印配置
// LogCondition & LogCallerSkip & LogLineSkip
// 1. 仅用来配置打印的条件, 这几个参数在log之前会被移除, 不会打印出来.
// 2. 不限制放在参数中的位置, 不过建议Condition放在最前面
type logOptions struct {
	LogCondition  *bool
	LogCallerSkip int
	LogLineSkip   int
}

// LogCondition 打印条件, 为true才打印, 默认true
func LogCondition(condition bool) LogOptionFunc {
	return func(o *logOptions) {
		o.LogCondition = &condition
	}
}

// LogCallerSkip 输出的方法名的层级, 默认为0, 代表输出当前方法名, 如果要输出上层方法名(比如闭包内打印), 则该参数设置为CallerLevel(1)即可, 以此类推.
func LogCallerSkip(callerSkip int) LogOptionFunc {
	return func(o *logOptions) {
		o.LogCallerSkip = callerSkip
	}
}

// LogLineSkip 输出的行号的层级, 默认为0, 代表输出当前所在代码块的行号, 如果要输出上层代码块的行号(比如闭包内打印), 则该参数设置为LineLevel(1)即可, 以此类推.
func LogLineSkip(lineSkip int) LogOptionFunc {
	return func(o *logOptions) {
		o.LogLineSkip = lineSkip
	}
}

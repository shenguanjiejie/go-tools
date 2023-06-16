package tools

import (
	"net/http"
	"reflect"
	"time"
)

// NetOptionFunc for web
type NetOptionFunc func(o *NetOptions)

var netOptionType = reflect.TypeOf(NetOptionFunc(nil))

// NetOptions 额外配置, 未进行配置的项, 会使用默认值
type NetOptions struct {
	Header        http.Header
	NetLogLevel   NetLogLevel   // default: NetLogAllWithoutObj
	LogCallerSkip int           // 默认为LogCallerSkip(0), 代表请求位置的方法所跳过的层数,如果想看tools内部打印所在的方法, 可以传tools.LogCallerSkip(-2) TODO: zap中用callerSkip
	LogLineSkip   int           // , 代表请求位置的行号所跳过的层数,如果想看tools内部的打印所在的行, 可以传tools.LogLineSkip(-2)
	Timeout       time.Duration // 为0会忽略
	UnmarshalPath []interface{} // 仅当obj参数不为nil时有效, eg:[]interface{}{"a",0,"b"}, 将会解析a下面的第1个元素的b节点
	contentType   string        // default: "application/json" , post only, 该参数不对外开放, 如有需求可以通过header进行设置.
}

// NetHeader header
func NetHeader(header http.Header) NetOptionFunc {
	return func(o *NetOptions) {
		o.Header = header
	}
}

// NetLogLevelOption default: NetLogAll
func NetLogLevelOption(netLogLevel NetLogLevel) NetOptionFunc {
	return func(o *NetOptions) {
		o.NetLogLevel = netLogLevel
	}
}

// LogCallerSkipOption 默认为0, 代表请求位置的方法所跳过的层数,如果想看tools内部打印所在的方法, 可以传-2
func LogCallerSkipOption(logCallerSkip int) NetOptionFunc {
	return func(o *NetOptions) {
		o.LogCallerSkip = logCallerSkip
	}
}

// LogLineSkipOption 默认为0, 代表请求位置的行号所跳过的层数,如果想看tools内部的打印所在的行, 可以传-2
func LogLineSkipOption(logLineSkip int) NetOptionFunc {
	return func(o *NetOptions) {
		o.LogLineSkip = logLineSkip
	}
}

// Timeout timeout
func Timeout(timeout time.Duration) NetOptionFunc {
	return func(o *NetOptions) {
		o.Timeout = timeout
	}
}

// UnmarshalPath 仅当obj参数不为nil时有效, eg:[]interface{}{"a",0,"b"}, 将会解析a下面的第1个元素的b节点
func UnmarshalPath(unmarshalPath []interface{}) NetOptionFunc {
	return func(o *NetOptions) {
		o.UnmarshalPath = unmarshalPath
	}
}

// // ContentType default: "application/json" , post only
// func ContentType(contentType string) NetOptionFunc {
// 	return func(o *NetOptions) {
// 		o.ContentType = contentType
// 	}
// }

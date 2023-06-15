package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// NetLogLevel 日志级别
type NetLogLevel int

// 日志配置, 默认LogAllWithoutObj
const (
	// 如果不需要打印, 请设置为NetLogNone, NetLogNil只是缺省值, 没有实际意义, 设置为NetLogNil相当于设置为NetLogAllWithoutObj
	NetLogNil  NetLogLevel = NetLogLevel(0)
	NetLogNone NetLogLevel = 1 << iota
	NetLogURL
	NetLogParams
	NetLogResponse
	// LogObj 打印反序列化后的obj
	NetLogObj
	NetLogError
	NetLogAllWithoutObj = NetLogURL | NetLogParams | NetLogResponse | NetLogError
	NetLogAll           = NetLogAllWithoutObj | NetLogObj
)

// HTTPConfig 额外配置, 未进行配置的项, 会使用默认值
type HTTPConfig struct {
	Header        http.Header
	NetLogLevel   NetLogLevel   // default: NetLogAll
	LogCallerSkip LogCallerSkip // 默认为LogCallerSkip(0), 代表请求位置的方法所跳过的层数,如果想看tools内部打印所在的方法, 可以传tools.LogCallerSkip(-2) TODO: zap中用callerSkip
	LogLineSkip   LogLineSkip   // 默认为LogLineSkip(0), 代表请求位置的行号所跳过的层数,如果想看tools内部的打印所在的行, 可以传tools.LogLineSkip(-2)
	Timeout       time.Duration // 为0会忽略
	UnmarshalPath []interface{} // 仅当obj参数不为nil时有效, eg:[]interface{}{"a",0,"b"}, 将会解析a下面的第1个元素的b节点
	contentType   string        // default: "application/json" , post only
}

type httpConfig struct {
	HTTPConfig
	Method string
	URL    string
	Params interface{}
	Body   io.Reader
}

var defaultConfig = &httpConfig{HTTPConfig: HTTPConfig{NetLogLevel: NetLogAllWithoutObj}}

// Get obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
func Get(urlStr string, values url.Values, obj interface{}, config ...*HTTPConfig) error {
	url := urlStr
	iconfig := configWithParams(config...)
	if values != nil {
		url = fmt.Sprintf("%s?%s", urlStr, values.Encode())
	}

	iconfig.Method = http.MethodGet
	iconfig.URL = url
	iconfig.Params = values
	err := request(obj, iconfig)
	return err
}

// Delete obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
func Delete(urlStr string, values url.Values, obj interface{}, config ...*HTTPConfig) error {
	url := urlStr
	iconfig := configWithParams(config...)
	if values != nil {
		url = fmt.Sprintf("%s?%s", urlStr, values.Encode())
	}

	iconfig.Method = http.MethodDelete
	iconfig.URL = url
	iconfig.Params = values
	err := request(obj, iconfig)
	return err
}

func requestWithData(method string, url string, data interface{}, obj interface{}, config ...*HTTPConfig) error {
	iconfig := configWithParams(config...)
	jsonParams, _ := json.Marshal(data)

	iconfig.URL = url
	iconfig.Method = http.MethodPost
	iconfig.Body = bytes.NewBuffer(jsonParams)
	iconfig.Params = data
	if iconfig.contentType == "" {
		iconfig.contentType = "application/json"
	}
	err := request(obj, iconfig)

	return err
}

// Post obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
func Post(url string, data interface{}, obj interface{}, config ...*HTTPConfig) error {
	return requestWithData(http.MethodPost, url, data, obj, config...)
}

// FormDataPost obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
func FormDataPost(url string, data map[string]string, obj interface{}, config ...*HTTPConfig) error {
	cmdResReqForm, contentType := createMultipartFormBody(data)
	iconfig := configWithParams(config...)

	iconfig.URL = url
	iconfig.Method = http.MethodPost
	iconfig.Body = cmdResReqForm
	iconfig.Params = data
	iconfig.contentType = contentType
	request(obj, iconfig)
	return nil
}

func createMultipartFormBody(params map[string]string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	// Add fields
	for key, val := range params {
		if key == "file" {
			// Open file
			f, err := os.Open(val)
			if err != nil {
				return nil, ""
			}
			defer f.Close()

			// Add file fields
			fw, err := w.CreateFormFile(key, val)
			if err != nil {
				return nil, ""
			}
			if _, err = io.Copy(fw, f); err != nil {
				return nil, ""
			}
		} else {
			// Add string fields
			fw, err := w.CreateFormField(key)
			if err != nil {
				return nil, ""
			}
			if _, err = fw.Write([]byte(val)); err != nil {
				return nil, ""
			}
		}
	}
	w.Close()

	return body, w.FormDataContentType()
}

// Put obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
func Put(url string, data interface{}, obj interface{}, config ...*HTTPConfig) error {
	return requestWithData(http.MethodPut, url, data, obj, config...)
}

// Patch obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
func Patch(url string, data interface{}, obj interface{}, config ...*HTTPConfig) error {
	return requestWithData(http.MethodPatch, url, data, obj, config...)
}

func configWithParams(config ...*HTTPConfig) *httpConfig {
	if len(config) > 0 {
		return &httpConfig{HTTPConfig: *config[0]}
	}
	return defaultConfig
}

func request(obj interface{}, config *httpConfig) error {
	client := http.DefaultClient
	if config.NetLogLevel == NetLogNil {
		config.NetLogLevel = NetLogAll
	}
	shouldLogError := LogCondition(config.NetLogLevel&NetLogError != 0)
	callerLevel := config.LogCallerSkip + 2
	lineLevel := config.LogLineSkip + 2
	Logln(LogCondition(config.NetLogLevel&NetLogURL != 0), callerLevel, lineLevel, config.Method, config.URL)
	Logln(LogCondition(config.NetLogLevel&NetLogParams != 0), callerLevel, lineLevel, config.Params)

	request, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		Error(shouldLogError, callerLevel, lineLevel, err)
		return err
	}

	if config.Header != nil {
		request.Header = config.Header
	}

	if config.Method != http.MethodGet && config.Method != http.MethodDelete && request.Header.Get("Content-Type") == "" {
		request.Close = true
		request.Header.Add("Content-Type", config.contentType)
	}

	if config.Timeout > 0 {
		client.Timeout = config.Timeout
	}

	response, err := client.Do(request)
	if err != nil {
		Error(shouldLogError, callerLevel, lineLevel, err)
		return err
	}

	if obj != nil && reflect.TypeOf(obj) == reflect.TypeOf(response) {
		*(obj.(*http.Response)) = *response
		Logln(LogCondition(config.NetLogLevel&NetLogResponse != 0), callerLevel, lineLevel, obj)
		return nil
	}
	result, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		Error(shouldLogError, callerLevel, lineLevel, err)
		return err
	}

	Logln(LogCondition(config.NetLogLevel&NetLogResponse != 0), callerLevel, lineLevel, string(result))

	if obj != nil {
		if len(config.UnmarshalPath) > 0 {
			value := jsoniter.Get(result, config.UnmarshalPath...)
			result = []byte(value.ToString())
		}
		err = jsoniter.Unmarshal(result, &obj)
		if err != nil {
			Error(shouldLogError, callerLevel, lineLevel, err)
			return err
		}
		Logln(LogCondition(config.NetLogLevel&NetLogObj != 0), callerLevel, lineLevel, obj)
	}

	return nil
}

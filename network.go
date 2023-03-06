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
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

//LogLevel 日志级别
type LogLevel int

// RJ 2022-10-14 11:43:33 日志配置, 默认LogAll
const (
	// 如果不需要打印, 请设置为LogNone, LogNil只是缺省值, 没有实际意义, 设置为LogNil相当于设置为LogAll
	LogNil  LogLevel = LogLevel(0)
	LogNone LogLevel = 1 << iota
	LogURL
	LogParams
	LogResponse
	//  LogObj 打印反序列化后的obj
	LogObj
	LogError
	LogAllWithoutObj = LogURL | LogParams | LogResponse | LogError
	LogAll           = LogAllWithoutObj | LogObj
)

const (
	HTTPMethodGet    = "GET"
	HTTPMethodPost   = "POST"
	HTTPMethodPut    = "PUT"
	HTTPMethodPatch  = "PATCH" // RFC 5789
	HTTPMethodDelete = "DELETE"
)

// HttpConfig 额外配置, 未进行配置的项, 会使用默认值
type HttpConfig struct {
	Header        http.Header
	Log           LogLevel       // default: LogAll
	LogCaller     LogCallerLevel // 默认为LogCallerLevel(0), 代表请求位置的方法,如果想看tools内部打印所在的方法, 可以传LogCallerLevel(-2)
	LogLine       LogLineLevel   // 默认为LogLineLevel(0), 代表请求位置的行号,如果想看tools内部的打印所在的行, 可以传LogLineLevel(-2)
	Timeout       time.Duration  // 为0会忽略
	ContentType   string         // default: "application/json" , post only
	UnmarshalNode string         //  仅当obj参数不为nil时有效, eg:"a.b.c", 则解析json中的a下面的b下面的c节点.
}

type httpConfig struct {
	HttpConfig
	Method string
	URL    string
	Params interface{}
	Body   io.Reader
}

var defaultConfig = &httpConfig{HttpConfig: HttpConfig{Log: LogAll}}
var jsonNodeRegex = `((\{.+\})|(\[.+\])|(("(?:[^"\\]|\\.)*")|((\d+\.\d+)|(\d+))))((?=,)|(^,.*\}|\])))`

/*
*HttpGet
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpGet(urlStr string, values url.Values, obj interface{}, config ...*HttpConfig) error {
	url := urlStr
	iconfig := configWithParams(config...)
	if values != nil {
		url = fmt.Sprintf("%s?%s", urlStr, values.Encode())
	}

	iconfig.Method = HTTPMethodGet
	iconfig.URL = url
	iconfig.Params = values
	err := request(obj, iconfig)
	return err
}

/*
*HttpGet
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpDelete(urlStr string, values url.Values, obj interface{}, config ...*HttpConfig) error {
	url := urlStr
	iconfig := configWithParams(config...)
	if values != nil {
		url = fmt.Sprintf("%s?%s", urlStr, values.Encode())
	}

	iconfig.Method = HTTPMethodDelete
	iconfig.URL = url
	iconfig.Params = values
	err := request(obj, iconfig)
	return err
}

func post(method string, url string, dataMap map[string]string, obj interface{}, config ...*HttpConfig) error {
	iconfig := configWithParams(config...)
	jsonParams, _ := json.Marshal(dataMap)

	iconfig.URL = url
	iconfig.Method = HTTPMethodPost
	iconfig.Body = bytes.NewBuffer(jsonParams)
	iconfig.Params = dataMap
	if iconfig.ContentType == "" {
		iconfig.ContentType = "application/json"
	}
	err := request(obj, iconfig)

	return err
}

/*
*HttpPost
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpPost(url string, dataMap map[string]string, obj interface{}, config ...*HttpConfig) error {
	return post(HTTPMethodPost, url, dataMap, obj, config...)
}

/*
*HttpFormDataPost
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpFormDataPost(url string, dataMap map[string]string, obj interface{}, config ...*HttpConfig) error {
	cmdResReqForm, contentType := createMultipartFormBody(dataMap)
	iconfig := configWithParams(config...)

	iconfig.URL = url
	iconfig.Method = HTTPMethodPost
	iconfig.Body = cmdResReqForm
	iconfig.Params = dataMap
	iconfig.ContentType = contentType
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

/*
*HttpPost
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpPut(url string, dataMap map[string]string, obj interface{}, config ...*HttpConfig) error {
	return post(HTTPMethodPut, url, dataMap, obj, config...)
}

/*
*HttpPost
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpPatch(url string, dataMap map[string]string, obj interface{}, config ...*HttpConfig) error {
	return post(HTTPMethodPatch, url, dataMap, obj, config...)
}

func configWithParams(config ...*HttpConfig) *httpConfig {
	if len(config) > 0 {
		return &httpConfig{HttpConfig: *config[0]}
	}
	return defaultConfig
}

func request(obj interface{}, config *httpConfig) error {
	client := http.DefaultClient
	if config.Log == LogNil {
		config.Log = LogAll
	}
	shouldLogError := LogCondition(config.Log&LogError != 0)
	callerLevel := config.LogCaller + 2
	lineLevel := config.LogLine + 2
	Logln(LogCondition(config.Log&LogURL != 0), callerLevel, lineLevel, config.Method, config.URL)
	Logln(LogCondition(config.Log&LogParams != 0), callerLevel, lineLevel, config.Params)

	request, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		Logln(shouldLogError, callerLevel, lineLevel, err)
		return err
	}

	if config.Header != nil {
		request.Header = config.Header
	}

	if config.Method != HTTPMethodGet && config.Method != HTTPMethodDelete && request.Header.Get("Content-Type") == "" {
		request.Close = true
		request.Header.Add("Content-Type", config.ContentType)
	}

	if config.Timeout > 0 {
		client.Timeout = config.Timeout
	}

	response, err := client.Do(request)
	if err != nil {
		Logln(shouldLogError, callerLevel, lineLevel, err)
		return err
	}

	if obj != nil && reflect.TypeOf(obj) == reflect.TypeOf(response) {
		*(obj.(*http.Response)) = *response
		Logln(LogCondition(config.Log&LogResponse != 0), callerLevel, lineLevel, obj)
		return nil
	}
	result, err := io.ReadAll(response.Body)
	if err != nil {
		Logln(shouldLogError, callerLevel, lineLevel, err)
		return err
	}
	if obj != nil {
		if config.UnmarshalNode != "" {
			nodes := strings.Split(config.UnmarshalNode, ".")
			regexStr := "((?<=("
			for _, key := range nodes {
				_, err := strconv.Atoi(key)
				if err != nil {
					// TODO: 数组index支持
				}
				regexStr += `(\{|\[).*"` + key + `":`
			}
			regexStr += "))"
			// RJ 2023-03-06 21:12:00 ((?<=((\{|\[).*"data":(\{|\[).*"sub_data":))((\{.+\})|(\[.+\])|(("(?:[^"\\]|\\.)*")|((\d+\.\d+)|(\d+))))((?=,)|(^,.*\}|\])))
			regexStr += jsonNodeRegex
			// RJ 2023-03-06 21:01:15 Go不支持正则的lookarround
			reg := regexp2.MustCompile(regexStr, regexp2.None)
			match, err := reg.FindStringMatch(string(result))
			matchStr := match.String()
			err = json.Unmarshal([]byte(matchStr), &obj)
			if err != nil {
				Logln(shouldLogError, callerLevel, lineLevel, err)
				return err
			}
		} else {
			err = json.Unmarshal(result, &obj)
			if err != nil {
				Logln(shouldLogError, callerLevel, lineLevel, err)
				return err
			}
		}

		err = json.Unmarshal(result, &obj)
		if err != nil {
			Logln(shouldLogError, callerLevel, lineLevel, err)
			return err
		}
		Logln(LogCondition(config.Log&LogObj != 0), callerLevel, lineLevel, obj)
	}
	Logln(LogCondition(config.Log&LogResponse != 0), callerLevel, lineLevel, string(result))

	return nil
}

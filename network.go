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
)

var defaultConfig = &httpConfig{HttpConfig: HttpConfig{Log: LogAll}}

type LogLevel int

// RJ 2022-10-14 11:43:33 日志配置, 默认LogAll
const (
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
	HttpMethodGET    = "GET"
	HttpMethodPOST   = "POST"
	HttpMethodPUT    = "PUT"
	HttpMethodDELETE = "DELETE"
)

type HttpConfig struct {
	Header      http.Header
	Log         LogLevel       // default: LogAll
	LogCaller   LogCallerLevel // 默认为LogCallerLevel(0), 代表请求位置的方法,如果想看tools内部打印所在的方法, 可以传LogCallerLevel(-2)
	LogLine     LogLineLevel   // 默认为LogLineLevel(0), 代表请求位置的行号,如果想看tools内部的打印所在的行, 可以传LogLineLevel(-2)
	Timeout     time.Duration  // 为0会忽略
	ContentType string         // default: "application/json" , post only
}

type httpConfig struct {
	HttpConfig
	Method string
	URL    string
	Params interface{}
	Body   io.Reader
}

/*
*HttpGet
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpGet(urlStr string, values url.Values, obj interface{}, config ...*HttpConfig) error {
	url := urlStr
	var iconfig *httpConfig
	if len(config) > 0 {
		iconfig = &httpConfig{HttpConfig: *config[0]}
	} else {
		iconfig = defaultConfig
	}
	if values != nil {
		url = fmt.Sprintf("%s?%s", urlStr, values.Encode())
	}

	iconfig.Method = HttpMethodPOST
	iconfig.URL = url
	iconfig.Params = values
	err := request(obj, iconfig)
	return err
}

/*
*HttpPost
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpPost(url string, dataMap map[string]string, obj interface{}, config ...*HttpConfig) error {
	var iconfig *httpConfig
	shouldLogError := LogCondition(false)
	if len(config) > 0 {
		iconfig = &httpConfig{HttpConfig: *config[0]}
		if iconfig.Log == LogNil {
			iconfig.Log = LogAll
		}
		shouldLogError = LogCondition(iconfig.Log&LogError != 0)
	} else {
		iconfig = defaultConfig
	}
	jsonParams, err := json.Marshal(dataMap)
	if err != nil {
		Logln(shouldLogError, LogCallerLevel(iconfig.LogCaller+1), LogLineLevel(iconfig.LogLine+1), err)
		return err
	}

	iconfig.URL = url
	iconfig.Method = HttpMethodPOST
	iconfig.Body = bytes.NewBuffer(jsonParams)
	iconfig.Params = dataMap
	if iconfig.ContentType == "" {
		iconfig.ContentType = "application/json"
	}
	err = request(obj, iconfig)

	return err
}

/*
*HttpFormDataPost
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpFormDataPost(url string, dataMap map[string]string, obj interface{}, config ...*HttpConfig) error {
	cmdResReqForm, contentType := createMultipartFormBody(dataMap)
	var iconfig *httpConfig
	if len(config) > 0 {
		iconfig = &httpConfig{HttpConfig: *config[0]}
	} else {
		iconfig = defaultConfig
	}

	iconfig.URL = url
	iconfig.Method = HttpMethodPOST
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

	if config.Method == HttpMethodPOST {
		request.Close = true
		request.Header.Add("Content-Type", config.ContentType)
	}
	if config.Header != nil {
		request.Header = config.Header
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

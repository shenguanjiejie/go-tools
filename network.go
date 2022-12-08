package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	Log         LogLevel
	Timeout     time.Duration
	ContentType string // post only
}

type httpConfig struct {
	HttpConfig
	Method string
	URL    string
	Params interface{}
	Body   io.Reader
}

/**HttpGet
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

/**HttpPost
@param obj : body所序列化的对象, 指针类型, 如果为*http.Response类型, 则直接返回*http.Response
*/
func HttpPost(url string, dataMap map[string]string, obj interface{}, config ...*HttpConfig) error {
	var iconfig *httpConfig
	shouldLogError := Condition(false)
	if len(config) > 0 {
		iconfig = &httpConfig{HttpConfig: *config[0]}
		shouldLogError = Condition(iconfig.Log&LogError != 0)
	} else {
		iconfig = defaultConfig
	}
	jsonParams, err := json.Marshal(dataMap)
	if err != nil {
		Logln(shouldLogError,CallerLevel(1), err)
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

/**HttpFormDataPost
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
	shouldLogError := Condition(config.Log&LogError != 0)
	callerLevel := CallerLevel(2)
	Logln(Condition(config.Log&LogURL != 0), callerLevel, config.Method, config.URL)
	Logln(Condition(config.Log&LogParams != 0), callerLevel, config.Params)

	request, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		Logln(shouldLogError, err)
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
		Logln(shouldLogError, callerLevel, err)
		return err
	}

	if obj != nil && reflect.TypeOf(obj) == reflect.TypeOf(response) {
		*(obj.(*http.Response)) = *response
		Logln(Condition(config.Log&LogResponse != 0), callerLevel, obj)
		return nil
	}
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Logln(shouldLogError, callerLevel, err)
		return err
	}
	if obj != nil {
		err = json.Unmarshal(result, &obj)
		if err != nil {
			Logln(shouldLogError, callerLevel, err)
			return err
		}
		Logln(Condition(config.Log&LogObj != 0), callerLevel, obj)
	}
	Logln(Condition(config.Log&LogResponse != 0), callerLevel, string(result))

	return nil
}

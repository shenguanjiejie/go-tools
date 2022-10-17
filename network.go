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
	"time"
)

var defaultConfig = []*HttpConfig{{Log: LogAll}}

type LogLevel int

// RJ 2022-10-14 11:43:33 日志配置, 默认LogAll
const (
	LogNone LogLevel = 1 << iota
	LogURL
	LogParams
	LogResponse
	LogError
	LogAll = LogURL | LogParams | LogResponse | LogError
)

type HttpConfig struct {
	Header  http.Header
	Log     LogLevel
	Timeout time.Duration
}

func HttpGet(urlStr string, values url.Values, config ...*HttpConfig) ([]byte, error) {
	url := urlStr
	shouldLogError := Condition(false)
	if len(config) > 0 {
		shouldLogError = Condition(config[0].Log&LogError != 0)
	} else {
		config = defaultConfig
	}
	if values != nil {
		url = fmt.Sprintf("%s?%s", urlStr, values.Encode())
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Slogln(shouldLogError, err)
		return nil, err
	}
	client, responseHandle := prepare(url, values, request, config)
	resp, err := client.Do(request)
	if err != nil {
		Slogln(shouldLogError, err)
		return nil, err
	}
	respBytes, _ := responseHandle(resp)
	return respBytes, nil
}

// RJ 2022-03-29 16:22:04 post请求
func HttpPost(url string, dataMap map[string]string, config ...*HttpConfig) ([]byte, error) {
	shouldLogError := Condition(false)
	if len(config) > 0 {
		shouldLogError = Condition(config[0].Log&LogError != 0)
	} else {
		config = defaultConfig
	}
	jsonParams, err := json.Marshal(dataMap)
	if err != nil {
		Slogln(shouldLogError, err)
		return nil, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
	if err != nil {
		Slogln(shouldLogError, err)
		return nil, err
	}
	client, responseHandle := prepare(url, dataMap, request, config)
	request.Close = true
	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		Slogln(shouldLogError, err)
		return nil, err
	}
	respBytes, _ := responseHandle(resp)

	return respBytes, nil
}

// RJ 2022-03-29 16:22:04 formdata请求
func HttpFormDataPost(url string, dataMap map[string]string, config ...*HttpConfig) ([]byte, error) {
	cmdResReqForm, contentType := createMultipartFormBody(dataMap)
	shouldLogError := Condition(false)
	if len(config) > 0 {
		shouldLogError = Condition(config[0].Log&LogError != 0)
	} else {
		config = defaultConfig
	}
	var err error
	if cmdResReqForm == nil {
		Slogln(shouldLogError, err)
		return nil, err
	}
	request, err := http.NewRequest("POST", url, cmdResReqForm)
	if err != nil {
		Slogln(shouldLogError, err)
		return nil, err
	}
	client, responseHandle := prepare(url, dataMap, request, config)
	request.Close = true
	request.Header.Add("Content-Type", contentType)
	resp, err := client.Do(request)
	if err != nil {
		Slogln(shouldLogError, err)
		return nil, err
	}
	respBytes, _ := responseHandle(resp)
	return respBytes, nil
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

func prepare(url string, params interface{}, request *http.Request, config []*HttpConfig) (*http.Client, func(response *http.Response) ([]byte, error)) {
	var conf *HttpConfig
	client := http.DefaultClient
	if len(config) > 0 {
		conf = config[0]
	}

	Slogln(Condition(conf.Log&LogURL != 0), request.Method, url)
	Slogln(Condition(conf.Log&LogParams != 0), params)

	if len(config) > 0 {
		httpConfig := config[0]
		if httpConfig.Header != nil {
			request.Header = config[0].Header
		}
		if httpConfig.Timeout > 0 {
			client.Timeout = httpConfig.Timeout
		}
	}

	return client, func(response *http.Response) ([]byte, error) {
		result, err := ioutil.ReadAll(response.Body)
		if err != nil {
			Slogln(Condition(conf.Log&LogError != 0), CallerLevel(1), err)
			return nil, err
		}
		if conf.Log&LogResponse != 0 {
			Slogln(Condition(conf.Log&LogResponse != 0), CallerLevel(1), string(result))
		}
		return result, nil
	}
}

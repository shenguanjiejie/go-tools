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

func HttpGet(urlStr string, values url.Values, header ...http.Header) ([]byte, error) {
	url := fmt.Sprintf("%s?%s", urlStr, values.Encode())
	Slogln(url)
	request, _ := http.NewRequest("GET", url, nil)
	if len(header) > 0 {
		request.Header = header[0]
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		Slogln(err)
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Slogln(err)
		return nil, err
	}
	return respBytes, nil
}

// RJ 2022-03-29 16:22:04 post请求
func HttpPost(url string, dataMap map[string]string) ([]byte, error) {
	Slogln(url)
	jsonParams, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("Content-Type", "application/json")
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RJ 2022-03-29 16:22:04 formdata请求
func HttpFormDataPost(url string, dataMap map[string]string) ([]byte, error) {
	Slogln(url)
	cmdResReqForm, contentType := createMultipartFormBody(dataMap)
	var err error
	if cmdResReqForm == nil {
		err = fmt.Errorf("%s", "create multipart form body error")
		return nil, err
	}
	req, err := http.NewRequest("POST", url, cmdResReqForm)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("Content-Type", contentType)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
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

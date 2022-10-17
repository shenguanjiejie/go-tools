package tools

import (
	"testing"
)

func TestHttpGet(t *testing.T) {
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, &HttpConfig{Log: LogURL})
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, &HttpConfig{Log: LogURL | LogParams})
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, &HttpConfig{Log: LogAll})
}

func TestHttpPost(t *testing.T) {

	params := map[string]string{
		"userId": "1",
		"id":     "101",
		"title":  "title test",
		"body":   "body test",
	}
	HttpPost("https://jsonplaceholder.typicode.com/posts", params, &HttpConfig{Log: LogAll})
}

func TestHttpFormDataPost(t *testing.T) {

	params := map[string]string{
		"userId": "1",
		"id":     "101",
		"title":  "title test",
		"body":   "body test",
	}
	HttpFormDataPost("https://jsonplaceholder.typicode.com/posts", params, &HttpConfig{Log: LogAll})
}

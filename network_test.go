/*
 * @Author: shenguanjiejie 835166018@qq.com
 * @Date: 2022-10-17 11:51:51
 * @LastEditors: shenguanjiejie 835166018@qq.com
 * @LastEditTime: 2022-10-19 19:15:42
 * @FilePath: /go-tools/network_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tools

import (
	"testing"
)

func TestHttpGet(t *testing.T) {
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, &HttpConfig{Log: LogNone})
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

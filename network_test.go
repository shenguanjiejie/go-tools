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

type User struct {
	ID     int    `json:"id" bson:"id"`
	UserID int    `json:"user_id" bson:"user_id"`
	Title  string `json:"title" bson:"title"`
	Body   string `json:"body" bson:"body"`
}

func TestHttpGet(t *testing.T) {
	user := new(User)
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, nil, &HttpConfig{Log: NetLogNone})
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, nil, &HttpConfig{Log: NetLogURL, LogCaller: LogCallerLevel(-2), LogLine: LogLineLevel(-2)})
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, nil, &HttpConfig{Log: NetLogURL | NetLogParams})
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, nil, &HttpConfig{Log: NetLogAll})
	HttpGet("https://jsonplaceholder.typicode.com/posts/1", nil, user, &HttpConfig{Log: NetLogAll})
}

func TestHttpPost(t *testing.T) {

	user := new(User)
	user.UserID = 1
	user.ID = 101
	user.Title = "title"
	user.Body = "body"

	newUser := new(User)

	// params := map[string]interface{}{
	// 	"userId": "1",
	// 	"id":     "101",
	// 	"title":  "title test",
	// 	"body":   "body test",
	// }
	HttpPost("https://jsonplaceholder.typicode.com/posts", user, newUser, &HttpConfig{Log: NetLogAll})
}

func TestHttpFormDataPost(t *testing.T) {
	user := new(User)
	params := map[string]string{
		"userId": "1",
		"id":     "101",
		"title":  "title test",
		"body":   "body test",
	}
	HttpFormDataPost("https://jsonplaceholder.typicode.com/posts", params, user, &HttpConfig{Log: NetLogAll})
}

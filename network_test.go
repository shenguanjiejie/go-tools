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
	Get("https://jsonplaceholder.typicode.com/posts/1", nil, nil, NetLogLevelOption(NetLogNone))
	Get("https://jsonplaceholder.typicode.com/posts/1", nil, nil, NetLogLevelOption(NetLogURL), LogCallerSkipOption(-2), LogLineSkipOption(-2))
	Get("https://jsonplaceholder.typicode.com/posts/1", nil, nil, NetLogLevelOption(NetLogURL|NetLogParams))
	Get("https://jsonplaceholder.typicode.com/posts/1", nil, nil, NetLogLevelOption(NetLogAll))
	Get("https://jsonplaceholder.typicode.com/posts/1", nil, user)
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
	Post("https://jsonplaceholder.typicode.com/posts", user, newUser)
}

func TestHttpFormDataPost(t *testing.T) {
	user := new(User)
	params := map[string]string{
		"userId": "1",
		"id":     "101",
		"title":  "title test",
		"body":   "body test",
	}
	FormDataPost("https://jsonplaceholder.typicode.com/posts", params, user)
}

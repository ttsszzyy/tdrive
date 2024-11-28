/*
 * Author: 李鸿胤 leeyfann@gmail.com
 * Date: 2023-10-27 14:28:34
 * LastEditors: 李鸿胤 leeyfann@gmail.com
 * Note: Need note condition
 */
package jwt

import (
	"fmt"
	"testing"
)

func TestCreateJwt(t *testing.T) {
	jwt := NewJWT(
		"561528c8-2d58-46df-570a-516bef5b7f89",
		259200,
	)
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDA1NTkwMzQsImlhdCI6MTcwMDQ3MjYzNCwidWlkIjoieHgtYWEtd3ciLCJ6eiI6MTIzNDU2fQ.NVI6-bt3AYKbk21j9lWzeTFZjrnTvm2AquS2QBTa88M
	payload := make(Payload)
	key := "xxx-222"
	payload["uid"] = key
	payload["tag"] = "tag"
	payload["roleName"] = "roleName"
	payload["account"] = "account"
	token, err := jwt.Token(payload)
	//err = jwt.Store(context.Background(), cacheKey(key), token)
	t.Log(token, err)
}

func TestParse(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDIxOTU1NjEsImlhdCI6MTcwMTkzNjM2MSwidWlkIjoiMjIyIiwienoiOjEyMzQ1Nn0.hzW0USKzXU7O_ZhMjR4ofMWkgRU7OxJdrO_cjac5PM8"
	parser := NewJWT("561528c8-2d58-46df-570a-516bef5b7f89", 259200)
	p, err := parser.Parse(token)
	fmt.Println(p, err)
}

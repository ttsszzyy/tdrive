/*
 * @Author: Young
 * @Date: 2022-05-11 10:58:40
 * LastEditors: lihy lihy@zhiannet.com
 * LastEditTime: 2023-10-26 14:34:52
 * @FilePath: /buyday/common/response/reponse.go
 */

package response

import (
	"T-driver/common/errors"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int32       `json:"code"`
	Time int64       `json:"time"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	body := Body{Data: resp, Time: time.Now().Unix()}

	if err != nil {
		errObj, _ := err.(*errors.Error)
		if errObj != nil {
			if errObj.Code() == errors.ErrUnauth {
				httpx.WriteJson(w, http.StatusUnauthorized, nil)
				return
			}
			body.Code = errObj.Code()
			body.Msg = errObj.Msg()
		} else {
			// 未知错误码
			body.Code = errors.ErrUnkonwn
			body.Msg = "系统错误"
		}
	}

	httpx.OkJson(w, body)
}

func ResponseBlob(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		Response(w, resp, err)
	} else {
		httpx.Ok(w)
	}
}

func ResponseRaw(w http.ResponseWriter, resp interface{}) {
	var body []byte
	switch v := resp.(type) {
	case []byte:
		body = v
	case string:
		body = []byte(v)
	default:
		body = []byte("1")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

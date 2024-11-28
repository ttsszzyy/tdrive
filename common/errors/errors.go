/*
 * @Author: Young
 * @Date: 2022-05-05 14:02:50
 * LastEditors: lihy lihy@zhiannet.com
 * LastEditTime: 2024-03-11 16:33:53
 * @FilePath: /buyday/common/errors/errors.go
 */

package errors

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Language 类型别名
type Language = string

const (
	// LanEn 英文
	LanEn Language = "en"
	// LanTw 繁体中文
	LanTw Language = "tw"
)

type Error struct {
	code int32  `json:"errCode"`
	msg  string `json:"errMsg"`
}

func (e *Error) Code() int32 {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Error() string {
	// return fmt.Sprintf("ErrCode:%d ErrMsg:%s", e.code, e.msg)
	val, _ := json.Marshal(e)
	return string(val)
}

// Parse 根据json格式的数据解释为错误对象
// 数据格式：{"errCode":0,"errMsg":"message"}
func Parse(err string) *Error {
	var e Error
	_ = json.Unmarshal([]byte(err), &e)
	return &e
}

func FromRpcError(errRpc error) *Error {
	s, _ := status.FromError(errRpc)
	if s == nil {
		return NewErrCodeMsg(ErrUnkonwn, "未知错误")
	}
	/*
		e := &Error{}
		err := json.Unmarshal([]byte(s.Message()), e)
		if err != nil {
			return NewErrCodeMsg(ErrUnkonwn, "未知错误")
		}
	*/
	return NewErrCodeMsg(int32(s.Code()), s.Message())
}

func NewErrCodeMsg(code int32, msg string) *Error {
	return &Error{code: code, msg: msg}
}

func NewErrCode(code int32, lans ...Language) *Error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}

	switch lan {
	case LanEn:
		msg = "unkonwn error message"
	case LanTw:
		msg = "未知錯誤"
	default:
		msg = "unkonwn error message"
	}
	return &Error{code: code, msg: msg}
}

func SystemError(lans ...Language) error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}
	switch lan {
	case LanEn:
		msg = "System error, please try again later"
	case LanTw:
		msg = "繫統錯誤，請稍後重試"
	default:
		msg = "System error, please try again later"
	}
	return NewErrCodeMsg(ErrSystem, msg)
}

func CacheError(lans ...Language) *Error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}
	switch lan {
	case LanEn:
		msg = "Cache error"
	case LanTw:
		msg = "緩存錯誤"
	default:
		msg = "Cache error"
	}
	return NewErrCodeMsg(ErrDB, msg)
}

func RpcSystemError(lans ...Language) error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}
	switch lan {
	case LanEn:
		msg = "System error, please try again later"
	case LanTw:
		msg = "繫統錯誤，請稍後重試"
	default:
		msg = "System error, please try again later"
	}
	return NewErrCodeMsg(ErrSystem, msg).RpcError()
}

func PermissionError(lans ...Language) error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}
	switch lan {
	case LanEn:
		msg = "Insufficient permissions"
	case LanTw:
		msg = "權限不足"
	default:
		msg = "Insufficient permissions"
	}
	return NewErrCodeMsg(ErrPermission, msg)
}

func ParamsError(lans ...Language) error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}
	switch lan {
	case LanEn:
		msg = "Parameter error"
	case LanTw:
		msg = "參數錯誤"
	default:
		msg = "Parameter error"
	}
	return NewErrCodeMsg(ErrCodeReqParam, msg)
}

func UnauthError(lans ...Language) error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}
	switch lan {
	case LanEn:
		msg = "Unauthorized login"
	case LanTw:
		msg = "未授權登錄"
	default:
		msg = "Unauthorized login"
	}
	return NewErrCodeMsg(401, msg)
}

func DbError(lans ...Language) *Error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}
	switch lan {
	case LanEn:
		msg = "Database error"
	case LanTw:
		msg = "數據庫錯誤"
	default:
		msg = "Database error"
	}
	return NewErrCodeMsg(ErrDB, msg)
}

func ErrorNotFound(lans ...Language) *Error {
	var msg, lan string

	if len(lans) > 0 {
		lan = lans[0]
	}
	switch lan {
	case LanEn:
		msg = "No records found"
	case LanTw:
		msg = "未查詢到記錄"
	default:
		msg = "No records found"
	}
	return NewErrCodeMsg(ErrCodeNotFount, msg)
}

func CustomError(msg string) *Error {
	return NewErrCodeMsg(ErrCodeCustom, msg)
}
func CustomErrorf(format string, msg ...any) *Error {
	return NewErrCodeMsg(ErrCodeCustom, fmt.Sprintf(format, msg...))
}

func (e *Error) RpcError() error {
	return status.Error(codes.Aborted, e.Msg())
}

func GetError(code int32, lans ...Language) *Error {
	var msg, lan, defaultMsg string

	if len(lans) > 0 {
		lan = lans[0]
	}

	switch lan {
	case LanEn:
		msg = errMsgEn[code]
		defaultMsg = "unkonwn error message"
	case LanTw:
		msg = errMsgTw[code]
		defaultMsg = "未知錯誤"
	default:
		msg = errMsgEn[code]
		defaultMsg = "unkonwn error message"
	}

	if msg == "" {
		msg = defaultMsg
	}

	return NewErrCodeMsg(code, msg)
}

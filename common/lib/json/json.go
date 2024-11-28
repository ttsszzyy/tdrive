/*
 * @Author: Young
 * @Date: 2021-05-26 20:04:56
 * LastEditors: lihy lihy@zhiannet.com
 * LastEditTime: 2023-03-15 20:55:06
 * @FilePath: /meta-nft/app/user/api/internal/lib/json/json.go
 */
package json

import (
	"bytes"
	"io"

	jsoniter "github.com/json-iterator/go"
)

var jsonLib = jsoniter.ConfigCompatibleWithStandardLibrary

func Unmarshal(data []byte, v interface{}) error {
	return jsonLib.Unmarshal(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return jsonLib.Marshal(v)
}

func MarshalReader(v interface{}) (io.Reader, error) {
	switch v.(type) {
	case string:
		return bytes.NewReader([]byte(v.(string))), nil
	default:

	}
	data, err := jsonLib.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

type Bint64 int64

func (B *Bint64) UnmarshalJSON(b []byte) error {
	txt := string(bytes.Trim(b, `"`))
	if txt == "0" || txt == "false" {
		*B = 0
	} else {
		*B = 1
	}
	return nil
}

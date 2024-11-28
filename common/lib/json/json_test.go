package json

import (
	"encoding/json"
	"testing"

	"github.com/jinzhu/copier"
)

type DeviceSoftwareInfoB struct {
	SoftwareName string `json:"NAME"`    //  软件名称
	Version      string `json:"Version"` //  软件版本
	Size         string `json:"Size"`    //  软件大小 如 17.5MB
	State        Bint64 `json:"State"`   //  运行状态 0未运行 1运行中
}

type DeviceSoftwareInfo struct {
	SoftwareName string `json:"Name"`    //  软件名称
	Version      string `json:"version"` //  软件版本
	Size         string `json:"size"`    //  软件大小 如 17.5MB
	State        int64  `json:"state"`   //  运行状态 0未运行 1运行中
}

func TestUnmarshal(t *testing.T) {
	str := "[{\"Name\":\"微信\",\"Version\":\"3.8.1.26\",\"State\":true},{\"Name\":\"钉钉\",\"Version\":\"7.0.10-Release.2189102\"}]"
	list := make([]DeviceSoftwareInfoB, 0)
	resp := make([]DeviceSoftwareInfo, 0)
	if err := json.Unmarshal([]byte(str), &list); err != nil {
		t.Fatal(err)
	} else {
		for _, v := range list {
			t.Logf("%+v", v)
		}
	}
	if err := copier.Copy(&resp, list); err != nil {
		t.Fatal(err)
	} else {
		for _, v := range resp {
			t.Logf("%+v", v)
		}
	}
}

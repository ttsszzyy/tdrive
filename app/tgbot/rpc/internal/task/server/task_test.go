package server

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

func TestName(t *testing.T) {
	// 定义模板字符串，其中包含占位符
	tmpl := "Hello,Today is {{.Day}}"

	// 创建一个新的模板
	parse, err := template.New("example").Parse(tmpl)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// 定义要替换的数据
	data := map[string]string{
		"Name": "Hi，{{.username}}\n欢迎来到TDrive\nTDrive是一个去中心化网盘\n原生于TON网络\n您可以直接转发Telegram内的视频和图片，就可以保存到T-Drive网盘，这样您可以不需要下载就快速查看；\n还可以获取积分和空投奖励\n⭐快嘗試一下吧！",
		"Day":  "Monday",
	}
	// 执行模板，并传入数据
	var out bytes.Buffer
	err = parse.Execute(&out, data)
	result := out.String()
	fmt.Println("Result:", result)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}

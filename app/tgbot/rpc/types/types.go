package types

type Markup struct {
	Button       string `json:"button,optional"`         //按钮名称
	Url          string `json:"url,optional"`            //按钮url
	CallbackData string `json:"callback_data,omitempty"` //回调数据
}

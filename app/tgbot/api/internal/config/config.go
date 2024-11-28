package config

import (
	"T-driver/common/asynq"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	AsynqConf asynq.AsynqConf
	Telegram  struct {
		Url        string
		Token      string
		WebhookUrl string
	}
	TiTan struct {
		TitanURL  string
		APIKey    string
		APISecret string
	}
}

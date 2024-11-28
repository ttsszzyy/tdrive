package config

import (
	"T-driver/common/asynq"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	BaseUrl          string
	RedisConf        redis.RedisConf
	AsynqConf        asynq.AsynqConf
	MaxUpload        int64
	UploadExpireTime int
	FastReward       struct {
		Storage      int64
		Distribution int
		Integral     int64
		TgUrl        string
		AssetUrl     string
		RecoveryDate int
		Describe     string
	}
	Telegram struct {
		Token       string
		WebhookUrl  string
		WebhookPort string
	}
	TiTan struct {
		TitanURL string
		APIKey   string
	}
	Action struct {
		Name     string
		Lon      float64
		Lat      float64
		Distance float64
		Points   int64
	}
	UserRpc zrpc.RpcClientConf
	BotRpc  zrpc.RpcClientConf
}

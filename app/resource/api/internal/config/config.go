package config

import (
	"T-driver/common/asynq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	RedisConf        redis.RedisConf
	AsynqConf        asynq.AsynqConf
	MaxUpload        int64
	UploadExpireTime int
	Telegram         struct {
		Token string
	}
	TiTan struct {
		TitanURL string
		APIKey   string
		AreaID   string
	}
	UserRpc zrpc.RpcClientConf
}

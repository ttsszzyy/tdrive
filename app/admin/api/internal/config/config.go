package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	RedisConf redis.RedisConf
	AuthCfg   struct {
		AccessSecret string
		AccessExpire int64
	}
	TiTan struct {
		TitanURL string
		APIKey   string
		AreaID   string
	}
	AwsS3 struct {
		AccessKey string
		SecretKey string
		EndPoint  string
		Bucket    string
		Region    string
	}
	UserRpc  zrpc.RpcClientConf
	BotRpc   zrpc.RpcClientConf
	VideoRpc zrpc.RpcClientConf
}

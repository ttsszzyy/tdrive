package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Telegram struct {
		Token string
	}
	RedisConf redis.RedisConf
	UserRpc   zrpc.RpcClientConf
	VideoRpc  zrpc.RpcClientConf

	S3 struct {
		Endpoint string
		CDN      string
	}
}

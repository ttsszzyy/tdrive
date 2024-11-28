package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		Dns     string
		DbCache cache.CacheConf
	}
	MonDb struct {
		Uri string
		Db  string
	}
	RedisConf redis.RedisConf
}

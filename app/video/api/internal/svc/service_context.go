package svc

import (
	"T-driver/app/user/rpc/user"
	"T-driver/app/video/api/internal/config"
	"T-driver/app/video/api/internal/middleware"
	"T-driver/app/video/rpc/video"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	Redis          *redis.Redis
	UserRpc        user.UserZrpcClient
	VideoRpc       video.VideoZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCli := redis.MustNewRedis(c.RedisConf)
	return &ServiceContext{
		Config:         c,
		Redis:          redisCli,
		AuthMiddleware: middleware.NewAuthMiddleware(c.Telegram.Token).Handle,
		UserRpc:        user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc:       video.NewVideoZrpcClient(zrpc.MustNewClient(c.VideoRpc)),
	}
}

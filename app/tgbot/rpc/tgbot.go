package main

import (
	"T-driver/app/tgbot/rpc/internal/task/server"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"T-driver/app/tgbot/rpc/internal/config"
	zrpcserver "T-driver/app/tgbot/rpc/internal/server"
	"T-driver/app/tgbot/rpc/internal/svc"
	"T-driver/app/tgbot/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/tgbot.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterTgbotServer(grpcServer, zrpcserver.NewTgbotServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	go func() {
		s.Start()
	}()
	defer s.Stop()

	serviceGroup := service.NewServiceGroup()
	defer func() {
		serviceGroup.Stop()
		logx.Close()
	}()
	serviceGroup.Add(server.NewTaskServer(ctx))
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	serviceGroup.Start()
}

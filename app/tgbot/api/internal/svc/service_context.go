package svc

import (
	"T-driver/app/tgbot/api/internal/config"
	"github.com/go-telegram/bot"
	"github.com/hibiken/asynq"
	storage "github.com/utopiosphe/titan-storage-sdk"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config      config.Config
	TgBot       *bot.Bot
	AsynqClient *asynq.Client
	Tenant      storage.Tenant
}

func NewServiceContext(c config.Config) *ServiceContext {
	tenant, err := storage.NewTenant(c.TiTan.TitanURL, c.TiTan.APIKey)
	if err != nil {
		logx.Must(err)
	}
	svcCtx := &ServiceContext{
		Config:      c,
		Tenant:      tenant,
		AsynqClient: c.AsynqConf.NewClient(),
	}

	return svcCtx
}

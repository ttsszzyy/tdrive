package bot

import (
	"T-driver/app/tgbot/api/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type BotWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录
func NewBotWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BotWebhookLogic {
	return &BotWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BotWebhookLogic) BotWebhook() error {

	return nil
}

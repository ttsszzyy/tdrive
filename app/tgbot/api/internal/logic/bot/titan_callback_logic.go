package bot

import (
	"context"

	"T-driver/app/tgbot/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type TitanCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTitanCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TitanCallbackLogic {
	return &TitanCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TitanCallbackLogic) TitanCallback() error {
	// todo: add your logic here and delete this line

	return nil
}

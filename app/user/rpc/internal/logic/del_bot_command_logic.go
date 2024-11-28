package logic

import (
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelBotCommandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelBotCommandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelBotCommandLogic {
	return &DelBotCommandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除bot命令
func (l *DelBotCommandLogic) DelBotCommand(in *pb.DelBotCommandReq) (*pb.Response, error) {
	err := l.svcCtx.BotCommandModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Response{}, nil
}

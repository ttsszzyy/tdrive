package parameter

import (
	"T-driver/app/tgbot/rpc/pb"
	"T-driver/app/user/rpc/user"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelBotCommandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelBotCommandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelBotCommandLogic {
	return &DelBotCommandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelBotCommandLogic) DelBotCommand(req *types.DelBotCommandReq) (resp *types.Response, err error) {
	command, err := l.svcCtx.Rpc.FindBotCommand(l.ctx, &user.FindBotCommandReq{})
	if err != nil {
		return nil, err
	}
	botCommand := make([]*pb.BotCommand, 0, len(command.List))
	botCommandMap := make(map[string]string)
	for _, v := range command.List {
		if v.Id == req.Id || v.SendType == 1 {
			continue
		}
		botCommandMap[v.Command] = v.Description
	}
	for k, v := range botCommandMap {
		botCommand = append(botCommand, &pb.BotCommand{
			Command:     k,
			Description: v,
		})
	}
	sendBotCommand, err := l.svcCtx.BotRpc.SendBotCommand(l.ctx, &pb.SendBotCommandReq{BotCommand: botCommand})
	if err != nil {
		return nil, err
	}
	if sendBotCommand.IsSuccess {
		_, err = l.svcCtx.Rpc.DelBotCommand(l.ctx, &user.DelBotCommandReq{Id: req.Id})
		if err != nil {
			return nil, err
		}
	}
	return &types.Response{}, nil
}

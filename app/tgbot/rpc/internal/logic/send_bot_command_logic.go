package logic

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"T-driver/app/tgbot/rpc/internal/svc"
	"T-driver/app/tgbot/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendBotCommandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendBotCommandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendBotCommandLogic {
	return &SendBotCommandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendBotCommandLogic) SendBotCommand(in *pb.SendBotCommandReq) (*pb.DelPinMsgResp, error) {
	commands := make([]models.BotCommand, 0, len(in.BotCommand))
	for _, v := range in.BotCommand {
		commands = append(commands, models.BotCommand{
			Command:     v.Command,
			Description: v.Description,
		})
	}
	myCommands, err := l.svcCtx.TgBot.SetMyCommands(l.ctx, &bot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DelPinMsgResp{
		IsSuccess: myCommands,
	}, nil
}

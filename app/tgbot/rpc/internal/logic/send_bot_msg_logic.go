package logic

import (
	"context"
	"github.com/go-telegram/bot"

	"T-driver/app/tgbot/rpc/internal/svc"
	"T-driver/app/tgbot/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendBotMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendBotMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendBotMsgLogic {
	return &SendBotMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 文本发送
func (l *SendBotMsgLogic) SendBotMsg(in *pb.SendBotMsgRequest) (*pb.SendBotMsgResponse, error) {
	_, err := l.svcCtx.TgBot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: in.ChatID,
		Text:   in.Text,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SendBotMsgResponse{}, nil
}

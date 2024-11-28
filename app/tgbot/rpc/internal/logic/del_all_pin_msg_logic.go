package logic

import (
	"context"
	"github.com/go-telegram/bot"

	"T-driver/app/tgbot/rpc/internal/svc"
	"T-driver/app/tgbot/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelAllPinMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelAllPinMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAllPinMsgLogic {
	return &DelAllPinMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelAllPinMsgLogic) DelAllPinMsg(in *pb.DelAllPinMsgReq) (*pb.DelPinMsgResp, error) {
	msg, err := l.svcCtx.TgBot.UnpinAllChatMessages(l.ctx, &bot.UnpinAllChatMessagesParams{
		ChatID: in.ChatID,
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	return &pb.DelPinMsgResp{
		IsSuccess: msg,
	}, nil
}

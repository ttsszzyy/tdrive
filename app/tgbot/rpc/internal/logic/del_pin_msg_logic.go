package logic

import (
	"T-driver/app/tgbot/rpc/internal/svc"
	"T-driver/app/tgbot/rpc/pb"
	"context"
	"github.com/go-telegram/bot"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelPinMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelPinMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelPinMsgLogic {
	return &DelPinMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelPinMsgLogic) DelPinMsg(in *pb.DelPinMsgReq) (*pb.DelPinMsgResp, error) {
	msg, err := l.svcCtx.TgBot.UnpinChatMessage(l.ctx, &bot.UnpinChatMessageParams{
		ChatID:    in.ChatID,
		MessageID: int(in.MsgID),
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	return &pb.DelPinMsgResp{
		IsSuccess: msg,
	}, nil
}

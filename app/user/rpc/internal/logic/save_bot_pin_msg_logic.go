package logic

import (
	"T-driver/app/user/model"
	"context"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveBotPinMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveBotPinMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveBotPinMsgLogic {
	return &SaveBotPinMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加pin消息
func (l *SaveBotPinMsgLogic) SaveBotPinMsg(in *pb.SaveBotPinMsgReq) (*pb.Response, error) {
	_, err := l.svcCtx.BotPinMessageModel.Insert(l.ctx, &model.BotPinMessage{
		ChatId:      in.ChatId,
		Message:     in.Msg,
		CreatedTime: time.Now().Unix(),
		Text:        in.Text,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Response{}, nil
}

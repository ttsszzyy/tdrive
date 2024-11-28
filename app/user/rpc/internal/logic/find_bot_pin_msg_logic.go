package logic

import (
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"
	"github.com/Masterminds/squirrel"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindBotPinMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindBotPinMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindBotPinMsgLogic {
	return &FindBotPinMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询pin消息
func (l *FindBotPinMsgLogic) FindBotPinMsg(in *pb.FindBotPinMsgReq) (*pb.FindBotPinMsgResp, error) {
	sb := squirrel.Select().Where("chat_id = ?", in.ChatId).Where("deleted_time = ?", 0)
	list, total, err := l.svcCtx.BotPinMessageModel.ListPage(l.ctx, in.Page, in.Size, sb)
	if err != nil {
		return nil, err
	}
	botPinMsgs := make([]*pb.BotPinMsg, 0, len(list))
	for _, v := range list {
		botPinMsgs = append(botPinMsgs, &pb.BotPinMsg{
			Id:          v.Id,
			ChatId:      v.ChatId,
			Msg:         v.Message,
			CreatedTime: v.CreatedTime,
			Text:        v.Text,
		})
	}
	return &pb.FindBotPinMsgResp{
		List:  botPinMsgs,
		Total: total,
	}, nil
}

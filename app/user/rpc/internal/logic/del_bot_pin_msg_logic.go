package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelBotPinMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelBotPinMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelBotPinMsgLogic {
	return &DelBotPinMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除pin消息
func (l *DelBotPinMsgLogic) DelBotPinMsg(in *pb.DelBotPinMsgReq) (*pb.Response, error) {
	if in.Id > 0 {
		err := l.svcCtx.BotPinMessageModel.Delete(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
	} else {
		list, err := l.svcCtx.BotPinMessageModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"chat_id": in.ChatId}).Where("deleted_time = ?", 0))
		if err != nil {
			return nil, err
		}
		ids := make([]int64, 0, len(list))
		for _, v := range list {
			ids = append(ids, v.Id)
		}
		err = l.svcCtx.BotPinMessageModel.Deletes(l.ctx, ids)
		if err != nil {
			return nil, err
		}
	}

	return &pb.Response{}, nil
}

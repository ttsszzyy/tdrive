package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindMessageLogic {
	return &FindMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询消息
func (l *FindMessageLogic) FindMessage(in *pb.FindMessageReq) (*pb.FindMessageResp, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Puid > 0 {
		sb = sb.Where("puid = ?", in.Puid)
	}
	messages, total, err := l.svcCtx.MessageModel.ListPage(l.ctx, in.Page, in.Size, sb)
	if err != nil {
		return nil, err
	}
	pbMessages := make([]*pb.MessageItem, 0)
	for _, message := range messages {
		pbMessages = append(pbMessages, &pb.MessageItem{
			Id:          message.Id,
			Puid:        message.Puid,
			Name:        message.Name,
			Uid:         message.Uid,
			Status:      message.Status,
			Remark:      message.Remark,
			CreatedTime: message.CreatedTime,
			UpdatedTime: message.UpdatedTime,
		})
	}
	return &pb.FindMessageResp{
		Messages: pbMessages,
		Total:    total,
	}, nil
}

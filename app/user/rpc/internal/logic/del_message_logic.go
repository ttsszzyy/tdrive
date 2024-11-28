package logic

import (
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelMessageLogic {
	return &DelMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除消息
func (l *DelMessageLogic) DelMessage(in *pb.DelMessageReq) (*pb.Response, error) {
	err := l.svcCtx.MessageModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Response{}, nil
}

package logic

import (
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelShareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelShareLogic {
	return &DelShareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除分享资源
func (l *DelShareLogic) DelShare(in *pb.FindOneShareReq) (*pb.Response, error) {
	err := l.svcCtx.ShareModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Response{}, nil
}

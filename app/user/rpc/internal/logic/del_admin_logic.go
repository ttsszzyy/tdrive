package logic

import (
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAdminLogic {
	return &DelAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存管理端用户
func (l *DelAdminLogic) DelAdmin(in *pb.DelTaskPoolByIdReq) (*pb.Response, error) {
	if in.Id > 0 {
		err := l.svcCtx.AdminModel.Delete(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
	}

	return &pb.Response{}, nil
}

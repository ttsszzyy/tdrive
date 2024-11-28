package admin

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAdminLogic {
	return &DelAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelAdminLogic) DelAdmin(req *types.DeleteTaskPoolReq) (resp *types.Response, err error) {
	resp = &types.Response{}
	time, err := l.svcCtx.Rpc.FindOneByIdAccountDeletedTime(l.ctx, &pb.AccountReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	if time.Account == "admin" {
		return nil, errors.CustomError("Admin user is prohibited from deleting")
	}
	_, err = l.svcCtx.Rpc.DelAdmin(l.ctx, &pb.DelTaskPoolByIdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

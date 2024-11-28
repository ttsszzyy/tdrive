package admin

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminLogic {
	return &UpdateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdminLogic) UpdateAdmin(req *types.UpdateAdminReq) (resp *types.Response, err error) {
	resp = &types.Response{}
	//检验用户名
	one, err := l.svcCtx.Rpc.FindOneByAccountDeletedTime(l.ctx, &pb.AccountReq{Account: req.Account})
	if err != nil {
		if err != model.ErrNotFound {
			return nil, err
		}
	}
	if one.Id != req.Id {
		return nil, errors.CustomError("Duplicate account name")
	}

	time, err := l.svcCtx.Rpc.FindOneByIdAccountDeletedTime(l.ctx, &pb.AccountReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.Rpc.SaveAdmin(l.ctx, &pb.Admin{
		Id:          req.Id,
		Account:     req.Account,
		Password:    time.Password,
		Avatar:      time.Avatar,
		IsDisable:   req.IsDisable,
		CreatedTime: time.CreatedTime,
		LastTime:    time.LastTime,
		Remark:      req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

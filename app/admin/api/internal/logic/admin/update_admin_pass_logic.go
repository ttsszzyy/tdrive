package admin

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"T-driver/common/utils"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminPassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminPassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminPassLogic {
	return &UpdateAdminPassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdminPassLogic) UpdateAdminPass(req *types.UpdateAdminPassReq) (resp *types.Response, err error) {
	resp = &types.Response{}
	time, err := l.svcCtx.Rpc.FindOneByIdAccountDeletedTime(l.ctx, &pb.AccountReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	if !utils.CheckPassword(req.Password, time.Password, model.PassSalt) {
		return nil, errors.CustomError("Password error")
	}
	_, err = l.svcCtx.Rpc.SaveAdmin(l.ctx, &pb.Admin{
		Id:          req.Id,
		Account:     time.Account,
		Password:    utils.GenPassword(req.Password, model.PassSalt),
		Avatar:      time.Avatar,
		IsDisable:   time.IsDisable,
		CreatedTime: time.CreatedTime,
		LastTime:    time.LastTime,
		Remark:      time.Remark,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

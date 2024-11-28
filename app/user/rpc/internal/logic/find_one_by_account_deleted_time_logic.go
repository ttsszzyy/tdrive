package logic

import (
	"T-driver/app/user/model"
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneByAccountDeletedTimeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneByAccountDeletedTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneByAccountDeletedTimeLogic {
	return &FindOneByAccountDeletedTimeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOneByAccountDeletedTimeLogic) FindOneByAccountDeletedTime(in *pb.AccountReq) (*pb.Admin, error) {
	admin, err := l.svcCtx.AdminModel.FindOneByAccountDeletedTime(l.ctx, in.Account, 0)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.Admin{}, nil
		}
		return nil, err
	}
	return &pb.Admin{
		Id:          admin.Id,
		Account:     admin.Account,
		Password:    admin.Password,
		Avatar:      admin.Avatar,
		IsDisable:   admin.IsDisable,
		CreatedTime: admin.CreatedTime,
		UpdatedTime: admin.UpdatedTime,
		LastTime:    admin.LastTime,
		Remark:      admin.Remark,
	}, err
}

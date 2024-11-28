package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneByIdAccountDeletedTimeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneByIdAccountDeletedTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneByIdAccountDeletedTimeLogic {
	return &FindOneByIdAccountDeletedTimeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理端
func (l *FindOneByIdAccountDeletedTimeLogic) FindOneByIdAccountDeletedTime(in *pb.AccountReq) (*pb.Admin, error) {
	admin, err := l.svcCtx.AdminModel.FindOne(l.ctx, in.Id)
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

package logic

import (
	"T-driver/app/user/model"
	"T-driver/common/errors"
	"context"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAdminLogic {
	return &SaveAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存管理端用户
func (l *SaveAdminLogic) SaveAdmin(in *pb.Admin) (*pb.Response, error) {
	if in.Id > 0 {
		findOne := &model.Admin{
			Id:          in.Id,
			Account:     in.Account,
			Password:    in.Password,
			Avatar:      in.Avatar,
			IsDisable:   in.IsDisable,
			Remark:      in.Remark,
			LastTime:    in.LastTime,
			CreatedTime: in.CreatedTime,
			UpdatedTime: time.Now().Unix(),
		}
		err := l.svcCtx.AdminModel.Update(l.ctx, findOne)
		if err != nil {
			return nil, errors.DbError()
		}
	} else {
		_, err := l.svcCtx.AdminModel.Insert(l.ctx, &model.Admin{
			Account:     in.Account,
			Password:    in.Password,
			Avatar:      in.Avatar,
			Remark:      in.Remark,
			LastTime:    in.LastTime,
			IsDisable:   in.IsDisable,
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		})
		if err != nil {
			return nil, errors.DbError()
		}
	}

	return &pb.Response{}, nil
}

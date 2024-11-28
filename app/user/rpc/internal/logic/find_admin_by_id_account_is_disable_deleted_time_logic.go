package logic

import (
	"T-driver/app/user/model"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAdminByIdAccountIsDisableDeletedTimeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAdminByIdAccountIsDisableDeletedTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAdminByIdAccountIsDisableDeletedTimeLogic {
	return &FindAdminByIdAccountIsDisableDeletedTimeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理端列表
func (l *FindAdminByIdAccountIsDisableDeletedTimeLogic) FindAdminByIdAccountIsDisableDeletedTime(in *pb.AdminReq) (*pb.QueryAdminRes, error) {
	sb := squirrel.Select().Where("deleted_time =?", 0)
	if in.Account != "" {
		sb = sb.Where(squirrel.Like{"account": fmt.Sprintf("%s%%", in.Account)})
	}
	if len(in.Id) > 0 {
		sb = sb.Where(squirrel.Eq{"id": in.Id})
	}
	if in.IsDisable > 0 {
		sb = sb.Where(squirrel.Eq{"is_disable": in.IsDisable})
	}
	var (
		list  []*model.Admin
		total int64
		err   error
	)
	if in.Size == 0 && in.Page == 0 {
		list, err = l.svcCtx.AdminModel.List(l.ctx, sb.OrderBy("created_time desc"))
	} else {
		list, total, err = l.svcCtx.AdminModel.ListPage(l.ctx, in.Page, in.Size, sb.OrderBy("created_time desc"))
	}
	if err != nil {
		return nil, err
	}
	resp := &pb.QueryAdminRes{
		Admins: make([]*pb.Admin, 0, len(list)),
		Total:  total,
	}
	for _, v := range list {
		resp.Admins = append(resp.Admins, &pb.Admin{
			Id:          v.Id,
			Account:     v.Account,
			Password:    v.Password,
			Avatar:      v.Avatar,
			IsDisable:   v.IsDisable,
			CreatedTime: v.CreatedTime,
			UpdatedTime: v.UpdatedTime,
			LastTime:    v.LastTime,
			Remark:      v.Remark,
		})
	}

	return resp, nil
}

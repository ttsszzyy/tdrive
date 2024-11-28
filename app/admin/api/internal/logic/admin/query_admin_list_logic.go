package admin

import (
	"T-driver/app/user/rpc/pb"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryAdminListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryAdminListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAdminListLogic {
	return &QueryAdminListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryAdminListLogic) QueryAdminList(req *types.QueryAdminListReq) (resp *types.QueryAdminListRes, err error) {
	list, err := l.svcCtx.Rpc.FindAdminByIdAccountIsDisableDeletedTime(l.ctx, &pb.AdminReq{
		Id:        []int64{req.Id},
		Account:   req.Account,
		IsDisable: req.IsDisable,
		Page:      req.Page,
		Size:      req.Size,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.QueryAdminListRes{
		Admins: make([]*types.Admin, 0, len(list.Admins)),
		Total:  list.Total,
	}
	for _, v := range list.Admins {
		resp.Admins = append(resp.Admins, &types.Admin{
			Id:          v.Id,
			Account:     v.Account,
			Remark:      v.Remark,
			IsDisable:   v.IsDisable,
			CreatedTime: v.CreatedTime,
			LastTime:    v.LastTime,
		})
	}

	return resp, nil
}

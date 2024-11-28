package user

import (
	"T-driver/app/user/rpc/user"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.IdReq) (resp *types.Response, err error) {
	if req.Id > 0 {
		_, err = l.svcCtx.Rpc.DisableUser(l.ctx, &user.DisableUserReq{
			Id:        req.Id,
			IsDisable: 1,
		})
		if err != nil {
			return nil, err
		}
	}
	return &types.Response{}, nil
}

package admin

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/utils"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminLogic {
	return &AddAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAdminLogic) AddAdmin(req *types.AddAdminReq) (resp *types.Response, err error) {
	resp = &types.Response{}
	//检验用户名
	_, err = l.svcCtx.Rpc.FindOneByAccountDeletedTime(l.ctx, &pb.AccountReq{Account: req.Account})
	if err != nil {
		if err != model.ErrNotFound {
			return nil, err
		}
	}
	_, err = l.svcCtx.Rpc.SaveAdmin(l.ctx, &pb.Admin{
		Account:   req.Account,
		Password:  utils.GenPassword(req.Password, model.PassSalt),
		IsDisable: req.IsDisable,
		Remark:    req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

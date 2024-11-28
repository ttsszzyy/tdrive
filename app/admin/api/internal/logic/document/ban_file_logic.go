package document

import (
	"T-driver/app/user/rpc/user"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BanFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBanFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanFileLogic {
	return &BanFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BanFileLogic) BanFile(req *types.BanFileReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.Rpc.UpdateAssetsName(l.ctx, &user.UpdateAssetsNameReq{
		Id:     req.Id,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return
}

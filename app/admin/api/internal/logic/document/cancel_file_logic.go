package document

import (
	"T-driver/app/user/rpc/user"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelFileLogic {
	return &CancelFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelFileLogic) CancelFile(req *types.CancelFileReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.Rpc.UpdateAssetsName(l.ctx, &user.UpdateAssetsNameReq{
		Id:       req.Id,
		IsReport: 2,
	})
	if err != nil {
		return nil, err
	}
	return &types.Response{}, nil
}

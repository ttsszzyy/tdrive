package taskPool

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTaskPoolLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTaskPoolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskPoolLogic {
	return &DeleteTaskPoolLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTaskPoolLogic) DeleteTaskPool(req *types.DeleteTaskPoolReq) (resp *types.Response, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	_, err = l.svcCtx.Rpc.DelTaskPoolById(l.ctx, &pb.DelTaskPoolByIdReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, errors.DbError(lan)
	}
	return &types.Response{}, nil
}

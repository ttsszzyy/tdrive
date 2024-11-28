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

type UpdateTaskPoolLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskPoolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskPoolLogic {
	return &UpdateTaskPoolLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTaskPoolLogic) UpdateTaskPool(req *types.UpdateTaskPoolReq) (resp *types.Response, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	_, err = l.svcCtx.Rpc.SaveTaskPool(l.ctx, &pb.SaveTaskPoolReq{
		Id:        req.Id,
		TaskType:  req.TaskType,
		TaskName:  req.TaskName,
		Integral:  req.Integral,
		IsDisable: req.IsDisable,
	})
	if err != nil {
		return nil, errors.DbError(lan)
	}
	return &types.Response{}, nil
}

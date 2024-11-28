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

type AddTaskPoolLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddTaskPoolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTaskPoolLogic {
	return &AddTaskPoolLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTaskPoolLogic) AddTaskPool(req *types.AddTaskPoolReq) (resp *types.Response, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	_, err = l.svcCtx.Rpc.SaveTaskPool(l.ctx, &pb.SaveTaskPoolReq{
		TaskType:  req.TaskType,
		TaskName:  req.TaskName,
		Integral:  req.Integral,
		IsDisable: 2,
	})
	if err != nil {
		return nil, errors.DbError(lan)
	}
	return &types.Response{}, nil
}

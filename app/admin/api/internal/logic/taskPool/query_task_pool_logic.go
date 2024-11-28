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

type QueryTaskPoolLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryTaskPoolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTaskPoolLogic {
	return &QueryTaskPoolLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTaskPoolLogic) QueryTaskPool(req *types.QueryTaskPoolReq) (resp *types.QueryTaskPoolRes, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	TaskPoolListResp, err := l.svcCtx.Rpc.FindTaskPoolByTaskType(l.ctx, &pb.FindTaskPoolByTaskTypeReq{
		TaskType: req.TaskType,
		TaskName: req.TaskName,
		Page:     req.Page,
		Size:     req.Size,
	})
	if err != nil {
		return nil, errors.DbError(lan)
	}
	resp = &types.QueryTaskPoolRes{
		TaskPools: []*types.TaskPool{},
		Count:     TaskPoolListResp.Count,
	}
	for _, pool := range TaskPoolListResp.TaskPools {
		resp.TaskPools = append(resp.TaskPools, &types.TaskPool{
			Id:          pool.Id,
			TaskType:    pool.TaskType,
			TaskName:    pool.TaskName,
			Integral:    pool.Integral,
			CreatedTime: pool.CreatedTime,
		})
	}

	return resp, nil
}

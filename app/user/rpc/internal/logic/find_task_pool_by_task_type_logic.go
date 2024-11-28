package logic

import (
	"T-driver/app/user/model"
	"context"
	"fmt"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindTaskPoolByTaskTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindTaskPoolByTaskTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindTaskPoolByTaskTypeLogic {
	return &FindTaskPoolByTaskTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindTaskPoolByTaskTypeLogic) FindTaskPoolByTaskType(in *pb.FindTaskPoolByTaskTypeReq) (*pb.TaskPoolListResp, error) {
	sb := squirrel.Select().Where("deleted_time =?", 0)
	if in.TaskName != "" {
		sb = sb.Where(squirrel.Like{"task_name": fmt.Sprintf("%s%%", in.TaskName)})
	}
	if in.TaskType > 0 {
		sb = sb.Where("task_type = ?", in.TaskType)
	}
	if in.Id > 0 {
		sb = sb.Where(squirrel.Eq{"id": in.Id})
	}
	if in.IsDisable > 0 {
		sb = sb.Where("is_disable = ?", in.IsDisable)
	}
	list := make([]*model.TaskPool, 0)
	var err error
	var n int64
	if in.Page == 0 && in.Size == 0 {
		if in.TaskName == "" && in.TaskType == 0 && in.IsDisable == 0 && in.Id == 0 {
			list, err = l.svcCtx.TaskPoolModel.FindAll(l.ctx, sb.OrderBy("sort desc"))
		} else {
			list, err = l.svcCtx.TaskPoolModel.List(l.ctx, sb.OrderBy("sort asc"))
		}
		n = int64(len(list))
	} else {
		list, n, err = l.svcCtx.TaskPoolModel.ListPage(l.ctx, in.Page, in.Size, sb.OrderBy("created_time desc"))
	}

	if err != nil {
		return nil, err
	}
	pools := make([]*pb.TaskPool, 0, len(list))
	for _, pool := range list {
		pools = append(pools, &pb.TaskPool{
			Id:          pool.Id,
			TaskType:    pool.TaskType,
			TaskName:    pool.TaskName,
			Integral:    pool.Integral,
			CreatedTime: pool.CreatedTime,
			Remark:      pool.Remark,
			IsDisable:   pool.IsDisable,
			Link:        pool.Link,
			Sort:        pool.Sort,
			TaskNameEn:  pool.TaskNameEn,
		})
	}
	return &pb.TaskPoolListResp{Count: n, TaskPools: pools}, nil
}

package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindTaskLogic {
	return &FindTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户任务列表
func (l *FindTaskLogic) FindTask(in *pb.FindTaskReq) (*pb.TaskResp, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Uid > 0 {
		sb = sb.Where("uid = ?", in.Uid)
	}
	if len(in.TaskPoolId) > 0 {
		sb = sb.Where(squirrel.Eq{"task_pool_id": in.TaskPoolId})
	}
	list, err := l.svcCtx.TaskModel.List(l.ctx, sb)
	if err != nil {
		return nil, err
	}
	resp := &pb.TaskResp{Total: int64(len(list)), Tasks: make([]*pb.Task, 0, len(list))}
	for _, task := range list {
		resp.Tasks = append(resp.Tasks, &pb.Task{
			Id:          task.Id,
			Uid:         task.Uid,
			FinishTime:  task.FinishTime,
			CreatedTime: task.CreatedTime,
			UpdatedTime: task.UpdatedTime,
			TaskPoolId:  task.TaskPoolId,
		})
	}
	return resp, nil
}

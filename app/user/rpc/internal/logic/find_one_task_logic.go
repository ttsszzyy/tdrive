package logic

import (
	"T-driver/app/user/model"
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneTaskLogic {
	return &FindOneTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOneTaskLogic) FindOneTask(in *pb.FindOneTaskReq) (*pb.Task, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Uid > 0 && in.TaskPoolId > 0 {
		sb = sb.Where("uid = ?", in.Uid).Where("task_pool_id = ?", in.TaskPoolId)
	}
	if in.Id > 0 {
		sb = sb.Where("id = ?", in.Id)
	}
	builder, err := l.svcCtx.TaskModel.FindOneByBuilder(l.ctx, sb)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.Task{}, nil
		}
		return nil, err
	}

	return &pb.Task{
		Id:          builder.Id,
		Uid:         builder.Uid,
		FinishTime:  builder.FinishTime,
		CreatedTime: builder.CreatedTime,
		UpdatedTime: builder.UpdatedTime,
	}, nil
}

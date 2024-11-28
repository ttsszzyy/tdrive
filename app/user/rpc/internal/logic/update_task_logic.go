package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
)

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改任务完成状态
func (l *UpdateTaskLogic) UpdateTask(in *pb.UpdateTaskReq) (*pb.Response, error) {
	one, err := l.svcCtx.TaskModel.FindOne(l.ctx, in.Id)
	if err != nil {
		logx.Errorf("查询任务失败: %v", err)
		return nil, errors.DbError()
	}
	if in.UpdatedTime > 0 {
		one.UpdatedTime = in.UpdatedTime
	}
	if in.FinishTime > 0 {
		one.FinishTime = in.FinishTime
	}
	taskPool, err := l.svcCtx.TaskPoolModel.FindOne(l.ctx, one.TaskPoolId)
	if err != nil {
		logx.Errorf("查询任务池失败: %v", err)
		return nil, errors.DbError()
	}
	u, err := l.svcCtx.UserPointsModel.FindOneByUidDeletedTime(l.ctx, one.Uid, 0)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return nil, errors.DbError()
	}
	err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		//更新任务状态为完成
		err = l.svcCtx.TaskModel.Update(l.ctx, one, session)
		if err != nil {
			logx.Errorf("修改任务失败:", err)
			return err
		}
		//更新用户积分
		u.Points += taskPool.Integral
		err = l.svcCtx.UserPointsModel.Update(l.ctx, u, session)
		if err != nil {
			return err
		}
		l.svcCtx.Redis.ZaddCtx(ctx, model.UserIntegral, u.Points, strconv.FormatInt(one.Uid, 10))
		return nil
	})
	if err != nil {
		logx.Error("修改任务状态失败: %v", err)
		return nil, err
	}
	return &pb.Response{}, nil
}

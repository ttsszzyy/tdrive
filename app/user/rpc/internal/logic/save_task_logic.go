package logic

import (
	"T-driver/app/user/model"
	"T-driver/common/errors"
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveTaskLogic {
	return &SaveTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改任务完成状态
func (l *SaveTaskLogic) SaveTask(in *pb.SaveTaskReq) (*pb.Response, error) {
	one, err := l.svcCtx.TaskPoolModel.FindOne(l.ctx, in.TaskPoolId)
	if err != nil {
		return nil, err
	}
	m := &model.Task{
		TaskPoolId:  in.TaskPoolId,
		Uid:         in.Uid,
		CreatedTime: time.Now().Unix(),
	}
	//签到需要添加时间
	if one.TaskType == 0 {
		m.UpdatedTime = time.Now().Unix()
	}
	//签到和邀请没有完成状态
	if one.TaskType != 0 && one.TaskType != 2 {
		m.FinishTime = time.Now().Unix()
	}
	u, err := l.svcCtx.UserPointsModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return nil, errors.DbError()
	}
	err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.TaskModel.Insert(l.ctx, m)
		if err != nil {
			return err
		}
		//更新用户积分
		u.Points += one.Integral
		err = l.svcCtx.UserPointsModel.Update(l.ctx, u, session)
		if err != nil {
			return err
		}
		l.svcCtx.Redis.ZaddCtx(ctx, model.UserIntegral, u.Points, strconv.FormatInt(u.Uid, 10))
		return nil
	})
	if err != nil {
		logx.Error("修改任务状态失败: %v", err)
		return nil, err
	}
	return &pb.Response{}, nil
}

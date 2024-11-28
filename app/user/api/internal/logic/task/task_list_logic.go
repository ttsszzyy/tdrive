package task

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/utils"
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/mr"
	"golang.org/x/sync/singleflight"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	sl     singleflight.Group
}

func NewTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskListLogic {
	return &TaskListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskListLogic) TaskList(req *types.Request) (resp *types.TaskListRes, err error) {

	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}

	// 查询任务池
	taskPools := &pb.TaskPoolListResp{}
	tasks := &pb.TaskResp{}
	one := &pb.User{}
	err = mr.Finish(func() error {
		var err error
		taskPools, err = l.queryTaskPools()
		return err
	}, func() error {
		var err error
		tasks, err = l.queryTasks(userData.User.ID)
		return err
	}, func() error {
		var err error
		one, err = l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
		return err
	})
	if err != nil {
		return nil, err
	}
	// 初始化响应结构
	resp = &types.TaskListRes{
		Url:   l.svcCtx.Config.FastReward.TgUrl + "?startapp=" + one.RecommendCode,
		Total: taskPools.Count,
		Tasks: make([]*types.Task, 0, taskPools.Count),
	}

	taskMap := make(map[int64]*pb.Task, len(tasks.Tasks))
	for _, task := range tasks.Tasks {
		taskMap[task.TaskPoolId] = task
	}
	for _, v := range taskPools.TaskPools {
		//校验任务是否完成
		task, ok := taskMap[v.Id]
		var id, updatedTime int64
		var isComplete int64 = 2 // 默认为未完成状态
		if ok {
			switch {
			//签到
			case v.TaskType == 0:
				if task.UpdatedTime > 0 && utils.IsSameDay(time.Now(), time.Unix(task.UpdatedTime, 0)) {
					isComplete = 1
				}
			//邀请 	保持默认值，即未完成状态
			case v.TaskType == 3:
			default:
				isComplete = 1
			}
			id = task.Id
			updatedTime = task.UpdatedTime
		}
		resp.Tasks = append(resp.Tasks, &types.Task{
			Id:          id,
			TaskPoolId:  v.Id,
			Uid:         userData.User.ID,
			TaskName:    v.TaskName,
			Integral:    v.Integral,
			IsComplete:  isComplete,
			TaskType:    v.TaskType,
			CreatedTime: v.CreatedTime,
			UpdatedTime: updatedTime,
			Url:         v.Link,
			Sort:        v.Sort,
			TaskNameEn:  v.TaskNameEn,
		})
	}
	return resp, nil
}

func (l *TaskListLogic) queryTaskPools() (*pb.TaskPoolListResp, error) {
	v, err, _ := l.sl.Do("FindTaskPoolByTaskType:2", func() (interface{}, error) {
		taskPools, err := l.svcCtx.Rpc.FindTaskPoolByTaskType(l.ctx, &pb.FindTaskPoolByTaskTypeReq{IsDisable: 2})
		if err != nil {
			return nil, err
		}
		return taskPools, nil
	})
	return v.(*pb.TaskPoolListResp), err
}

func (l *TaskListLogic) queryTasks(uid int64) (*pb.TaskResp, error) {
	tasks, err := l.svcCtx.Rpc.FindTask(l.ctx, &pb.FindTaskReq{
		Uid: uid,
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

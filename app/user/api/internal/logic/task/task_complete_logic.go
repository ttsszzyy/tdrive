package task

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	sg     singleflight.Group
}

func NewTaskCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskCompleteLogic {
	return &TaskCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskCompleteLogic) TaskComplete(req *types.TaskCompleteReq) (resp *types.TaskCompleteResp, err error) {
	resp = &types.TaskCompleteResp{
		Integral: req.Integral,
	}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	//添加分布式锁，防止重复
	err = l.svcCtx.Redis.Setex("user_task_complete_"+strconv.FormatInt(userData.User.ID, 10), "1", 2)
	if err != nil {
		return nil, errors.CustomError("Please try again later")
	}
	defer l.svcCtx.Redis.Del("user_task_complete_" + strconv.FormatInt(userData.User.ID, 10))

	_, err, _ = l.sg.Do(strconv.FormatInt(userData.User.ID, 10), func() (interface{}, error) {
		//校验是否完成
		task, err := l.svcCtx.Rpc.FindOneTask(l.ctx, &pb.FindOneTaskReq{TaskPoolId: req.TaskPoolId, Uid: userData.User.ID})
		if err != nil {
			return nil, errors.FromRpcError(err)
		}
		if req.Id > 0 {
			list, err := l.svcCtx.Rpc.FindTaskPoolByTaskType(l.ctx, &pb.FindTaskPoolByTaskTypeReq{Id: req.TaskPoolId})
			if err != nil {
				return nil, err
			}
			//修改签到时间
			if list.TaskPools[0].TaskType == 0 {
				// 获取当前时间
				now := time.Now()
				midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				if task.UpdatedTime >= midnight.Unix() {
					return nil, errors.CustomError("今日已簽到")
				}
				_, err := l.svcCtx.Rpc.UpdateTask(l.ctx, &pb.UpdateTaskReq{
					Id:          req.Id,
					UpdatedTime: time.Now().Unix(),
				})
				if err != nil {
					return nil, err
				}
			}

		} else {
			if task.Id > 0 {
				return nil, errors.CustomError("任務已完成")
			}
			p := &pb.SaveTaskReq{
				Uid:        userData.User.ID,
				TaskPoolId: req.TaskPoolId,
			}
			_, err = l.svcCtx.Rpc.SaveTask(l.ctx, p)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

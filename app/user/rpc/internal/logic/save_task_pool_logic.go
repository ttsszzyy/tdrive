package logic

import (
	"T-driver/app/user/model"
	"T-driver/common/errors"
	"context"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveTaskPoolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveTaskPoolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveTaskPoolLogic {
	return &SaveTaskPoolLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveTaskPoolLogic) SaveTaskPool(in *pb.SaveTaskPoolReq) (*pb.Response, error) {
	if in.Id > 0 {
		one, err := l.svcCtx.TaskPoolModel.FindOne(l.ctx, in.Id)
		if err != nil {
			return nil, errors.DbError()
		}
		err = l.svcCtx.TaskPoolModel.Update(l.ctx, &model.TaskPool{
			Id:          in.Id,
			TaskType:    in.TaskType,
			TaskName:    in.TaskName,
			Integral:    in.Integral,
			CreatedTime: one.CreatedTime,
			UpdatedTime: time.Now().Unix(),
			Remark:      in.Remark,
			IsDisable:   in.IsDisable,
			Link:        in.Link,
			Sort:        in.Sort,
			TaskNameEn:  in.TaskNameEn,
		})
		if err != nil {
			logx.Error(err)
			return nil, errors.DbError()
		}
	} else {
		_, err := l.svcCtx.TaskPoolModel.Insert(l.ctx, &model.TaskPool{
			TaskType:    in.TaskType,
			TaskName:    in.TaskName,
			Integral:    in.Integral,
			CreatedTime: time.Now().Unix(),
			Remark:      in.Remark,
			IsDisable:   in.IsDisable,
			Link:        in.Link,
			Sort:        in.Sort,
			TaskNameEn:  in.TaskNameEn,
		})
		if err != nil {
			logx.Error(err)
			return nil, errors.DbError()
		}
	}
	//清理查询缓存
	l.svcCtx.TaskPoolModel.ClearCache()

	return &pb.Response{}, nil
}

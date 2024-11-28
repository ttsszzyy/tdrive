package logic

import (
	"T-driver/common/errors"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelTaskPoolByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelTaskPoolByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelTaskPoolByIdLogic {
	return &DelTaskPoolByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelTaskPoolByIdLogic) DelTaskPoolById(in *pb.DelTaskPoolByIdReq) (*pb.Response, error) {
	if in.Id > 0 {
		//删除任务
		list, err := l.svcCtx.TaskModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"task_pool_id": in.Id}))
		if err != nil {
			return nil, errors.DbError()
		}
		//删除任务
		ids := make([]int64, 0, len(list))
		for _, v := range list {
			ids = append(ids, v.Id)
		}
		err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			//删除任务池
			err := l.svcCtx.TaskPoolModel.Delete(ctx, in.Id, session)
			if err != nil {
				return err
			}
			//清理查询缓存
			l.svcCtx.TaskPoolModel.ClearCache()
			//删除任务
			err = l.svcCtx.TaskModel.Deletes(ctx, ids, session)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			logx.Error(err)
			return nil, err
		}

	}

	return &pb.Response{}, nil
}

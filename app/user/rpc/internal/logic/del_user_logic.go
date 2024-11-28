package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserLogic {
	return &DelUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除用户
func (l *DelUserLogic) DelUser(in *pb.DelUserReq) (*pb.Response, error) {
	//删除用户任务
	list3, err := l.svcCtx.TaskModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"uid": in.Uid}).Where("deleted_time = ?", 0))
	if err != nil {
		return nil, err
	}
	taskIds := make([]int64, 0, len(list3))
	for _, task := range list3 {
		taskIds = append(taskIds, task.Id)
	}
	//删除用户分享
	list4, err := l.svcCtx.ShareModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"uid": in.Uid}).Where("deleted_time = ?", 0))
	if err != nil {
		return nil, err
	}
	shareIds := make([]int64, 0, len(list4))
	for _, task := range list4 {
		shareIds = append(shareIds, task.Id)
	}
	//删除用户邀请
	list5, err := l.svcCtx.UserInviteModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"pid": in.Uid}).Where("deleted_time = ?", 0))
	if err != nil {
		return nil, err
	}
	inviteIds := make([]int64, 0, len(list5))
	for _, v := range list5 {
		inviteIds = append(inviteIds, v.Id)
	}
	//删除用户邀请奖励
	list6, err := l.svcCtx.UserInviteRewardModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"uid": in.Uid}).Where("deleted_time = ?", 0))
	if err != nil {
		return nil, err
	}
	inviteRewardIds := make([]int64, 0, len(list6))
	for _, v := range list6 {
		inviteRewardIds = append(inviteRewardIds, v.Id)
	}
	UserPoints, err := l.svcCtx.UserPointsModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		return nil, err
	}
	UserStorage, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		return nil, err
	}
	//删除用户兑换记录
	list7, err := l.svcCtx.UserStorageExchangeModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"uid": in.Uid}).Where("deleted_time = ?", 0))
	if err != nil {
		return nil, err
	}
	exchangeIds := make([]int64, 0, len(list7))
	for _, v := range list7 {
		exchangeIds = append(exchangeIds, v.Id)
	}
	//删除用户
	user, err := l.svcCtx.UserModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err = l.svcCtx.TaskModel.Deletes(ctx, taskIds, session)
		if err != nil {
			return err
		}
		err = l.svcCtx.ShareModel.Deletes(ctx, shareIds, session)
		if err != nil {
			return err
		}
		err = l.svcCtx.UserInviteModel.Deletes(ctx, inviteIds, session)
		if err != nil {
			return err
		}
		err = l.svcCtx.UserInviteRewardModel.Deletes(ctx, inviteRewardIds, in.Uid, session)
		if err != nil {
			return err
		}
		err = l.svcCtx.UserPointsModel.Delete(ctx, UserPoints.Id, session)
		if err != nil {
			return err
		}
		err = l.svcCtx.UserStorageModel.Delete(ctx, UserStorage.Id, session)
		if err != nil {
			return err
		}
		err = l.svcCtx.UserStorageExchangeModel.Deletes(ctx, exchangeIds, session)
		if err != nil {
			return err
		}
		err = l.svcCtx.UserModel.Delete(ctx, user.Id, session)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	//清理redis
	l.svcCtx.Redis.ZremCtx(l.ctx, model.UserIntegral, strconv.FormatInt(in.Uid, 10))
	return &pb.Response{}, nil
}

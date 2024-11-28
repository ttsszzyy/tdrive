package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimInviteRewardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClaimInviteRewardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimInviteRewardLogic {
	return &ClaimInviteRewardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 领取邀请奖励
func (l *ClaimInviteRewardLogic) ClaimInviteReward(in *pb.ClaimInviteRewardReq) (*pb.Response, error) {
	one := &model.UserPoints{}
	var invite int64
	userInviteReward := &model.UserInviteReward{}
	var err error
	err = mr.Finish(func() error {
		one, err = l.svcCtx.UserPointsModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
		if err != nil {
			return err
		}
		return err
	}, func() error {
		invite, err = l.svcCtx.UserModel.CountInvite(l.ctx, squirrel.Select().Where(squirrel.Eq{"pid": in.Uid}))
		if err != nil {
			return err
		}
		return err
	}, func() error {
		userInviteReward, err = l.svcCtx.UserInviteRewardModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
		if err != nil {
			if err == model.ErrNotFound {
				userInviteReward = &model.UserInviteReward{}
			} else {
				return errors.DbError()
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	//查询可以领取的范围
	where := squirrel.Select().Where("param_type = ?  and backup_value > ? and backup_value  <= ?", 5, userInviteReward.InviteNum, invite)
	list, err := l.svcCtx.DictModel.List(l.ctx, where)
	if err != nil {
		return nil, err
	}
	var points int64
	for _, dict := range list {
		num, _ := strconv.ParseInt(dict.Value, 10, 64)
		points += num
		userInviteReward.InviteNum = dict.BackupValue
	}
	if points > 0 {
		//更新用户邀请奖励数量
		err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			if userInviteReward.Id > 0 {
				userInviteReward.UpdatedTime = time.Now().Unix()
				err = l.svcCtx.UserInviteRewardModel.Update(l.ctx, userInviteReward, session)
				if err != nil {
					return err
				}
			} else {
				userInviteReward.Uid = in.Uid
				userInviteReward.CreatedTime = time.Now().Unix()
				_, err = l.svcCtx.UserInviteRewardModel.Insert(l.ctx, userInviteReward, session)
				if err != nil {
					return err
				}
			}

			//更新用户积分
			one.Points += points
			one.ReqPoints += points
			err = l.svcCtx.UserPointsModel.Update(l.ctx, one, session)
			if err != nil {
				return err
			}
			//添加用户到积分排行榜
			l.svcCtx.Redis.ZaddCtx(l.ctx, model.UserIntegral, one.Points, strconv.FormatInt(one.Uid, 10))
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		return nil, status.Error(errors.ErrCodeNotTask, "未達標")
	}

	return &pb.Response{}, nil
}

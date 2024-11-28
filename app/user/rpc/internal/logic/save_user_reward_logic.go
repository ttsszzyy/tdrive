package logic

import (
	"T-driver/app/user/model"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveUserRewardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveUserRewardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveUserRewardLogic {
	return &SaveUserRewardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存邀请奖励
func (l *SaveUserRewardLogic) SaveUserReward(in *pb.SaveUserRewardReq) (*pb.Response, error) {
	ids := []int64{in.Uid, in.Pid}
	if in.Points > 0 {
		list, err := l.svcCtx.UserPointsModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"uid": ids}).Where("deleted_time = ?", 0))
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			for _, v := range list {
				if in.Uid == v.Uid {
					v.Points += in.Points
					v.ReqPoints += in.Points
				} else if in.Pid == v.Uid {
					v.Points += in.Points + in.PointsVip
					v.ReqPoints += in.Points + in.PointsVip
				}

				err = l.svcCtx.UserPointsModel.Update(l.ctx, v, session)
				if err != nil {
					return err
				}
				//添加用户到积分排行榜
				l.svcCtx.Redis.ZaddCtx(l.ctx, model.UserIntegral, v.Points, strconv.FormatInt(v.Uid, 10))
			}
			_, err = l.svcCtx.UserInviteModel.Insert(l.ctx, &model.UserInvite{
				Pid:          in.Pid,
				Uid:          in.Uid,
				InvitePoints: in.Points + in.PointsVip,
				CreatedTime:  time.Now().Unix(),
			}, session)
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	return &pb.Response{}, nil
}

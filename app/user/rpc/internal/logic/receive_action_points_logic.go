package logic

import (
	"context"
	"strconv"
	"time"

	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReceiveActionPointsLogic 用户领取活动积分
type ReceiveActionPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewReceiveActionPointsLogic 新建 用户领取活动积分
func NewReceiveActionPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveActionPointsLogic {
	return &ReceiveActionPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReceiveActionPoints 实现用户领取活动积分
func (l *ReceiveActionPointsLogic) ReceiveActionPoints(in *pb.ReceiveActionPointsReq) (*pb.Response, error) {
	var lastUsedTime int64

	// 先判断用户是否符合领取过
	_, err := l.svcCtx.ActionRecordModel.FindOneByUidAction(l.ctx, in.Uid, in.Name)
	switch err {
	case nil:
		return nil, status.Error(codes.Internal, errors.GetError(errors.ErrReceivedActionPoints, in.Lan).Msg())
	case model.ErrNotFound:
	default:
		logx.Error(err)
		return nil, status.Error(codes.Internal, errors.DbError(in.Lan).Msg())
	}
	// 判断此次活动中，该ip是否已经被使用过
	num, err := l.svcCtx.ActionRecordModel.GetIPNumByAction(l.ctx, in.Name, in.Ip)
	if err != nil {
		logx.Error(err)
		return nil, status.Error(codes.Internal, errors.DbError(in.Lan).Msg())
	}
	if num >= 1 {
		return nil, status.Error(codes.Internal, errors.GetError(errors.ErrSameIPForAction, in.Lan).Msg())
	}
	// 获取该用户的注册时间
	user, err := l.svcCtx.UserModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		logx.Error(err)
		return nil, status.Error(codes.Internal, errors.DbError(in.Lan).Msg())
	}
	// 获取该用户上一次的使用时间
	asset, err := l.svcCtx.AssetsModel.GetLastAssetInfo(l.ctx, in.Uid)
	switch err {
	case nil:
		lastUsedTime = asset.CreatedTime
	case model.ErrNotFound:
		lastUsedTime = 0
	default:
		logx.Error(err)
		return nil, status.Error(codes.Internal, errors.DbError(in.Lan).Msg())
	}

	// 修改用户积分信息
	err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, s sqlx.Session) error {
		// 添加用户领取记录
		_, err := l.svcCtx.ActionRecordModel.Insert(ctx, &model.ActionRecord{
			Uid:          in.Uid,
			Location:     in.Location,
			Ip:           in.Ip,
			Points:       in.Points,
			RegisterTime: user.CreatedTime,
			LastUsedTime: lastUsedTime,
			Action:       in.Name,
			CreatedTime:  time.Now().Unix(),
		}, s)
		if err != nil {
			logx.Error(err)
			return status.Error(codes.Internal, errors.DbError(in.Lan).Msg())
		}
		// 获取用户积分信息，不存在则添加
		one, err := l.svcCtx.UserPointsModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
		switch err {
		case nil:
			one.Points += in.Points
			one.UpdatedTime = time.Now().Unix()
			err = l.svcCtx.UserPointsModel.Update(l.ctx, one, s)
			if err != nil {
				logx.Error(err)
				return status.Error(codes.Internal, errors.DbError(in.Lan).Msg())
			}
		case model.ErrNotFound:
			one = &model.UserPoints{Uid: in.Uid, Points: in.Points, CreatedTime: time.Now().Unix(), UpdatedTime: time.Now().Unix()}
			_, err = l.svcCtx.UserPointsModel.Insert(l.ctx, one, s)
			if err != nil {
				logx.Error(err)
				return status.Error(codes.Internal, errors.DbError(in.Lan).Msg())
			}
		default:
			logx.Error(err)
			return status.Error(codes.Internal, errors.DbError(in.Lan).Msg())
		}

		//添加用户到积分排行榜
		l.svcCtx.Redis.ZaddCtx(l.ctx, model.UserIntegral, one.Points, strconv.FormatInt(one.Uid, 10))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.Response{}, nil
}

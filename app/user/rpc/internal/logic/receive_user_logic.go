package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveUserLogic {
	return &ReceiveUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 领取用户奖励空间
func (l *ReceiveUserLogic) ReceiveUser(in *pb.ReceiveUserReq) (*pb.Response, error) {
	u, err := l.svcCtx.UserModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		u.IsReceive = in.IsReceive
		err := l.svcCtx.UserModel.Update(ctx, u, session)
		if err != nil {
			return err
		}
		switch {
		case in.Storage > 0:
			p, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
			if err != nil {
				return err
			}
			p.Storage += in.Storage
			err = l.svcCtx.UserStorageModel.Update(ctx, p, session)
			if err != nil {
				return err
			}
		default:
			p, err := l.svcCtx.UserPointsModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
			if err != nil {
				return err
			}
			p.Points += in.Points
			err = l.svcCtx.UserPointsModel.Update(ctx, p, session)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.Response{}, nil
}

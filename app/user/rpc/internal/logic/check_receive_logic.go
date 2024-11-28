package logic

import (
	"context"

	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CheckReceiveLogic 校验是否满足领取条件
type CheckReceiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewCheckReceiveLogic 新建 校验是否满足领取条件
func NewCheckReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckReceiveLogic {
	return &CheckReceiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CheckReceive 实现 校验是否满足领取条件
func (l *CheckReceiveLogic) CheckReceive(in *pb.ReceiveActionPointsReq) (*pb.Response, error) {
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

	return &pb.Response{}, nil
}

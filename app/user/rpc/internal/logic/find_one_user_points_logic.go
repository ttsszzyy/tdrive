package logic

import (
	"T-driver/app/user/model"
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneUserPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneUserPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneUserPointsLogic {
	return &FindOneUserPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户积分
func (l *FindOneUserPointsLogic) FindOneUserPoints(in *pb.FindOneUserPointsReq) (*pb.FindOneUserPointsResp, error) {
	one, err := l.svcCtx.UserPointsModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.FindOneUserPointsResp{}, nil
		}
		return nil, err
	}

	return &pb.FindOneUserPointsResp{
		Id:          one.Id,
		CreatedTime: one.CreatedTime,
		Points:      one.Points,
		ReqPoints:   one.ReqPoints,
		Uid:         one.Uid,
	}, nil
}

package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserPointsLogic {
	return &FindUserPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserPointsLogic) FindUserPoints(in *pb.FindUserPointsReq) (*pb.FindUserPointsResp, error) {
	sb := squirrel.Select().Where("deleted_time =?", 0)
	if len(in.Uid) > 0 {
		sb = sb.Where(squirrel.Eq{"uid": in.Uid})
	}
	list, err := l.svcCtx.UserPointsModel.List(l.ctx, sb)
	if err != nil {
		return nil, err
	}
	pbList := make([]*pb.FindOneUserPointsResp, 0)
	for _, v := range list {
		pbList = append(pbList, &pb.FindOneUserPointsResp{
			Uid:         v.Uid,
			Points:      v.Points,
			CreatedTime: v.CreatedTime,
			ReqPoints:   v.ReqPoints,
		})
	}
	return &pb.FindUserPointsResp{List: pbList}, nil
}

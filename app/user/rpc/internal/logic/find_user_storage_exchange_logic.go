package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserStorageExchangeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserStorageExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserStorageExchangeLogic {
	return &FindUserStorageExchangeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户空间兑换列表
func (l *FindUserStorageExchangeLogic) FindUserStorageExchange(in *pb.FindUserStorageExchangeReq) (*pb.FindUserStorageExchangeRespList, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Uid > 0 {
		sb = sb.Where("uid = ?", in.Uid)
	}
	list, total, err := l.svcCtx.UserStorageExchangeModel.ListPage(l.ctx, in.Page, in.Size, sb.OrderBy("created_time desc"))
	if err != nil {
		return nil, err
	}
	pbList := make([]*pb.FindUserStorageExchangeResp, 0)
	for _, v := range list {
		pbList = append(pbList, &pb.FindUserStorageExchangeResp{
			CreatedTime:     v.CreatedTime,
			Id:              v.Id,
			StorageExchange: v.ExchangeStorage,
			Uid:             v.Uid,
		})
	}
	return &pb.FindUserStorageExchangeRespList{
		List:  pbList,
		Total: total,
	}, nil
}

package logic

import (
	"T-driver/common/lib/json"
	"context"

	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindShareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindShareLogic {
	return &FindShareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看分析任务
func (l *FindShareLogic) FindShare(in *pb.FindShareReq) (*pb.FindShareResp, error) {
	list, total, err := l.svcCtx.ShareModel.ListPage(l.ctx, in.Page, in.Size, squirrel.Select().Where("uid = ?", in.Uid).Where("deleted_time = ?", 0).OrderBy("created_time desc"))
	if err != nil {
		return nil, err
	}
	resp := &pb.FindShareResp{Total: total, Shares: make([]*pb.Share, 0, len(list))}
	for _, share := range list {
		assetIds := make([]int64, 0)
		json.Unmarshal([]byte(share.AssetIds), &assetIds)
		resp.Shares = append(resp.Shares, &pb.Share{
			Id:            share.Id,
			Uid:           share.Uid,
			Uuid:          share.Uuid,
			AssetIds:      assetIds,
			AssetName:     share.AssetName,
			AssetSize:     share.AssetSize,
			AssetType:     share.AssetType,
			Link:          share.Link,
			EffectiveTime: share.EffectiveTime,
			ReadNum:       share.ReadNum,
			SaveNum:       share.SaveNum,
			CreatedTime:   share.CreatedTime,
			Password:      share.Password,
		})
	}
	return resp, nil
}

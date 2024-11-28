package logic

import (
	"context"

	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountAssetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountAssetsLogic {
	return &CountAssetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看资源数据量
func (l *CountAssetsLogic) CountAssets(in *pb.CountAssetsReq) (*pb.CountAssetsResp, error) {
	sb := squirrel.Select()
	if in.AssetName != "" {
		sb = sb.Where("asset_name = ?", in.AssetName)
	}
	if in.Pid > 0 {
		sb = sb.Where("pid = ?", in.Pid)
	}
	if in.Uid > 0 {
		sb = sb.Where("uid = ?", in.Uid)
	}
	if in.Cid != "" {
		sb = sb.Where("cid = ?", in.Cid)
	}
	if in.IsTag > 0 {
		sb = sb.Where("is_tag = ?", in.IsTag)
	}
	if in.IsDel {
		sb = sb.Where("deleted_time > ?", 0)
	} else {
		sb = sb.Where("deleted_time = ?", 0)
	}
	if in.AssetTypes > 0 {
		sb = sb.Where(squirrel.Eq{"asset_type": in.AssetTypes})
	}
	total, err := l.svcCtx.AssetsModel.Count(l.ctx, sb)
	if err != nil {
		return nil, err
	}
	return &pb.CountAssetsResp{Total: total}, nil
}

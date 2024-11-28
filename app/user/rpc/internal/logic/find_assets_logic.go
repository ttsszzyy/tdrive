package logic

import (
	"T-driver/app/user/model"
	"context"

	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAssetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAssetsLogic {
	return &FindAssetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看资源
func (l *FindAssetsLogic) FindAssets(in *pb.FindAssetsReq) (*pb.FindAssetsResp, error) {
	sb := squirrel.Select()
	if in.AssetName != "" {
		sb = sb.Where(squirrel.Like{"asset_name": in.AssetName + "%"})
	}
	if in.Pid > 0 {
		sb = sb.Where("pid = ?", in.Pid)
	}
	if len(in.Ids) > 0 {
		sb = sb.Where(squirrel.Eq{"id": in.Ids})
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
	if len(in.Status) > 0 {
		sb = sb.Where(squirrel.Eq{"status": in.Status})
	}
	if in.IsReport > 0 {
		sb = sb.Where("is_report = ?", in.IsReport)
	}
	if in.StartTime > 0 {
		sb = sb.Where(squirrel.GtOrEq{"created_time": in.StartTime})
	}
	if in.EndTime > 0 {
		sb = sb.Where(squirrel.LtOrEq{"created_time": in.EndTime})
	}
	if len(in.AssetTypes) > 0 {
		sb = sb.Where(squirrel.Eq{"asset_type": in.AssetTypes})
	}
	if in.IsAdd {
		sb = sb.Where("created_time != ?", 0)
	}
	sort := "desc"
	if in.Sort == 1 {
		sort = "asc"
	}
	switch in.Order {
	case 1:
		sb = sb.OrderBy("updated_time " + sort)
	case 2:
		sb = sb.OrderBy("asset_size " + sort)
	case 3:
		sb = sb.OrderBy("asset_name " + sort)
	default:
		sb = sb.OrderBy("created_time " + sort)
	}
	sb = sb.OrderBy("id " + sort)
	list := make([]*model.Assets, 0)
	var total int64
	var err error
	if in.Page == 0 && in.Size == 0 {
		list, err = l.svcCtx.AssetsModel.List(l.ctx, sb)
		total = int64(len(list))
	} else {
		list, total, err = l.svcCtx.AssetsModel.ListPage(l.ctx, in.Page, in.Size, sb)
	}
	if err != nil {
		return nil, err
	}
	resp := &pb.FindAssetsResp{
		Total:  total,
		Assets: make([]*pb.Assets, 0, len(list)),
	}
	for _, asset := range list {

		resp.Assets = append(resp.Assets, &pb.Assets{
			Id:          asset.Id,
			Uid:         asset.Uid,
			Cid:         asset.Cid,
			AssetName:   asset.AssetName,
			AssetSize:   asset.AssetSize,
			AssetType:   asset.AssetType,
			TransitType: asset.TransitType,
			IsTag:       asset.IsTag,
			Pid:         asset.Pid,
			Source:      asset.Source,
			CreatedTime: asset.CreatedTime,
			UpdatedTime: asset.UpdatedTime,
			DeletedTime: asset.DeletedTime,
			Status:      asset.Status,
			IsReport:    asset.IsReport,
			ReportType:  asset.ReportType,
			Link:        asset.Link,
			IsDefault:   asset.IsDefault,
		})
	}

	return resp, nil
}

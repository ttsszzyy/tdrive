package logic

import (
	"T-driver/app/user/model"
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAssetsNoDelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAssetsNoDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAssetsNoDelLogic {
	return &FindAssetsNoDelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindAssetsNoDelLogic) FindAssetsNoDel(in *pb.FindAssetsReq) (*pb.FindAssetsResp, error) {
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

package logic

import (
	"T-driver/app/user/model"
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneAssetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneAssetsLogic {
	return &FindOneAssetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看资源
func (l *FindOneAssetsLogic) FindOneAssets(in *pb.FindOneAssetsReq) (*pb.Assets, error) {
	one := &model.Assets{}
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Cid != "" {
		sb = sb.Where("cid = ?", in.Cid)
	}
	if in.Uid > 0 {
		sb = sb.Where("uid = ?", in.Uid)
	}
	if in.Uid > 0 {
		sb = sb.Where("uid = ?", in.Uid)
	}
	if in.IsDefault > 0 {
		sb = sb.Where("is_default = ?", in.IsDefault)
	}
	if in.Id > 0 {
		sb = sb.Where("id = ?", in.Id)
	}
	if in.AssetName != "" {
		sb = sb.Where("asset_name = ?", in.AssetName)
	}
	if in.AssetType > 0 {
		sb = sb.Where("asset_type = ?", in.AssetType)
	}
	one, err := l.svcCtx.AssetsModel.FindOneByBuilder(l.ctx, sb)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.Assets{}, nil
		}
		return nil, err
	}
	return &pb.Assets{
		Id:          one.Id,
		Uid:         one.Uid,
		Cid:         one.Cid,
		AssetName:   one.AssetName,
		AssetSize:   one.AssetSize,
		AssetType:   one.AssetType,
		TransitType: one.TransitType,
		IsTag:       one.IsTag,
		Pid:         one.Pid,
		Source:      one.Source,
		CreatedTime: one.CreatedTime,
		UpdatedTime: one.UpdatedTime,
		Link:        one.Link,
	}, nil
}

package logic

import (
	"T-driver/app/user/model"
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAssetsNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAssetsNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAssetsNameLogic {
	return &UpdateAssetsNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改资源name,tag,status
func (l *UpdateAssetsNameLogic) UpdateAssetsName(in *pb.UpdateAssetsNameReq) (*pb.Response, error) {
	one, err := l.svcCtx.AssetsModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.Response{}, nil
		}
		return nil, err
	}
	if in.AssetName != "" {
		one.AssetName = in.AssetName
	}
	if in.IsTag > 0 {
		one.IsTag = in.IsTag
	}
	if in.Status > 0 {
		one.Status = in.Status
	}
	if in.AssetSize > 0 {
		one.AssetSize = in.AssetSize
	}
	if in.IsReport > 0 {
		one.IsReport = in.IsReport
	}
	if in.ReportType > 0 {
		one.ReportType = in.ReportType
	}
	if in.Cid != "" {
		one.Cid = in.Cid
	}
	if in.Link != "" {
		one.Link = in.Link
	}
	one.UpdatedTime = time.Now().Unix()
	err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err = l.svcCtx.AssetsModel.Update(l.ctx, one, session)
		if err != nil {
			return err
		}
		if in.AssetSize > 0 {
			u, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, one.Uid, 0)
			if err != nil {
				return err
			}
			u.StorageUse += in.AssetSize
			if u.StorageUse < 0 {
				u.StorageUse = 0
			}
			err = l.svcCtx.UserStorageModel.Update(l.ctx, u, session)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logx.Error("修改资源失败")
		return nil, err
	}

	return &pb.Response{}, nil
}

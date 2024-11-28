package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/utils"
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAssetFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAssetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAssetFileLogic {
	return &UpdateAssetFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改传输中资源tag,status,cid,link,assetId
func (l *UpdateAssetFileLogic) UpdateAssetFile(in *pb.UpdateAssetFileReq) (*pb.UpdateAssetResp, error) {
	resp := new(pb.UpdateAssetResp)

	one, err := l.svcCtx.AssetFileModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Tag > 0 {
		one.IsTag = in.Tag
	}
	if in.Cid != "" {
		one.Cid = in.Cid
	}
	if in.Link != "" {
		one.Link = in.Link
	}
	if in.Status > 0 {
		one.Status = in.Status
	}
	if in.AssetSize > 0 {
		one.AssetSize = in.AssetSize
	}

	if in.Status == 3 {
		asset := &model.Assets{
			Uid:         one.Uid,
			Cid:         one.Cid,
			AssetName:   one.AssetName,
			AssetSize:   one.AssetSize,
			AssetType:   one.AssetType,
			TransitType: 1,
			IsTag:       one.IsTag,
			Pid:         one.Pid,
			Source:      one.Source,
			IsReport:    2,
			Status:      in.Status,
			Link:        one.Link,
			CreatedTime: one.CreateAt.Unix(),
			UpdatedTime: one.UpdateAt.Unix(),
		}
		//todo 待优化 校验文件名是否重复
		if asset.AssetName != "" && asset.AssetType != 1 {
			//计数器
			num := 0
		NEXT:
			count, err := l.svcCtx.AssetsModel.Count(l.ctx, squirrel.Select().Where("asset_name = ? and deleted_time = ? and pid = ?", asset.AssetName, 0, asset.Pid))
			if err != nil {
				return nil, err
			}
			if count >= 1 {
				num++
				if asset.AssetType != 1 {
					filename, ext := utils.GetFileName(asset.AssetName)
					if strings.HasSuffix(filename, fmt.Sprintf("(%d)", num-1)) {
						replace := utils.ReplaceLastOccurrence(filename, fmt.Sprintf("(%d)", num-1), fmt.Sprintf("(%d)", num))
						asset.AssetName = fmt.Sprintf("%s%s", replace, ext)
						goto NEXT
					} else {
						asset.AssetName = fmt.Sprintf("%s(%d)%s", filename, num, ext)
						goto NEXT
					}

				} else {
					//文件夹
					if strings.HasSuffix(asset.AssetName, fmt.Sprintf("(%d)", num-1)) {
						replace := utils.ReplaceLastOccurrence(asset.AssetName, fmt.Sprintf("(%d)", num-1), fmt.Sprintf("(%d)", num))
						asset.AssetName = replace
						goto NEXT
					} else {
						asset.AssetName = fmt.Sprintf("%s(%d)", asset.AssetName, num)
						goto NEXT
					}
				}
			}
		}

		//添加数据库
		err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			data, err := l.svcCtx.AssetsModel.Insert(l.ctx, asset, session)
			if err != nil {
				return err
			}
			if asset.AssetSize > 0 {
				u, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, asset.Uid, 0)
				if err != nil {
					return err
				}
				u.StorageUse += asset.AssetSize
				if u.StorageUse < 0 {
					u.StorageUse = 0
				}
				err = l.svcCtx.UserStorageModel.Update(l.ctx, u, session)
				if err != nil {
					return err
				}
			}
			id, err := data.LastInsertId()
			if err != nil {
				return err
			}
			resp.Id = id
			one.AssetId = id
			_, err = l.svcCtx.AssetFileModel.Update(l.ctx, one)
			if err != nil {
				return err
			}
			return err
		})
		if err != nil {
			logx.Error("修改资源失败")
			return nil, err
		}
	} else {
		err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			if in.Tag > 0 && in.AssetId > 0 {
				assets, err := l.svcCtx.AssetsModel.FindOne(ctx, in.AssetId)
				if err != nil {
					return err
				}
				assets.IsTag = in.Tag
				err = l.svcCtx.AssetsModel.Update(ctx, assets, session)
				if err != nil {
					return err
				}
			}
			_, err = l.svcCtx.AssetFileModel.Update(ctx, one)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			logx.Error("修改资源失败")
			return nil, err
		}
	}

	return resp, nil
}

package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"T-driver/common/utils"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAssetsCopyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAssetsCopyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAssetsCopyLogic {
	return &UpdateAssetsCopyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 资源复制
func (l *UpdateAssetsCopyLogic) UpdateAssetsCopy(in *pb.UpdateAssetsCopyReq) (*pb.Response, error) {
	err := l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var use int64
		var uid int64
		for _, id := range in.Ids {
			one, err := l.svcCtx.AssetsModel.FindOne(l.ctx, id, session)
			if err != nil {
				return err
			}
			one.Id = 0
			one.Pid = in.Pid
			one.CreatedTime = 0
			one.UpdatedTime = time.Now().Unix()
			one.IsDefault = 0
			//校验文件名是否重复
			if one.AssetName != "" {
				//计数器
				num := 0
			NEXT:
				count, err := l.svcCtx.AssetsModel.Count(l.ctx, squirrel.Select().Where("asset_name = ? and deleted_time = ? and pid = ?", one.AssetName, 0, in.Pid))
				if err != nil {
					return err
				}
				if count >= 1 {
					num++
					if one.AssetType != 1 {
						//文件
						filename, ext := utils.GetFileName(one.AssetName)
						if strings.HasSuffix(filename, fmt.Sprintf("(%d)", num-1)) {
							replace := utils.ReplaceLastOccurrence(filename, fmt.Sprintf("(%d)", num-1), fmt.Sprintf("(%d)", num))
							one.AssetName = fmt.Sprintf("%s%s", replace, ext)
							goto NEXT
						} else {
							one.AssetName = fmt.Sprintf("%s(%d)%s", filename, num, ext)
							goto NEXT
						}
					} else {
						//文件夹
						if strings.HasSuffix(one.AssetName, fmt.Sprintf("(%d)", num-1)) {
							replace := utils.ReplaceLastOccurrence(one.AssetName, fmt.Sprintf("(%d)", num-1), fmt.Sprintf("(%d)", num))
							one.AssetName = replace
							goto NEXT
						} else {
							one.AssetName = fmt.Sprintf("%s(%d)", one.AssetName, num)
							goto NEXT
						}
					}

				}
			}
			newId, err := l.svcCtx.AssetsModel.Insert(l.ctx, one, session)
			if err != nil {
				return err
			}
			if one.AssetType == 1 {
				lastInsertId, err := newId.LastInsertId()
				if err != nil {
					logx.Errorf("获取新id失败: %v", err)
					return err
				}
				//复制文件夹下面的文件到新的目录下
				subdirectories, err := l.Subdirectories(lastInsertId, id, session)
				if err != nil {
					logx.Errorf("复制文件夹下面的文件到新的目录下失败: %v", err)
					return err
				}
				use += subdirectories
			}
			use += one.AssetSize
			if uid == 0 {
				uid = one.Uid
			}
		}
		//修改用户使用空间
		if use > 0 {

			u, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, uid, 0)
			if err != nil {
				return err
			}
			u.StorageUse += use
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
		logx.Errorf("资源复制失败: %v", err)
		return nil, errors.DbError()
	}
	return &pb.Response{}, nil
}

func (l *UpdateAssetsCopyLogic) Subdirectories(pid, old int64, s sqlx.Session) (use int64, err error) {
	list, err := l.svcCtx.AssetsModel.List(l.ctx, squirrel.Select().Where("pid = ? and deleted_time = ?", old, 0))
	if err != nil {
		logx.Errorf("获取子目录失败: %v", err)
		return 0, err
	}
	for _, v := range list {
		one := &model.Assets{}
		copier.Copy(&one, v)
		one.Pid = pid
		one.IsDefault = 0
		newId, err := l.svcCtx.AssetsModel.Insert(l.ctx, one, s)
		if err != nil {
			logx.Errorf("插入子目录失败: %v", err)
			return 0, err
		}
		id, err := newId.LastInsertId()
		if err != nil {
			logx.Errorf("获取子目录失败: %v", err)
			return 0, err
		}
		if v.AssetType == 1 {
			subdirectories, err := l.Subdirectories(id, v.Id, s)
			if err != nil {
				logx.Errorf("复制文件夹下面的文件到新的目录下失败: %v", err)
				return 0, err
			}
			use += subdirectories
		}
		use += v.AssetSize
	}
	return use, nil
}

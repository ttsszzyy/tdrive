package logic

import (
	"T-driver/common/errors"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearAssetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearAssetsLogic {
	return &ClearAssetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清理资源
func (l *ClearAssetsLogic) ClearAssets(in *pb.DelAssetsReq) (*pb.Response, error) {
	var use int64
	var uid int64
	list, err := l.svcCtx.AssetsModel.List(l.ctx, squirrel.Select().Where(squirrel.Eq{"id": in.Ids}))
	if err != nil {
		logx.Errorf("查询资源失败: %v", err)
		return nil, errors.DbError()
	}
	err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		logx.Infof("清理资源: %v", in.Ids)
		for _, v := range list {
			//清楚文件夹下面的所有文件
			if v.AssetType == 1 {
				subdirectories, err := l.Subdirectories(v.Id, session)
				if err != nil {
					logx.Errorf("获取子目录失败: %v", err)
					return err
				}
				use += subdirectories
			}
			//清理磁盘空间
			use += v.AssetSize

			if uid == 0 {
				uid = v.Uid
			}
		}
		err := l.svcCtx.AssetsModel.Clear(l.ctx, in.Ids, session)
		if err != nil {
			return err
		}
		if use > 0 {
			u, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, uid, 0)
			if err != nil {
				return err
			}
			u.StorageUse += -use
			if u.StorageUse < 0 {
				u.StorageUse = 0
			}
			err = l.svcCtx.UserStorageModel.Update(l.ctx, u, session)
			if err != nil {
				return err
			}
		}
		//todo 修改mongo文件状态为不可重试状态
		/*_, err = l.svcCtx.AssetFileModel.UpdateByAssetIds(l.ctx, in.Ids,6)
		if err != nil {
			logx.Errorf("删除文件失败: %v", err)
			return err
		}*/
		return nil
	})
	if err != nil {
		logx.Errorf("清理资源失败: %v", err)
		return nil, errors.DbError()
	}

	return &pb.Response{}, nil
}

func (l *ClearAssetsLogic) Subdirectories(old int64, s sqlx.Session) (use int64, err error) {
	list, err := l.svcCtx.AssetsModel.List(l.ctx, squirrel.Select().Where("pid = ? and deleted_time = ?", old, 0))
	if err != nil {
		logx.Errorf("获取子目录失败: %v", err)
		return 0, err
	}
	ids := make([]int64, 0, len(list))
	for _, v := range list {
		if v.AssetType == 1 {
			subdirectories, err := l.Subdirectories(v.Id, s)
			if err != nil {
				logx.Errorf("获取子目录失败: %v", err)
				return 0, err
			}
			use += subdirectories
		}
		use += v.AssetSize
		ids = append(ids, v.Id)
	}
	err = l.svcCtx.AssetsModel.Clear(l.ctx, ids, s)
	if err != nil {
		logx.Errorf("删除子目录失败: %v", err)
		return 0, err
	}
	return use, nil
}

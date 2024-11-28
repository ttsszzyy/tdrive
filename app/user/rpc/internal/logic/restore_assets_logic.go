package logic

import (
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/utils"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

type RestoreAssetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRestoreAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreAssetsLogic {
	return &RestoreAssetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 还原资源
func (l *RestoreAssetsLogic) RestoreAssets(in *pb.DelAssetsReq) (*pb.Response, error) {
	err := l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, id := range in.Ids {
			one, err := l.svcCtx.AssetsModel.FindOne(l.ctx, id, session)
			if err != nil {
				return err
			}
			one.DeletedTime = 0
			//校验文件名是否重复
			if one.AssetType == 1 {
				//计数器
				num := 0
			NEXT:
				count, err := l.svcCtx.AssetsModel.Count(l.ctx, squirrel.Select().Where("asset_name = ? and deleted_time = ? and pid = ?", one.AssetName, 0, one.Pid))
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
			err = l.svcCtx.AssetsModel.Update(l.ctx, one, session)
			if err != nil {
				return err
			}
			if one.AssetSize > 0 {
				u, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, one.Uid, 0)
				if err != nil {
					return err
				}
				if one.AssetSize > u.Storage-u.StorageUse {
					return fmt.Errorf("磁盘空间不足")
				}
				//添加空间
				u.StorageUse += one.AssetSize
				if u.StorageUse < 0 {
					u.StorageUse = 0
				}
				err = l.svcCtx.UserStorageModel.Update(l.ctx, u, session)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		logx.Errorf("还原资源失败: %v", err)
		return nil, err
	}

	return &pb.Response{}, nil
}

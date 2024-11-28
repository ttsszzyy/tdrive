package logic

import (
	"T-driver/common/errors"
	"T-driver/common/utils"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAssetsMoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAssetsMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAssetsMoveLogic {
	return &UpdateAssetsMoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 资源移动
func (l *UpdateAssetsMoveLogic) UpdateAssetsMove(in *pb.UpdateAssetsMoveReq) (*pb.Response, error) {
	err := l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, id := range in.Ids {
			one, err := l.svcCtx.AssetsModel.FindOne(l.ctx, id, session)
			if err != nil {
				return err
			}
			one.Pid = in.Pid
			one.UpdatedTime = time.Now().Unix()
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
			err = l.svcCtx.AssetsModel.Update(l.ctx, one, session)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logx.Errorf("资源移动失败: %v", err)
		return nil, errors.DbError()
	}

	return &pb.Response{}, nil
}

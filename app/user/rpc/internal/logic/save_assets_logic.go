package logic

import (
	"T-driver/app/user/model"
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

type SaveAssetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAssetsLogic {
	return &SaveAssetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存资源
func (l *SaveAssetsLogic) SaveAssets(in *pb.SaveAssetsReq) (*pb.SaveAssetsResp, error) {
	var id int64
	//校验文件名是否重复
	if in.AssetName != "" && in.AssetType != 1 {
		//计数器
		num := 0
	NEXT:
		count, err := l.svcCtx.AssetsModel.Count(l.ctx, squirrel.Select().Where("asset_name = ? and deleted_time = ? and pid = ?", in.AssetName, 0, in.Pid))
		if err != nil {
			return nil, err
		}
		if count >= 1 {
			num++
			if in.AssetType != 1 {
				filename, ext := utils.GetFileName(in.AssetName)
				if strings.HasSuffix(filename, fmt.Sprintf("(%d)", num-1)) {
					replace := utils.ReplaceLastOccurrence(filename, fmt.Sprintf("(%d)", num-1), fmt.Sprintf("(%d)", num))
					in.AssetName = fmt.Sprintf("%s%s", replace, ext)
					goto NEXT
				} else {
					in.AssetName = fmt.Sprintf("%s(%d)%s", filename, num, ext)
					goto NEXT
				}

			} else {
				//文件夹
				if strings.HasSuffix(in.AssetName, fmt.Sprintf("(%d)", num-1)) {
					replace := utils.ReplaceLastOccurrence(in.AssetName, fmt.Sprintf("(%d)", num-1), fmt.Sprintf("(%d)", num))
					in.AssetName = replace
					goto NEXT
				} else {
					in.AssetName = fmt.Sprintf("%s(%d)", in.AssetName, num)
					goto NEXT
				}
			}

		}
	}
	err := l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		asset, err := l.svcCtx.AssetsModel.Insert(l.ctx, &model.Assets{
			Uid:         in.Uid,
			Cid:         in.Cid,
			AssetName:   in.AssetName,
			AssetSize:   in.AssetSize,
			AssetType:   in.AssetType,
			TransitType: in.TransitType,
			IsTag:       in.IsTag,
			Pid:         in.Pid,
			Source:      in.Source,
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
			Status:      in.Status,
			IsReport:    2,
		}, session)
		if err != nil {
			return err
		}
		id, err = asset.LastInsertId()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logx.Error("保存资源失败:", err.Error())
		return nil, err
	}

	return &pb.SaveAssetsResp{Id: id}, nil
}

package document

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type CopyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCopyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CopyLogic {
	return &CopyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CopyLogic) Copy(req *types.CopyReq) (resp *types.Response, err error) {
	var (
		msg string
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
	)

	resp = &types.Response{}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{Ids: req.Ids})
	if err != nil {
		return nil, err
	}
	var size int64
	for _, asset := range assets.Assets {
		if asset.AssetType == 1 {
			use, err := l.Subdirectories(asset.Id)
			if err != nil {
				return nil, err
			}
			size += use
		}
		size += asset.AssetSize
	}
	//校验磁盘空间是否足够
	s, err := l.svcCtx.Rpc.FindOneUserStorage(l.ctx, &pb.FindOneUserStorageReq{Uid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	if size > s.SurStorage {
		switch lan {
		case errors.LanEn:
			msg = "Insufficient disk space"
		case errors.LanTw:
			msg = "磁盤空間不足"
		default:
			msg = "Insufficient disk space"
		}
		return nil, errors.CustomError(msg)
	}
	// todo: 上传文件
	_, err = l.svcCtx.Rpc.UpdateAssetsCopy(l.ctx, &pb.UpdateAssetsCopyReq{
		Ids:  req.Ids,
		Pid:  req.Pid,
		Cids: nil,
	})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	return
}

func (l *CopyLogic) Subdirectories(id int64) (use int64, err error) {
	list, err := l.svcCtx.Rpc.FindAssets(context.TODO(), &pb.FindAssetsReq{Pid: id})
	if err != nil {
		logx.Errorf("获取子目录失败: %v", id, err)
		return 0, err
	}
	for _, v := range list.Assets {
		if v.AssetType == 1 {
			subdirectories, err := l.Subdirectories(v.Id)
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

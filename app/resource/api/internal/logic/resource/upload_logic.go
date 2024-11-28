package resource

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/utils"
	"bytes"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"T-driver/app/resource/api/internal/svc"
	"T-driver/app/resource/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.UploadReq) (resp *types.Response, err error) {
	var (
		msg string
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
	)
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	//检查用户当前的上传数量
	maxUpload, err := l.svcCtx.Redis.ScardCtx(l.ctx, fmt.Sprintf("%s%d", model.UserUpload, userData.User.ID))
	if err != nil {
		switch lan {
		case errors.LanEn:
			msg = "Failed to retrieve the number of user uploads"
		case errors.LanTw:
			msg = "獲取用戶上傳數量失敗"
		default:
			msg = "Failed to retrieve the number of user uploads"
		}
		return nil, errors.CustomError(msg)
	}
	if maxUpload >= l.svcCtx.Config.MaxUpload {
		return nil, nil
	}

	// 校验文件大小不能超过限制
	if req.AssetSize > l.svcCtx.Config.MaxBytes {
		maxFileSizeMB := l.svcCtx.Config.MaxBytes / (1024 * 1024)
		switch lan {
		case errors.LanEn:
			msg = "The file size cannot exceed"
		case errors.LanTw:
			msg = "文件大小不能超過"
		default:
			msg = "The file size cannot exceed"
		}
		return nil, errors.CustomError(fmt.Sprintf("%s%d MB", msg, maxFileSizeMB))
	}
	//校验用户空间是否足够
	u, err := l.svcCtx.Rpc.FindOneUserStorage(l.ctx, &pb.FindOneUserStorageReq{Uid: userData.User.ID})
	if err != nil {
		switch lan {
		case errors.LanEn:
			msg = "Unable to retrieve user storage information"
		case errors.LanTw:
			msg = "無法獲取用戶存儲信息"
		default:
			msg = "Unable to retrieve user storage information"
		}
		return nil, errors.CustomError(msg)
	}
	if req.AssetSize > u.SurStorage {
		switch lan {
		case errors.LanEn:
			msg = "Insufficient space"
		case errors.LanTw:
			msg = "空間不足"
		default:
			msg = "Insufficient space"
		}
		return nil, errors.CustomError(msg)
	}
	//查询用户文件夹
	if req.Pid == 0 {
		pid, err := l.svcCtx.AssetsFolder(model.MyUpload, userData.User.ID)
		if err != nil {
			return nil, err
		}
		req.Pid = pid
	}

	asset, err := l.svcCtx.Rpc.SaveAssets(l.ctx, &pb.SaveAssetsReq{
		Uid:         userData.User.ID,
		AssetName:   req.AssetName,
		AssetSize:   req.AssetSize,
		AssetType:   utils.IsFileType(req.AssetName),
		TransitType: req.TransitType,
		IsTag:       model.AssetStatusAfoot,
		Pid:         req.Pid,
		Source:      req.Source, //本地上传
		Status:      model.AssetStatusAfoot,
	})
	if err != nil {
		switch lan {
		case errors.LanEn:
			msg = "Failed to save asset information"
		case errors.LanTw:
			msg = "保存資産信息失敗"
		default:
			msg = "Failed to save asset information"
		}
		return nil, errors.CustomError(msg)
	}
	// 将文件ID添加到上传集合中
	_, err = l.svcCtx.Redis.SaddCtx(l.ctx, fmt.Sprintf("%s%d", model.UserUpload, userData.User.ID), asset.Id)
	if err != nil {
		return nil, err
	}
	//设置过期时间
	err = l.svcCtx.Redis.Expire(fmt.Sprintf("%s%d", model.UserUpload, userData.User.ID), l.svcCtx.Config.UploadExpireTime)
	if err != nil {
		return nil, err
	}

	//本地上传
	go func(svcCtx *svc.ServiceContext, file []byte, id int64, name string, uid int64, assetSize int64) { // 使用匿名函数确保id不会被覆盖
		//刪除上传集合中的文件ID
		defer func(Redis *redis.Redis, key string, values ...any) {
			_, err := Redis.Srem(key, values)
			if err != nil {
				return
			}
		}(svcCtx.Redis, fmt.Sprintf("%s%d", model.UserUpload, uid), id)
		reader := bytes.NewReader(file)
		svcCtx.Upload(reader, name, id, assetSize)
	}(l.svcCtx, req.Flie, asset.Id, req.AssetName, userData.User.ID, req.AssetSize)
	return &types.Response{}, nil
}

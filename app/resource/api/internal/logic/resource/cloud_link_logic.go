package resource

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/utils"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"T-driver/app/resource/api/internal/svc"
	"T-driver/app/resource/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudLinkLogic {
	return &CloudLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudLinkLogic) CloudLink(req *types.CloudLinkReq) (resp *types.Response, err error) {
	var (
		msg string
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
	)

	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	var assetSize int64
	for _, link := range req.Link {
		// 创建一个HTTP请求来获取文件
		u, err := url.Parse(link)
		if err != nil {
			switch lan {
			case errors.LanEn:
				msg = "Failed to parse file"
			case errors.LanTw:
				msg = "解析文件失敗"
			default:
				msg = "Failed to parse file"
			}
			return nil, errors.CustomError(msg)
		}

		// 获取基本的文件名
		filename := path.Base(u.Path)
		// 创建一个HTTP HEAD请求
		request, err := http.NewRequest(http.MethodHead, link, nil)
		if err != nil {
			switch lan {
			case errors.LanEn:
				msg = "Download file failed"
			case errors.LanTw:
				msg = "下載文件失敗"
			default:
				msg = "Download file failed"
			}
			return nil, errors.CustomError(msg)
		}

		// 发送请求
		client := &http.Client{}
		res, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		req.AssetName = append(req.AssetName, filename)
		// 获取文件大小
		contentLength := res.Header.Get("Content-Length")
		fileSize, _ := strconv.ParseInt(contentLength, 10, 64)
		req.AssetSize = append(req.AssetSize, fileSize)
		//校验文件大小不能超过限制
		if fileSize > l.svcCtx.Config.MaxBytes {
			switch lan {
			case errors.LanEn:
				msg = "The file size cannot exceed"
			case errors.LanTw:
				msg = "文件大小不能超過"
			default:
				msg = "The file size cannot exceed"
			}
			return nil, errors.CustomError(msg + strconv.Itoa(int(l.svcCtx.Config.MaxBytes/1024/1024)) + "MB" + ",url:" + link)
		}
		assetSize += fileSize
	}
	//校验用户空间是否足够
	u, err := l.svcCtx.Rpc.FindOneUserStorage(l.ctx, &pb.FindOneUserStorageReq{Uid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	if assetSize > u.SurStorage {
		switch lan {
		case errors.LanEn:
			msg = "Insufficient disk space"
		case errors.LanTw:
			msg = "磁盤空間不足"
		default:
			msg = "Insufficient disk space"
		}
		return nil, errors.CustomError("磁盘空间不足")
	}
	//查询用户文件夹
	if req.Pid == 0 {
		pid, err := l.svcCtx.AssetsFolder(model.MyUpload, userData.User.ID)
		if err != nil {
			return nil, err
		}
		req.Pid = pid
	}

	for i, v := range req.Link {

		asset, err := l.svcCtx.Rpc.SaveAssets(l.ctx, &pb.SaveAssetsReq{
			Uid:         userData.User.ID,
			AssetName:   req.AssetName[i],
			AssetSize:   req.AssetSize[i],
			AssetType:   utils.IsFileType(req.AssetName[i]),
			IsTag:       model.AssetStatusAfoot,
			Pid:         req.Pid,
			Source:      req.Source,
			Status:      model.AssetStatusAfoot,
			TransitType: 3,
		})
		if err != nil {
			return nil, err
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
		req.Id = asset.Id
		//云上传
		go func(svc *svc.ServiceContext, link string, id int64, uid int64, assetSize int64) { // 使用匿名函数确保id不会被覆盖
			//刪除上传集合中的文件ID
			defer func(Redis *redis.Redis, key string, values ...any) {
				Redis.Srem(key, values)
			}(svc.Redis, fmt.Sprintf("%s%d", model.UserUpload, uid), id)
			svc.LinkUpload(link, id, assetSize)
		}(l.svcCtx, v, req.Id, u.Uid, req.AssetSize[i])
	}
	return &types.Response{}, nil
}

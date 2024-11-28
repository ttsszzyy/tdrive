package document

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"context"
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

// CloudLink 处理文件上传逻辑 ------废弃
// req: 上传请求体，包含文件上传的相关信息。
// 返回值:
// - resp: 响应体，包含上传操作的结果信息。
// - err: 错误信息，如果上传过程中出现错误。
func (l *CloudLinkLogic) CloudLink(req *types.CloudLinkReq) (resp *types.Response, err error) {
	//获取用户信息
	/*userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	var assetSize int64
	for _, link := range req.Link {
		// 创建一个HTTP请求来获取文件
		u, err := url.Parse(link)
		if err != nil {
			return nil, errors.CustomError("解析文件失败")
		}

		// 获取基本的文件名
		filename := path.Base(u.Path)
		// 创建一个HTTP HEAD请求
		request, err := http.NewRequest("HEAD", link, nil)
		if err != nil {
			return nil, errors.CustomError("下载文件失败")
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
			return nil, errors.CustomError("文件大小不能超过" + strconv.Itoa(int(l.svcCtx.Config.MaxBytes/1024/1024)) + "MB" + ",url:" + link)
		}
		assetSize += fileSize
	}
	//校验用户空间是否足够
	s, err := l.svcCtx.Rpc.FindOneUserStorage(l.ctx, &pb.FindOneUserStorageReq{Uid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	if assetSize > s.SurStorage {
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
		}(l.svcCtx, v, req.Id, userData.User.ID, req.AssetSize[i])
	}*/
	return
}

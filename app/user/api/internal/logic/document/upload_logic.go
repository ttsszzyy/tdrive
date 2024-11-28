package document

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

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

// Upload 处理文件上传逻辑----废弃
// req: 上传请求体，包含文件上传的相关信息。
// 返回值:
// - resp: 响应体，包含上传操作的结果信息。
// - err: 错误信息，如果上传过程中出现错误。
func (l *UploadLogic) Upload(req *types.UploadReq) (resp *types.Response, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	if req.Id == 0 {
		return nil, errors.ErrorNotFound(lan)
	}
	//检查用户当前的上传数量
	/*maxUpload, err := l.svcCtx.Redis.ScardCtx(l.ctx, fmt.Sprintf("%s%d", model.UserUpload, userData.User.ID))
	if err != nil {
		return nil, errors.CustomError("获取用户上传数量失败")
	}
	u, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.CustomError("获取用户信息失败")
	}
	if maxUpload >= l.svcCtx.Config.MaxUpload {
		text := "您正在上傳多個檔案，可能會耗費較長時間。 請確保網絡穩定，以免上傳中斷影響你的使用體驗！"
		if u.LanguageCode == "en" {
			text = "You're uploading multiple files, which may take some time. Ensure your network connection is stable to avoid interruptions that could affect your upload process!"
		}
		return nil, errors.NewErrCodeMsg(errors.ErrUserUploadLimit, text)
	}*/
	//种子上传
	/*if req.TransitType == 7 {
		tf, err := torrent.ParseFile(bytes.NewReader(req.Flie))
		if err != nil {
			return nil, errors.CustomError("解析种子失败")
		}
		req.AssetName = tf.FileName
		req.AssetSize = int64(tf.FileLen)
		// random peerId
		var peerId [torrent.IDLEN]byte
		_, _ = rand.Read(peerId[:])
		//connect tracker & find peers
		peers := torrent.FindPeers(tf, peerId)
		if len(peers) == 0 {
			return nil, errors.CustomError("解析种子失败")
		}
		// build torrent task
		task := &torrent.TorrentTask{
			PeerId:   peerId,
			PeerList: peers,
			InfoSHA:  tf.InfoSHA,
			FileName: tf.FileName,
			FileLen:  tf.FileLen,
			PieceLen: tf.PieceLen,
			PieceSHA: tf.PieceSHA,
		}
		//download from peers & make file
		req.Flie, err = torrent.Download(task)
		if err != nil {
			fmt.Println("download error")
			return nil, err
		}

	}*/

	// 校验文件大小不能超过限制
	/*if req.AssetSize > l.svcCtx.Config.MaxBytes {
		maxFileSizeMB := l.svcCtx.Config.MaxBytes / (1024 * 1024)
		return nil, errors.CustomError(fmt.Sprintf("文件大小不能超过%d MB", maxFileSizeMB))
	}
	//校验用户空间是否足够
	s, err := l.svcCtx.Rpc.FindOneUserStorage(l.ctx, &pb.FindOneUserStorageReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.CustomError("无法获取用户存储信息")
	}
	if req.AssetSize > s.SurStorage {
		return nil, errors.NewErrCodeMsg(errors.ErrCodeNotSpace, "空间不足")
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
		return nil, errors.CustomError("保存资产信息失败")
	}*/
	/*assets, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	// 将文件ID添加到上传集合中
	_, err = l.svcCtx.Redis.SaddCtx(l.ctx, fmt.Sprintf("%s%d", model.UserUpload, userData.User.ID), assets.Id)
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
	}(l.svcCtx, req.Flie, req.Id, req.AssetName, userData.User.ID, req.AssetSize)*/
	return &types.Response{}, nil
}

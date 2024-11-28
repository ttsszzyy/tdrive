package transmission

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/utils"
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAssetFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAssetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAssetFileLogic {
	return &CreateAssetFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAssetFileLogic) CreateAssetFile(req *types.CreateAssetFileReq) (resp *types.CreateAssetFileResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}

	//检查用户当前的上传数量
	maxUpload, err := l.svcCtx.Redis.ScardCtx(l.ctx, fmt.Sprintf("%s%d", model.UserUpload, userData.User.ID))
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
	}
	// 校验文件大小不能超过限制
	if req.AssetSize > l.svcCtx.Config.MaxBytes {
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

	asset, err := l.svcCtx.Rpc.SaveAssetFile(l.ctx, &pb.SaveAssetFileReq{
		Uid:       userData.User.ID,
		AssetName: req.AssetName,
		AssetSize: req.AssetSize,
		AssetType: utils.IsFileType(req.AssetName),
		Pid:       req.Pid,
		Source:    req.Source, //本地上传
		Path:      req.Path,
	})
	if err != nil {
		return nil, errors.CustomError("保存资产信息失败")
	}
	return &types.CreateAssetFileResp{Id: asset.Id}, nil
}

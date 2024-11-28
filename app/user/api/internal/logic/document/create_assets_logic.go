package document

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

type CreateAssetsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAssetsLogic {
	return &CreateAssetsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAssetsLogic) CreateAssets(req *types.CreateAssetsReq) (resp *types.CreateAssetsResp, err error) {
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
	u, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		switch lan {
		case errors.LanEn:
			msg = "Failed to obtain user information"
		case errors.LanTw:
			msg = "獲取用戶信息失敗"
		default:
			msg = "Failed to obtain user information"
		}
		return nil, errors.CustomError(msg)
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
	s, err := l.svcCtx.Rpc.FindOneUserStorage(l.ctx, &pb.FindOneUserStorageReq{Uid: userData.User.ID})
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
	if req.AssetSize > s.SurStorage {
		switch lan {
		case errors.LanEn:
			msg = "Insufficient space"
		case errors.LanTw:
			msg = "空間不足"
		default:
			msg = "Insufficient space"
		}
		return nil, errors.NewErrCodeMsg(errors.ErrCodeNotSpace, msg)
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
	return &types.CreateAssetsResp{Id: asset.Id}, nil
}

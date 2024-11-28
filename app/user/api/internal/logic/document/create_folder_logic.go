package document

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFolderLogic {
	return &CreateFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFolderLogic) CreateFolder(req *types.CreateFolderReq) (resp *types.Response, err error) {
	var (
		msg string
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
	)
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	if req.Pid == 0 {
		req.Pid = 1
	}
	//查询文件夹是否存在
	count, err := l.svcCtx.Rpc.CountAssets(l.ctx, &pb.CountAssetsReq{AssetName: req.AssetName, Uid: userData.User.ID, Pid: req.Pid})
	if err != nil {
		return nil, err
	}
	if count.Total > 0 {
		switch lan {
		case errors.LanEn:
			msg = "The file name cannot be duplicated"
		case errors.LanTw:
			msg = "文件名稱不能重複"
		default:
			msg = "The file name cannot be duplicated"
		}
		return nil, errors.CustomError(msg)
	}
	//保存文件夹
	_, err = l.svcCtx.Rpc.SaveAssets(l.ctx, &pb.SaveAssetsReq{
		Uid:         userData.User.ID,
		AssetName:   req.AssetName,
		AssetType:   1,
		TransitType: 1,
		IsTag:       2,
		Pid:         req.Pid,
		Source:      1,
		Status:      model.AssetStatusEnable,
	})
	if err != nil {
		return nil, err
	}

	return
}

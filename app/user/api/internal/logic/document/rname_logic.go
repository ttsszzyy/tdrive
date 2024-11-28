package document

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RnameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRnameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RnameLogic {
	return &RnameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RnameLogic) Rname(req *types.RnameReq) (resp *types.Response, err error) {
	var (
		msg string
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
	)
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	assets, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	//查询文件夹是否存在
	count, err := l.svcCtx.Rpc.CountAssets(l.ctx, &pb.CountAssetsReq{AssetName: req.AssetName, Uid: userData.User.ID, Pid: assets.Pid})
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
	_, err = l.svcCtx.Rpc.UpdateAssetsName(l.ctx, &pb.UpdateAssetsNameReq{
		Id:        req.Id,
		AssetName: req.AssetName,
	})
	if err != nil {
		return nil, err
	}
	// 调用sdk修改文件名
	tcli, err := l.svcCtx.GetStorage(userData.User.ID)
	if err != nil {
		logx.Errorf("get titan storage client error:%v", err)
	}
	if tcli != nil {
		if err := tcli.RenameAsset(l.ctx, assets.Cid, req.AssetName); err != nil {
			logx.Errorf("dail titan storeage client rename asset error:%w", err)
		}
	}
	return resp, nil
}

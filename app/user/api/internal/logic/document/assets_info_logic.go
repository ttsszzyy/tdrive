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

type AssetsInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssetsInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetsInfoLogic {
	return &AssetsInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssetsInfoLogic) AssetsInfo(req *types.AssetsInfoReq) (resp *types.AssetItem, err error) {
	var (
		msg string
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
	)

	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(msg)
	}
	one, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	if one.Id == 0 {
		switch lan {
		case errors.LanEn:
			msg = "The resource does not exist"
		case errors.LanTw:
			msg = "資源不存在"
		default:
			msg = "The resource does not exist"
		}
		return nil, errors.NewErrCodeMsg(10001, msg)
	}
	passet, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{Id: one.Pid})
	if err != nil {
		return nil, err
	}
	//获取资源链接
	url := one.Link
	if one.AssetType != 1 && one.Link == "" {
		/*shareAssetResult, err := l.svcCtx.Storage.GetURL(l.ctx, one.Cid)
		if err != nil {
			logx.Error("获取资源链接失败：", err)
			return nil, errors.NewErrCodeMsg(10001, err.Error())
		}
		if shareAssetResult != nil && len(shareAssetResult.URLs) > 0 {
			url = shareAssetResult.URLs[0]
			l.svcCtx.Rpc.UpdateAssetsName(l.ctx, &pb.UpdateAssetsNameReq{Id: one.Id, Link: url})
		}*/
	}

	resp = &types.AssetItem{
		Uid:         one.Uid,
		Cid:         one.Cid,
		TransitType: one.TransitType,
		AssetName:   one.AssetName,
		AssetSize:   one.AssetSize,
		AssetType:   one.AssetType,
		CreatedTime: one.CreatedTime,
		Pid:         one.Pid,
		Pname:       passet.AssetName,
		UpdatedTime: one.UpdatedTime,
		IsTag:       one.IsTag,
		Source:      one.Source,
		Url:         url,
		Id:          one.Id,
		Status:      one.Status,
		IsReport:    one.IsReport,
	}
	return resp, nil
}

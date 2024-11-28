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

type ReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportLogic {
	return &ReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportLogic) Report(req *types.ReportReq) (resp *types.Response, err error) {
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	assets, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{Cid: req.Cid})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.Rpc.UpdateAssetsName(l.ctx, &pb.UpdateAssetsNameReq{
		Id:         assets.Id,
		IsReport:   1,
		ReportType: req.ReportType,
	})
	if err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}

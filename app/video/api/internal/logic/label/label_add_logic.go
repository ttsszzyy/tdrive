package label

import (
	"T-driver/app/video/rpc/videoPb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LabelAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLabelAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelAddLogic {
	return &LabelAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LabelAddLogic) LabelAdd(req *types.LabelAddReq) (resp *types.LabelAddResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	createLabel, err := l.svcCtx.VideoRpc.CreateLabel(l.ctx, &videoPb.CreateLabelReq{
		Title: req.Title,
		Vid:   req.Vid,
		Uid:   userData.User.ID,
	})
	return &types.LabelAddResp{Repeat: createLabel.Repeat}, err
}

package comment

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

type CommentAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentAddLogic {
	return &CommentAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentAddLogic) CommentAdd(req *types.CommentAddReq) (resp *types.Response, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	_, err = l.svcCtx.VideoRpc.CreateComment(l.ctx, &videoPb.CreateCommentReq{
		Content: req.Content,
		Vid:     req.Vid,
		Uid:     userData.User.ID,
	})
	return &types.Response{}, err
}

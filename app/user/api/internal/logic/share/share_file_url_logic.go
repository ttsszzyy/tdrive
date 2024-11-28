package share

import (
	"context"
	"fmt"
	"time"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

// ShareFileUrlLogic 获取分享文件链接地址
type ShareFileUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewShareFileUrlLogic 新建 获取分享文件链接地址
func NewShareFileUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareFileUrlLogic {
	return &ShareFileUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ShareFileUrl 实现 获取分享文件链接地址
func (l *ShareFileUrlLogic) ShareFileUrl(req *types.GetShareURLReq) (resp string, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	// 获取分享信息
	logx.Error("uuid", req.UUID)
	share, err := l.svcCtx.Rpc.FindOneShare(l.ctx, &pb.FindOneShareReq{Uuid: req.UUID})
	if err != nil {
		return "", errors.SystemError(lan)
	}
	if share.Id == 0 || (share.EffectiveTime > 0 && share.EffectiveTime < time.Now().Unix()) || share.DeletedTime > 0 {
		return "", errors.CustomError("The shared file has expired")
	}

	return share.Link, nil
}

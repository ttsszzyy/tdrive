package reply

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/app/video/rpc/videoPb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplyAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyAddLogic {
	return &ReplyAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplyAddLogic) ReplyAdd(req *types.ReplyAddReq) (resp *types.VideoReply, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	reply, err := l.svcCtx.VideoRpc.CreateReply(l.ctx, &videoPb.CreateReplyReq{
		CommentId:  req.CommentId,
		Content:    req.Content,
		FromUserId: req.FromUserId,
		ToUserId:   userData.User.ID,
	})
	userList, err := l.svcCtx.UserRpc.FindUser(l.ctx, &pb.QueryUserReq{
		Uid: []int64{reply.ToUserId, reply.FromUserId},
	})
	if err != nil {
		return nil, err
	}
	umap := make(map[int64]*pb.User)
	for _, u := range userList.Users {
		umap[u.Uid] = u
	}
	u, ok := umap[reply.FromUserId]
	fromUserName := ""
	fromUserAvatar := ""
	if ok {
		fromUserName = u.Name
		fromUserAvatar = u.Avatar
	}
	toU, ok := umap[reply.ToUserId]
	toUserName := ""
	toUserAvatar := ""
	if ok {
		toUserName = toU.Name
		toUserAvatar = toU.Avatar
	}
	return &types.VideoReply{
		CommentId:      reply.CommentId,
		Content:        reply.Content,
		CreatedTime:    reply.CreatedTime,
		FromUserId:     reply.FromUserId,
		FromUserName:   fromUserName,
		FromUserAvatar: fromUserAvatar,
		Id:             reply.Id,
		ToUserId:       reply.ToUserId,
		ToUserName:     toUserName,
		ToUserAvatar:   toUserAvatar,
	}, err
}

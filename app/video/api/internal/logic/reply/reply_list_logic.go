package reply

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/app/video/rpc/videoPb"
	"context"

	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyListLogic {
	return &ReplyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplyListLogic) ReplyList(req *types.ReplyListReq) (resp *types.ReplyListResp, err error) {
	if req.Size > 100 {
		req.Size = 100
	}
	list, err := l.svcCtx.VideoRpc.GetReplyList(l.ctx, &videoPb.GetReplyListReq{
		CommentId: req.Commentid,
		Page:      req.Page,
		Size:      req.Size,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.ReplyListResp{
		List:  make([]*types.VideoReply, 0, len(list.List)),
		Total: list.Total,
	}
	ids := make([]int64, 0, len(list.List)*2)
	for _, v := range list.List {
		ids = append(ids, v.FromUserId)
		ids = append(ids, v.ToUserId)
	}
	userList, err := l.svcCtx.UserRpc.FindUser(l.ctx, &pb.QueryUserReq{
		Uid: ids,
	})
	if err != nil {
		return nil, err
	}
	umap := make(map[int64]*pb.User)
	for _, u := range userList.Users {
		umap[u.Uid] = u
	}
	for _, v := range list.List {
		u, ok := umap[v.FromUserId]
		fromUserName := ""
		fromUserAvatar := ""
		if ok {
			fromUserName = u.Name
			fromUserAvatar = u.Avatar
		}
		toU, ok := umap[v.ToUserId]
		toUserName := ""
		toUserAvatar := ""
		if ok {
			toUserName = toU.Name
			toUserAvatar = toU.Avatar
		}
		resp.List = append(resp.List, &types.VideoReply{
			CommentId:      v.CommentId,
			Content:        v.Content,
			CreatedTime:    v.CreatedTime,
			FromUserId:     v.FromUserId,
			FromUserName:   fromUserName,
			FromUserAvatar: fromUserAvatar,
			Id:             v.Id,
			ToUserId:       v.ToUserId,
			ToUserName:     toUserName,
			ToUserAvatar:   toUserAvatar,
		})
	}

	return resp, nil
}

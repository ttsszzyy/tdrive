package comment

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/app/video/rpc/videoPb"
	"T-driver/common/errors"
	"context"

	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	if req.Size > 100 {
		req.Size = 100
	}
	list, err := l.svcCtx.VideoRpc.GetCommentList(l.ctx, &videoPb.GetCommentListReq{
		Vid:  req.Vid,
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	resp = &types.CommentListResp{
		List:  make([]*types.VideoComment, 0, len(list.List)),
		Total: list.Total,
	}
	ids := make([]int64, 0, len(list.List))
	for _, v := range list.List {
		ids = append(ids, v.Uid)
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
		resp.List = append(resp.List, &types.VideoComment{
			Content:     v.Content,
			CreatedTime: v.CreatedTime,
			Id:          v.Id,
			Uid:         v.Uid,
			Vid:         v.Vid,
			UserName:    umap[v.Uid].Name,
			Avatar:      umap[v.Uid].Avatar,
			ReplyCount:  v.ReplyCount,
		})
	}

	return resp, nil
}

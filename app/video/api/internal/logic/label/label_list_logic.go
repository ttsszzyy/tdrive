package label

import (
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/videoPb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"strconv"

	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LabelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLabelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelListLogic {
	return &LabelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LabelListLogic) LabelList(req *types.LabelListReq) (resp *types.LabelListResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	list, err := l.svcCtx.VideoRpc.GetLabelList(l.ctx, &videoPb.GetLabelListReq{
		Vid:  req.Vid,
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	lid := make([]string, 0, len(list.List))
	labelUserCountMap := make(map[string]*videoPb.CountLabelUserItem)
	for _, v := range list.List {
		lid = append(lid, v.Id)
		//获取标签点赞信息
		labelLikeNum, err := l.svcCtx.Redis.GetCtx(l.ctx, fmt.Sprintf("%s%v", model.LabelLikeKey, v.Id))
		if err != nil {
			return nil, errors.CustomError("获取标签点赞数量失败")
		}
		i, _ := strconv.ParseInt(labelLikeNum, 10, 64)
		labelNoLikeNum, err := l.svcCtx.Redis.GetCtx(l.ctx, fmt.Sprintf("%s%v", model.LabelNoLikeKey, v.Id))
		if err != nil {
			return nil, errors.CustomError("获取标签点踩数量失败")
		}
		noi, _ := strconv.ParseInt(labelNoLikeNum, 10, 64)
		labelUserCountMap[v.Id] = &videoPb.CountLabelUserItem{
			Lid:        v.Id,
			LikesNum:   i,
			NoLikesNum: noi,
		}
	}
	labelUser := &videoPb.GetLabelUserResp{}
	if len(lid) > 0 {
		//获取用户标签点赞信息
		labelUser, err = l.svcCtx.VideoRpc.GetLabelUser(l.ctx, &videoPb.GetLabelUserReq{
			Lid: lid,
			Uid: userData.User.ID,
		})
		if err != nil {
			return nil, errors.FromRpcError(err)
		}
	}
	labelUserMap := make(map[string]*videoPb.VideosLabelUser)
	for _, v := range labelUser.List {
		labelUserMap[v.Lid] = v
	}

	res := make([]*types.Label, 0, len(list.List))
	for _, v := range list.List {
		var likes, noLikes bool
		var labelUserId string
		videosLabelUser, ok := labelUserMap[v.Id]
		if ok {
			labelUserId = videosLabelUser.Id
			likes = videosLabelUser.Likes
			noLikes = videosLabelUser.NoLikes
		}
		res = append(res, &types.Label{
			Lid:         v.Id,
			LabelUserId: labelUserId,
			Vid:         v.Vid,
			Uid:         v.Uid,
			Title:       v.Title,
			Likes:       likes,
			NoLikes:     noLikes,
			LikesNum:    labelUserCountMap[v.Id].LikesNum,
			NoLikesNum:  labelUserCountMap[v.Id].NoLikesNum,
			CreatedTime: v.CreatedTime,
		})
	}

	return &types.LabelListResp{
		List:  res,
		Total: list.Total,
	}, nil
}

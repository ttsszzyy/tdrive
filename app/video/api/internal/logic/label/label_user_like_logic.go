package label

import (
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/videoPb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/mr"
	"strconv"

	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LabelUserLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLabelUserLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelUserLikeLogic {
	return &LabelUserLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LabelUserLikeLogic) LabelUserLike(req *types.LabelUserLikeReq) (resp *types.LabelUserLikeResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	res, err := l.svcCtx.VideoRpc.SaveLabelUser(l.ctx, &videoPb.SaveLabelUserReq{
		Id:      req.LabelUserId,
		Lid:     req.Lid,
		Uid:     userData.User.ID,
		Likes:   req.Likes,
		NoLikes: req.NoLikes,
	})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	labelUser := &videoPb.VideosLabelUser{}
	label := &videoPb.VideosLabel{}
	labelUserCount := &videoPb.CountLabelUserItem{}
	err = mr.Finish(func() error {
		labelUser, err = l.svcCtx.VideoRpc.GetOneLabelUser(l.ctx, &videoPb.GetOneLabelUserReq{Id: res.Id})
		if err != nil {
			return errors.FromRpcError(err)
		}
		return nil
	}, func() error {
		label, err = l.svcCtx.VideoRpc.GetOneLabel(l.ctx, &videoPb.GetOneLabelReq{Id: req.Lid})
		if err != nil {
			return errors.FromRpcError(err)
		}
		return nil
	}, func() error {
		//获取标签总数
		//获取标签点赞信息
		labelLikeNum, err := l.svcCtx.Redis.GetCtx(l.ctx, fmt.Sprintf("%s%v", model.LabelLikeKey, req.Lid))
		if err != nil {
			return errors.CustomError("获取标签点赞数量失败")
		}
		i, _ := strconv.ParseInt(labelLikeNum, 10, 64)
		labelNoLikeNum, err := l.svcCtx.Redis.GetCtx(l.ctx, fmt.Sprintf("%s%v", model.LabelNoLikeKey, req.Lid))
		if err != nil {
			return errors.CustomError("获取标签点踩数量失败")
		}
		noi, _ := strconv.ParseInt(labelNoLikeNum, 10, 64)
		labelUserCount.LikesNum = i
		labelUserCount.NoLikesNum = noi
		labelUserCount.Lid = req.Lid
		return nil
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	return &types.LabelUserLikeResp{
		Lid:         labelUser.Lid,
		LabelUserId: labelUser.Id,
		Vid:         label.Vid,
		Uid:         labelUser.Uid,
		Title:       label.Title,
		Likes:       labelUser.Likes,
		NoLikes:     labelUser.NoLikes,
		LikesNum:    labelUserCount.LikesNum,
		NoLikesNum:  labelUserCount.NoLikesNum,
		CreatedTime: labelUser.CreatedTime,
	}, nil
}

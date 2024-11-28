package logic

import (
	"T-driver/app/video/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveLabelUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveLabelUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveLabelUserLogic {
	return &SaveLabelUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标签点赞
func (l *SaveLabelUserLogic) SaveLabelUser(in *videoPb.SaveLabelUserReq) (*videoPb.SaveLabelUserResp, error) {
	id := in.Id
	if in.Id != "" {
		one, err := l.svcCtx.VideosLabelUserModel.FindOne(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
		one.Likes = in.Likes
		one.NoLikes = in.NoLikes
		_, err = l.svcCtx.VideosLabelUserModel.Update(l.ctx, one)
		if err != nil {
			return nil, err
		}
	} else {
		objectID := primitive.NewObjectID()
		err := l.svcCtx.VideosLabelUserModel.Insert(l.ctx, &model.Videos_label_user{
			ID:       objectID,
			Uid:      in.Uid,
			Lid:      in.Lid,
			Likes:    in.Likes,
			NoLikes:  in.NoLikes,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
		id = objectID.Hex()
	}
	//校验用户是否点赞
	likesismemberCtx, err := l.svcCtx.Redis.SismemberCtx(l.ctx, fmt.Sprintf("%s%v", model.UserLabelLikeKey, in.Lid), in.Uid)
	if err != nil {
		return nil, err
	}
	//更新点赞缓存
	if in.Likes {
		//用户点赞
		_, err := l.svcCtx.Redis.SaddCtx(l.ctx, fmt.Sprintf("%s%v", model.UserLabelLikeKey, in.Lid), in.Uid)
		if err != nil {
			return nil, err
		}
		if !likesismemberCtx {
			//点赞数量+1
			_, err = l.svcCtx.Redis.IncrCtx(l.ctx, fmt.Sprintf("%s%v", model.LabelLikeKey, in.Lid))
			if err != nil {
				return nil, err
			}
		}
	} else {
		// 取消点赞
		_, err = l.svcCtx.Redis.SremCtx(l.ctx, fmt.Sprintf("%s%v", model.UserLabelLikeKey, in.Lid), in.Uid)
		if err != nil {
			return nil, err
		}
		if likesismemberCtx {
			//点赞数量-1
			_, err := l.svcCtx.Redis.DecrCtx(l.ctx, fmt.Sprintf("%s%v", model.LabelLikeKey, in.Lid))
			if err != nil {
				return nil, err
			}
		}
	}

	//校验用户是否点踩
	sismemberCtx, err := l.svcCtx.Redis.SismemberCtx(l.ctx, fmt.Sprintf("%s%v", model.UserLabelNoLikeKey, in.Lid), in.Uid)
	if err != nil {
		return nil, err
	}
	if in.NoLikes {
		//用户点踩
		_, err := l.svcCtx.Redis.SaddCtx(l.ctx, fmt.Sprintf("%s%v", model.UserLabelNoLikeKey, in.Lid), in.Uid)
		if err != nil {
			return nil, err
		}
		//更新点踩缓存
		if !sismemberCtx {
			//点踩数量+1
			_, err = l.svcCtx.Redis.IncrCtx(l.ctx, fmt.Sprintf("%s%v", model.LabelNoLikeKey, in.Lid))
			if err != nil {
				return nil, err
			}
		}
	} else {
		// 取消点踩
		_, err = l.svcCtx.Redis.SremCtx(l.ctx, fmt.Sprintf("%s%v", model.UserLabelNoLikeKey, in.Lid), in.Uid)
		if err != nil {
			return nil, err
		}
		if sismemberCtx {
			//点踩数量-1
			_, err := l.svcCtx.Redis.DecrCtx(l.ctx, fmt.Sprintf("%s%v", model.LabelNoLikeKey, in.Lid))
			if err != nil {
				return nil, err
			}
		}
	}

	return &videoPb.SaveLabelUserResp{Id: id}, nil
}

package video

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/videoPb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"strconv"
	"strings"

	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoListLogic {
	return &VideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoListLogic) VideoList(req *types.VideoListReq) (resp *types.VideoListResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	if req.Size > 20 {
		req.Size = 5
	}
	//第一次获取随机获取视频数据
	if req.Page == 1 {
		//删除集合
		l.svcCtx.Redis.DelCtx(l.ctx, fmt.Sprintf("%s%v", model.VideoListUser, userData.User.ID))

		//获取短视频集合是否存在
		exists, err := l.svcCtx.Redis.ExistsCtx(l.ctx, model.VideoList)
		if err != nil {
			return nil, err
		}
		if !exists {
			//获取所有的视频id，并且保存到集合中
			list, err := l.svcCtx.VideoRpc.GetVideoList(l.ctx, &videoPb.GetVideoListReq{
				Status: 3,
			})
			if err != nil {
				return nil, errors.FromRpcError(err)
			}
			ids := make([]any, 0, len(list.VideoList))
			for _, v := range list.VideoList {
				ids = append(ids, v.Id)
			}
			l.svcCtx.Redis.SaddCtx(l.ctx, model.VideoList, ids...)
		}
		// 获取集合中的所有元素
		members, err := l.svcCtx.Redis.SmembersCtx(l.ctx, model.VideoList)
		if err != nil {
			return nil, err
		}
		vids := make([]any, 0, len(members))
		for _, member := range members {
			vids = append(vids, member)
		}
		_, err = l.svcCtx.Redis.SaddCtx(l.ctx, fmt.Sprintf("%s%v", model.VideoListUser, userData.User.ID), vids...)
		if err != nil {
			return nil, err
		}
	}

	//校验剩余数量
	cardinality, err := l.svcCtx.Redis.ScardCtx(l.ctx, fmt.Sprintf("%s%v", model.VideoListUser, userData.User.ID))
	if err != nil {
		return nil, err
	}
	if req.Size > cardinality {
		req.Size = cardinality
		// 获取集合中的所有元素
		/*members, err := l.svcCtx.Redis.SmembersCtx(l.ctx, model.VideoList)
		if err != nil {
			return nil, err
		}
		vids := make([]any, 0, len(members))
		for _, member := range members {
			vids = append(vids, member)
		}
		_, err = l.svcCtx.Redis.SaddCtx(l.ctx, fmt.Sprintf("%s%v", model.VideoListUser, userData.User.ID), vids...)
		if err != nil {
			return nil, err
		}*/
	}
	//随机获取视频数量
	vids := make([]int64, 0, req.Size)
	for i := 0; i < int(req.Size); i++ {
		total, err := l.svcCtx.Redis.SpopCtx(l.ctx, fmt.Sprintf("%s%v", model.VideoListUser, userData.User.ID))
		if err != nil {
			fmt.Println("spop err :", err)
			return nil, err
		}
		id, _ := strconv.ParseInt(total, 10, 64)
		vids = append(vids, id)
	}

	list, err := l.svcCtx.VideoRpc.GetVideoList(l.ctx, &videoPb.GetVideoListReq{
		Ids:    vids,
		Size:   req.Size,
		Status: 3,
	})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	ids := make([]int64, 0, len(list.VideoList))
	for _, v := range list.VideoList {
		ids = append(ids, v.Uid)
	}
	admins, err := l.svcCtx.UserRpc.FindAdminByIdAccountIsDisableDeletedTime(l.ctx, &pb.AdminReq{Id: ids})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	adminsMap := make(map[int64]*pb.Admin)
	for _, v := range admins.Admins {
		adminsMap[v.Id] = v
	}
	List := make([]*types.Video, 0)
	for _, v := range list.VideoList {
		var fp, avatar string
		if v.Uid != 0 {
			admin, ok := adminsMap[v.Uid]
			if ok {
				avatar = admin.Avatar
			}
		}
		fp = fmt.Sprintf("%s%s", l.svcCtx.Config.S3.CDN, strings.TrimPrefix(v.FilePath, l.svcCtx.Config.S3.Endpoint))
		List = append(List, &types.Video{
			Id:          v.Id,
			Title:       v.Title,
			Desc:        v.Desc,
			FilePath:    fp,
			CreatedTime: v.CreatedTime,
			Avatar:      avatar,
		})
	}
	//获取短视频的总数
	total, err := l.svcCtx.Redis.ScardCtx(l.ctx, model.VideoList)
	if err != nil {
		return nil, err
	}
	return &types.VideoListResp{
		List:  List,
		Total: total,
	}, nil
}

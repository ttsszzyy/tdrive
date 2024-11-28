package user

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type RewardListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRewardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RewardListLogic {
	return &RewardListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RewardListLogic) RewardList(req *types.RewardListReq) (resp *types.RewardListRes, err error) {
	//获取用户信息
	/*userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	if req.Size > 100 {
		req.Size = 100
	}*/
	/*//获取推荐用户数据
	UserList, err := l.svcCtx.Rpc.FindUserByPid(l.ctx, &pb.FindUserByPidReq{Pid: userData.User.ID})
	if err != nil {
		return nil, errors.DbError()
	}
	resp.UserTotal = UserList.Total
	//获取引荐积分
	userPoints, err := l.svcCtx.Rpc.FindUserPointsByUid(l.ctx, &pb.UidSourceReq{Uid: userData.User.ID, Source: model.PointsRecommend})
	if err != nil {
		return nil, errors.DbError()
	}
	resp.PointsTotal = userPoints.Total
	//获取引荐空间
	userStorage, err := l.svcCtx.Rpc.FindUserStorageByUid(l.ctx, &pb.UidSourceReq{Uid: userData.User.ID, Source: model.StorageRecommend})
	if err != nil {
		return nil, errors.DbError()
	}
	resp.StorageTotal = userStorage.Total

	//获取推荐用户数据分页
	list, err := l.svcCtx.Rpc.FindUserByNameIsDisable(l.ctx, &pb.QueryUserReq{Pid: userData.User.ID, Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, errors.DbError()
	}
	for _, info := range list.Users {
		resp.Rewards = append(resp.Rewards, &types.Reward{
			Uid:         info.Uid,
			Name:        info.Name,
			Source:      info.Source,
			CreatedTime: info.CreatedTime,
		})
	}
	resp.Total = UserList.Total*/

	/*generateFunc := func(source chan<- int64) {
		source <- userData.User.ID
	}

	mapperFunc := func(item int64, writer mr.Writer[*types.RewardListRes], cancel func(error)) {
		resp = &types.RewardListRes{Rewards: make([]*types.Reward, 0)}
		//获取推荐用户数据
		UserList, err := l.svcCtx.Rpc.FindUserByPid(l.ctx, &pb.FindUserByPidReq{Pid: userData.User.ID})
		if err != nil {
			cancel(errors.DbError())
			return
		}
		u, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
		if err != nil {
			cancel(errors.DbError())
			return
		}
		resp.UserTotal = UserList.Total
		resp.Total = UserList.Total
		resp.PointsTotal = u.ReqPoints
		resp.StorageTotal = u.ReqStorage
		writer.Write(resp)
	}

	reducerFunc := func(pipe <-chan *types.RewardListRes, writer mr.Writer[*types.RewardListRes], cancel func(error)) {
		//获取推荐用户数据分页
		list, err := l.svcCtx.Rpc.FindUserByNameIsDisable(l.ctx, &pb.QueryUserReq{Pid: userData.User.ID, Page: req.Page, Size: req.Size})
		if err != nil {
			cancel(errors.DbError())
			return
		}
		res := <-pipe
		for _, info := range list.Users {
			res.Rewards = append(res.Rewards, &types.Reward{
				Uid:         info.Uid,
				Name:        info.Name,
				Source:      info.Source,
				CreatedTime: info.CreatedTime,
			})
		}
		writer.Write(res)
	}
	reduce, err := mr.MapReduce(generateFunc, mapperFunc, reducerFunc, mr.WithContext(l.ctx))
	if err != nil {
		return nil, err
	}

	return reduce, nil*/
	return nil, err
}

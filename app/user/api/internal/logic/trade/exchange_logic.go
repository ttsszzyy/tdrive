package trade

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeLogic {
	return &ExchangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeLogic) Exchange(req *types.ExchangeReq) (resp *types.Response, err error) {
	//获取用户信息
	/*userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	dict, err := l.svcCtx.Rpc.FindDictByName(l.ctx, &pb.FindDictByNameReq{Code: []string{model.SpaceMerchantExchangeDictCode}})
	if err != nil {
		return nil, err
	}
	if len(dict.DictList) == 0 {
		return nil, errors.CustomError("未配置空间商人兑换参数")
	}
	if req.Storage < 0 {
		return nil, errors.CustomError("存储空间不能小于0")
	}
	value := dict.DictList[0].Value
	storage, _ := strconv.ParseInt(value, 10, 64)
	points := dict.DictList[0].BackupValue

	i := req.Storage / storage * points
	u, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.Rpc.SaveUser(l.ctx, &pb.User{
		Id:              u.Id,
		Points:          i,
		Storage:         -(req.Storage * 1024 * 1024 * 1024),
		ExchangeStorage: req.Storage,
	})
	if err != nil {
		return nil, err
	}*/

	return
}

// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package user

import (
	"context"

	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AccountReq                      = pb.AccountReq
	Admin                           = pb.Admin
	AdminReq                        = pb.AdminReq
	AllAssetIDsRes                  = pb.AllAssetIDsRes
	AssetFile                       = pb.AssetFile
	Assets                          = pb.Assets
	BotCommandItem                  = pb.BotCommandItem
	BotPinMsg                       = pb.BotPinMsg
	CheckUserInviteReq              = pb.CheckUserInviteReq
	CheckUserInviteResp             = pb.CheckUserInviteResp
	ClaimInviteRewardReq            = pb.ClaimInviteRewardReq
	CountAssetsReq                  = pb.CountAssetsReq
	CountAssetsResp                 = pb.CountAssetsResp
	CountUserByExchangeStorageReq   = pb.CountUserByExchangeStorageReq
	CountUserByExchangeStorageResp  = pb.CountUserByExchangeStorageResp
	CountUserInviteReq              = pb.CountUserInviteReq
	CountUserInviteResp             = pb.CountUserInviteResp
	DelAssetFileReq                 = pb.DelAssetFileReq
	DelAssetsReq                    = pb.DelAssetsReq
	DelBotCommandReq                = pb.DelBotCommandReq
	DelBotPinMsgReq                 = pb.DelBotPinMsgReq
	DelMessageReq                   = pb.DelMessageReq
	DelTaskPoolByIdReq              = pb.DelTaskPoolByIdReq
	DelUserReq                      = pb.DelUserReq
	Dict                            = pb.Dict
	DisableUserReq                  = pb.DisableUserReq
	FindAssetFileReq                = pb.FindAssetFileReq
	FindAssetFileResp               = pb.FindAssetFileResp
	FindAssetsReq                   = pb.FindAssetsReq
	FindAssetsResp                  = pb.FindAssetsResp
	FindBotCommandReq               = pb.FindBotCommandReq
	FindBotCommandResp              = pb.FindBotCommandResp
	FindBotPinMsgReq                = pb.FindBotPinMsgReq
	FindBotPinMsgResp               = pb.FindBotPinMsgResp
	FindDictByNameReq               = pb.FindDictByNameReq
	FindDictByNameResp              = pb.FindDictByNameResp
	FindMessageReq                  = pb.FindMessageReq
	FindMessageResp                 = pb.FindMessageResp
	FindOneAssetFileReq             = pb.FindOneAssetFileReq
	FindOneAssetsReq                = pb.FindOneAssetsReq
	FindOneBotCommandReq            = pb.FindOneBotCommandReq
	FindOneShareReq                 = pb.FindOneShareReq
	FindOneTaskReq                  = pb.FindOneTaskReq
	FindOneUserPointsReq            = pb.FindOneUserPointsReq
	FindOneUserPointsResp           = pb.FindOneUserPointsResp
	FindOneUserStorageReq           = pb.FindOneUserStorageReq
	FindOneUserStorageResp          = pb.FindOneUserStorageResp
	FindOneUserTitanTokenReq        = pb.FindOneUserTitanTokenReq
	FindShareReq                    = pb.FindShareReq
	FindShareResp                   = pb.FindShareResp
	FindTaskPoolByTaskTypeReq       = pb.FindTaskPoolByTaskTypeReq
	FindTaskReq                     = pb.FindTaskReq
	FindUserByPidReq                = pb.FindUserByPidReq
	FindUserInviteReq               = pb.FindUserInviteReq
	FindUserInviteResp              = pb.FindUserInviteResp
	FindUserPointsReq               = pb.FindUserPointsReq
	FindUserPointsResp              = pb.FindUserPointsResp
	FindUserRewardIntegralReq       = pb.FindUserRewardIntegralReq
	FindUserRewardIntegralResp      = pb.FindUserRewardIntegralResp
	FindUserStorageExchangeReq      = pb.FindUserStorageExchangeReq
	FindUserStorageExchangeResp     = pb.FindUserStorageExchangeResp
	FindUserStorageExchangeRespList = pb.FindUserStorageExchangeRespList
	GetAllAssetIDsReq               = pb.GetAllAssetIDsReq
	IsOldUserResp                   = pb.IsOldUserResp
	MessageItem                     = pb.MessageItem
	QueryAdminRes                   = pb.QueryAdminRes
	QueryUserReq                    = pb.QueryUserReq
	QueryUserRes                    = pb.QueryUserRes
	ReceiveActionPointsReq          = pb.ReceiveActionPointsReq
	ReceiveUserReq                  = pb.ReceiveUserReq
	Response                        = pb.Response
	SaveAssetFileReq                = pb.SaveAssetFileReq
	SaveAssetFileResp               = pb.SaveAssetFileResp
	SaveAssetsReq                   = pb.SaveAssetsReq
	SaveAssetsResp                  = pb.SaveAssetsResp
	SaveBotCommandReq               = pb.SaveBotCommandReq
	SaveBotPinMsgReq                = pb.SaveBotPinMsgReq
	SaveDictReq                     = pb.SaveDictReq
	SaveMessageReq                  = pb.SaveMessageReq
	SaveShareReq                    = pb.SaveShareReq
	SaveTaskPoolReq                 = pb.SaveTaskPoolReq
	SaveTaskReq                     = pb.SaveTaskReq
	SaveUserPointsAndStorageReq     = pb.SaveUserPointsAndStorageReq
	SaveUserRewardReq               = pb.SaveUserRewardReq
	SaveUserTitanTokenReq           = pb.SaveUserTitanTokenReq
	SaveUserTokenReq                = pb.SaveUserTokenReq
	Share                           = pb.Share
	SignInTaskResp                  = pb.SignInTaskResp
	Task                            = pb.Task
	TaskPool                        = pb.TaskPool
	TaskPoolListResp                = pb.TaskPoolListResp
	TaskResp                        = pb.TaskResp
	UidReq                          = pb.UidReq
	UpdateAssetFileReq              = pb.UpdateAssetFileReq
	UpdateAssetResp                 = pb.UpdateAssetResp
	UpdateAssetsCopyReq             = pb.UpdateAssetsCopyReq
	UpdateAssetsMoveReq             = pb.UpdateAssetsMoveReq
	UpdateAssetsNameReq             = pb.UpdateAssetsNameReq
	UpdateAssetsStatusReq           = pb.UpdateAssetsStatusReq
	UpdateShareReq                  = pb.UpdateShareReq
	UpdateTaskReq                   = pb.UpdateTaskReq
	UpdateUserTitanTokenReq         = pb.UpdateUserTitanTokenReq
	User                            = pb.User
	UserCodeReq                     = pb.UserCodeReq
	UserInfo                        = pb.UserInfo
	UserInvite                      = pb.UserInvite
	UserList                        = pb.UserList
	UserStorageAndToken             = pb.UserStorageAndToken
	UserTitanToken                  = pb.UserTitanToken

	UserZrpcClient interface {
		// 用户端
		FindOneByUid(ctx context.Context, in *UidReq, opts ...grpc.CallOption) (*User, error)
		// 保存用户
		SaveUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*Response, error)
		// 删除用户
		DelUser(ctx context.Context, in *DelUserReq, opts ...grpc.CallOption) (*Response, error)
		// 查询推荐用户
		FindUserByPid(ctx context.Context, in *FindUserByPidReq, opts ...grpc.CallOption) (*UserList, error)
		FindUser(ctx context.Context, in *QueryUserReq, opts ...grpc.CallOption) (*UserList, error)
		// 根据用户邀请码查询用户
		FindUserByCode(ctx context.Context, in *UserCodeReq, opts ...grpc.CallOption) (*User, error)
		// 查询用戶邀請數量
		CountUserInvite(ctx context.Context, in *CountUserInviteReq, opts ...grpc.CallOption) (*CountUserInviteResp, error)
		// 查询奖励空间总数
		CountUserByExchangeStorage(ctx context.Context, in *CountUserByExchangeStorageReq, opts ...grpc.CallOption) (*CountUserByExchangeStorageResp, error)
		// 校验用户是否邀请过
		CheckUserInvite(ctx context.Context, in *CheckUserInviteReq, opts ...grpc.CallOption) (*CheckUserInviteResp, error)
		// 用户积分
		FindOneUserPoints(ctx context.Context, in *FindOneUserPointsReq, opts ...grpc.CallOption) (*FindOneUserPointsResp, error)
		FindUserPoints(ctx context.Context, in *FindUserPointsReq, opts ...grpc.CallOption) (*FindUserPointsResp, error)
		// 用户存储空间
		FindOneUserStorage(ctx context.Context, in *FindOneUserStorageReq, opts ...grpc.CallOption) (*FindOneUserStorageResp, error)
		// 用户空间兑换列表
		FindUserStorageExchange(ctx context.Context, in *FindUserStorageExchangeReq, opts ...grpc.CallOption) (*FindUserStorageExchangeRespList, error)
		// 保存用户积分和空间
		SaveUserPointsAndStorage(ctx context.Context, in *SaveUserPointsAndStorageReq, opts ...grpc.CallOption) (*Response, error)
		// 领取用户奖励空间
		ReceiveUser(ctx context.Context, in *ReceiveUserReq, opts ...grpc.CallOption) (*Response, error)
		// 校验是否满足领取条件
		CheckReceive(ctx context.Context, in *ReceiveActionPointsReq, opts ...grpc.CallOption) (*Response, error)
		// 用户领取活动积分
		ReceiveActionPoints(ctx context.Context, in *ReceiveActionPointsReq, opts ...grpc.CallOption) (*Response, error)
		CheckIsOldUser(ctx context.Context, in *UidReq, opts ...grpc.CallOption) (*IsOldUserResp, error)
		// 用户预订空投
		GetUserStorageAndToken(ctx context.Context, in *UidReq, opts ...grpc.CallOption) (*UserStorageAndToken, error)
		SaveUserToken(ctx context.Context, in *SaveUserTokenReq, opts ...grpc.CallOption) (*Response, error)
		// 发送消息
		SaveMessage(ctx context.Context, in *SaveMessageReq, opts ...grpc.CallOption) (*Response, error)
		// 删除消息
		DelMessage(ctx context.Context, in *DelMessageReq, opts ...grpc.CallOption) (*Response, error)
		// 查询消息
		FindMessage(ctx context.Context, in *FindMessageReq, opts ...grpc.CallOption) (*FindMessageResp, error)
		// 查询用户任务列表
		FindTask(ctx context.Context, in *FindTaskReq, opts ...grpc.CallOption) (*TaskResp, error)
		FindOneTask(ctx context.Context, in *FindOneTaskReq, opts ...grpc.CallOption) (*Task, error)
		// 修改任务完成状态
		UpdateTask(ctx context.Context, in *UpdateTaskReq, opts ...grpc.CallOption) (*Response, error)
		SaveTask(ctx context.Context, in *SaveTaskReq, opts ...grpc.CallOption) (*Response, error)
		// 保存分享资源
		SaveShare(ctx context.Context, in *SaveShareReq, opts ...grpc.CallOption) (*Response, error)
		// 修改分享资源密码和有效期
		UpdateShare(ctx context.Context, in *UpdateShareReq, opts ...grpc.CallOption) (*Response, error)
		// 查看分享资源
		FindShare(ctx context.Context, in *FindShareReq, opts ...grpc.CallOption) (*FindShareResp, error)
		// 查看分享资源
		FindOneShare(ctx context.Context, in *FindOneShareReq, opts ...grpc.CallOption) (*Share, error)
		// 删除分享资源
		DelShare(ctx context.Context, in *FindOneShareReq, opts ...grpc.CallOption) (*Response, error)
		// 保存资源
		SaveAssets(ctx context.Context, in *SaveAssetsReq, opts ...grpc.CallOption) (*SaveAssetsResp, error)
		// 修改资源name,*tag,status,cid,link
		UpdateAssetsName(ctx context.Context, in *UpdateAssetsNameReq, opts ...grpc.CallOption) (*Response, error)
		// 资源移动
		UpdateAssetsMove(ctx context.Context, in *UpdateAssetsMoveReq, opts ...grpc.CallOption) (*Response, error)
		// 资源复制
		UpdateAssetsCopy(ctx context.Context, in *UpdateAssetsCopyReq, opts ...grpc.CallOption) (*Response, error)
		// 资源删除
		DelAssets(ctx context.Context, in *DelAssetsReq, opts ...grpc.CallOption) (*Response, error)
		// 查看资源
		FindAssets(ctx context.Context, in *FindAssetsReq, opts ...grpc.CallOption) (*FindAssetsResp, error)
		FindAssetsNoDel(ctx context.Context, in *FindAssetsReq, opts ...grpc.CallOption) (*FindAssetsResp, error)
		// 查看资源
		FindOneAssets(ctx context.Context, in *FindOneAssetsReq, opts ...grpc.CallOption) (*Assets, error)
		// 还原资源
		RestoreAssets(ctx context.Context, in *DelAssetsReq, opts ...grpc.CallOption) (*Response, error)
		// 清理资源
		ClearAssets(ctx context.Context, in *DelAssetsReq, opts ...grpc.CallOption) (*Response, error)
		// 查看资源数据量
		CountAssets(ctx context.Context, in *CountAssetsReq, opts ...grpc.CallOption) (*CountAssetsResp, error)
		// 获取资源目录下的所有子资源id
		GetAllAssetIds(ctx context.Context, in *GetAllAssetIDsReq, opts ...grpc.CallOption) (*AllAssetIDsRes, error)
		// 查看邀请奖励
		FindUserRewardIntegral(ctx context.Context, in *FindUserRewardIntegralReq, opts ...grpc.CallOption) (*FindUserRewardIntegralResp, error)
		// 保存邀请奖励
		SaveUserReward(ctx context.Context, in *SaveUserRewardReq, opts ...grpc.CallOption) (*Response, error)
		// 领取邀请奖励
		ClaimInviteReward(ctx context.Context, in *ClaimInviteRewardReq, opts ...grpc.CallOption) (*Response, error)
		// 朋友列表
		FindUserInvite(ctx context.Context, in *FindUserInviteReq, opts ...grpc.CallOption) (*FindUserInviteResp, error)
		// 管理端
		FindOneByAccountDeletedTime(ctx context.Context, in *AccountReq, opts ...grpc.CallOption) (*Admin, error)
		FindOneByIdAccountDeletedTime(ctx context.Context, in *AccountReq, opts ...grpc.CallOption) (*Admin, error)
		// 管理端列表
		FindAdminByIdAccountIsDisableDeletedTime(ctx context.Context, in *AdminReq, opts ...grpc.CallOption) (*QueryAdminRes, error)
		// 保存管理端用户
		SaveAdmin(ctx context.Context, in *Admin, opts ...grpc.CallOption) (*Response, error)
		// 保存管理端用户
		DelAdmin(ctx context.Context, in *DelTaskPoolByIdReq, opts ...grpc.CallOption) (*Response, error)
		// 查询用户信息 返回剩余积分、总空间、剩余空间
		FindUserByNameIsDisable(ctx context.Context, in *QueryUserReq, opts ...grpc.CallOption) (*QueryUserRes, error)
		// 禁用用户
		DisableUser(ctx context.Context, in *DisableUserReq, opts ...grpc.CallOption) (*Response, error)
		// 保存任务池
		SaveTaskPool(ctx context.Context, in *SaveTaskPoolReq, opts ...grpc.CallOption) (*Response, error)
		// 查询任务池
		FindTaskPoolByTaskType(ctx context.Context, in *FindTaskPoolByTaskTypeReq, opts ...grpc.CallOption) (*TaskPoolListResp, error)
		// 删除任务池
		DelTaskPoolById(ctx context.Context, in *DelTaskPoolByIdReq, opts ...grpc.CallOption) (*Response, error)
		// 保存数据字典
		SaveDict(ctx context.Context, in *SaveDictReq, opts ...grpc.CallOption) (*Response, error)
		// 查询数据字典 name And code
		FindDictByName(ctx context.Context, in *FindDictByNameReq, opts ...grpc.CallOption) (*FindDictByNameResp, error)
		// 添加pin消息
		SaveBotPinMsg(ctx context.Context, in *SaveBotPinMsgReq, opts ...grpc.CallOption) (*Response, error)
		// 删除pin消息
		DelBotPinMsg(ctx context.Context, in *DelBotPinMsgReq, opts ...grpc.CallOption) (*Response, error)
		// 查询pin消息
		FindBotPinMsg(ctx context.Context, in *FindBotPinMsgReq, opts ...grpc.CallOption) (*FindBotPinMsgResp, error)
		// 保存bot命令
		SaveBotCommand(ctx context.Context, in *SaveBotCommandReq, opts ...grpc.CallOption) (*Response, error)
		// 查询bot命令
		FindBotCommand(ctx context.Context, in *FindBotCommandReq, opts ...grpc.CallOption) (*FindBotCommandResp, error)
		// 删除bot命令
		DelBotCommand(ctx context.Context, in *DelBotCommandReq, opts ...grpc.CallOption) (*Response, error)
		// 查询bot命令
		FindOneBotCommand(ctx context.Context, in *FindOneBotCommandReq, opts ...grpc.CallOption) (*BotCommandItem, error)
		// 保存传输中资源
		SaveAssetFile(ctx context.Context, in *SaveAssetFileReq, opts ...grpc.CallOption) (*SaveAssetFileResp, error)
		// 修改传输中资源tag,status,cid,link,assetId
		UpdateAssetFile(ctx context.Context, in *UpdateAssetFileReq, opts ...grpc.CallOption) (*UpdateAssetResp, error)
		// 资源传输中删除
		DelAssetFile(ctx context.Context, in *DelAssetFileReq, opts ...grpc.CallOption) (*Response, error)
		// 查看传输中资源
		FindAssetFile(ctx context.Context, in *FindAssetFileReq, opts ...grpc.CallOption) (*FindAssetFileResp, error)
		// 查看传输中资源
		FindOneAssetFile(ctx context.Context, in *FindOneAssetFileReq, opts ...grpc.CallOption) (*AssetFile, error)
		// 保存用户titan token
		SaveUserTitanToken(ctx context.Context, in *SaveUserTitanTokenReq, opts ...grpc.CallOption) (*Response, error)
		// 查询用户titan token
		FindOneUserTitanToken(ctx context.Context, in *FindOneUserTitanTokenReq, opts ...grpc.CallOption) (*UserTitanToken, error)
		// 修改用户titan token
		UpdateUserTitanToken(ctx context.Context, in *UpdateUserTitanTokenReq, opts ...grpc.CallOption) (*Response, error)
	}

	defaultUserZrpcClient struct {
		cli zrpc.Client
	}
)

func NewUserZrpcClient(cli zrpc.Client) UserZrpcClient {
	return &defaultUserZrpcClient{
		cli: cli,
	}
}

// 用户端
func (m *defaultUserZrpcClient) FindOneByUid(ctx context.Context, in *UidReq, opts ...grpc.CallOption) (*User, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneByUid(ctx, in, opts...)
}

// 保存用户
func (m *defaultUserZrpcClient) SaveUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveUser(ctx, in, opts...)
}

// 删除用户
func (m *defaultUserZrpcClient) DelUser(ctx context.Context, in *DelUserReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelUser(ctx, in, opts...)
}

// 查询推荐用户
func (m *defaultUserZrpcClient) FindUserByPid(ctx context.Context, in *FindUserByPidReq, opts ...grpc.CallOption) (*UserList, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindUserByPid(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) FindUser(ctx context.Context, in *QueryUserReq, opts ...grpc.CallOption) (*UserList, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindUser(ctx, in, opts...)
}

// 根据用户邀请码查询用户
func (m *defaultUserZrpcClient) FindUserByCode(ctx context.Context, in *UserCodeReq, opts ...grpc.CallOption) (*User, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindUserByCode(ctx, in, opts...)
}

// 查询用戶邀請數量
func (m *defaultUserZrpcClient) CountUserInvite(ctx context.Context, in *CountUserInviteReq, opts ...grpc.CallOption) (*CountUserInviteResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.CountUserInvite(ctx, in, opts...)
}

// 查询奖励空间总数
func (m *defaultUserZrpcClient) CountUserByExchangeStorage(ctx context.Context, in *CountUserByExchangeStorageReq, opts ...grpc.CallOption) (*CountUserByExchangeStorageResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.CountUserByExchangeStorage(ctx, in, opts...)
}

// 校验用户是否邀请过
func (m *defaultUserZrpcClient) CheckUserInvite(ctx context.Context, in *CheckUserInviteReq, opts ...grpc.CallOption) (*CheckUserInviteResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.CheckUserInvite(ctx, in, opts...)
}

// 用户积分
func (m *defaultUserZrpcClient) FindOneUserPoints(ctx context.Context, in *FindOneUserPointsReq, opts ...grpc.CallOption) (*FindOneUserPointsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneUserPoints(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) FindUserPoints(ctx context.Context, in *FindUserPointsReq, opts ...grpc.CallOption) (*FindUserPointsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindUserPoints(ctx, in, opts...)
}

// 用户存储空间
func (m *defaultUserZrpcClient) FindOneUserStorage(ctx context.Context, in *FindOneUserStorageReq, opts ...grpc.CallOption) (*FindOneUserStorageResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneUserStorage(ctx, in, opts...)
}

// 用户空间兑换列表
func (m *defaultUserZrpcClient) FindUserStorageExchange(ctx context.Context, in *FindUserStorageExchangeReq, opts ...grpc.CallOption) (*FindUserStorageExchangeRespList, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindUserStorageExchange(ctx, in, opts...)
}

// 保存用户积分和空间
func (m *defaultUserZrpcClient) SaveUserPointsAndStorage(ctx context.Context, in *SaveUserPointsAndStorageReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveUserPointsAndStorage(ctx, in, opts...)
}

// 领取用户奖励空间
func (m *defaultUserZrpcClient) ReceiveUser(ctx context.Context, in *ReceiveUserReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.ReceiveUser(ctx, in, opts...)
}

// 校验是否满足领取条件
func (m *defaultUserZrpcClient) CheckReceive(ctx context.Context, in *ReceiveActionPointsReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.CheckReceive(ctx, in, opts...)
}

// 用户领取活动积分
func (m *defaultUserZrpcClient) ReceiveActionPoints(ctx context.Context, in *ReceiveActionPointsReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.ReceiveActionPoints(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) CheckIsOldUser(ctx context.Context, in *UidReq, opts ...grpc.CallOption) (*IsOldUserResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.CheckIsOldUser(ctx, in, opts...)
}

// 用户预订空投
func (m *defaultUserZrpcClient) GetUserStorageAndToken(ctx context.Context, in *UidReq, opts ...grpc.CallOption) (*UserStorageAndToken, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.GetUserStorageAndToken(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) SaveUserToken(ctx context.Context, in *SaveUserTokenReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveUserToken(ctx, in, opts...)
}

// 发送消息
func (m *defaultUserZrpcClient) SaveMessage(ctx context.Context, in *SaveMessageReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveMessage(ctx, in, opts...)
}

// 删除消息
func (m *defaultUserZrpcClient) DelMessage(ctx context.Context, in *DelMessageReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelMessage(ctx, in, opts...)
}

// 查询消息
func (m *defaultUserZrpcClient) FindMessage(ctx context.Context, in *FindMessageReq, opts ...grpc.CallOption) (*FindMessageResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindMessage(ctx, in, opts...)
}

// 查询用户任务列表
func (m *defaultUserZrpcClient) FindTask(ctx context.Context, in *FindTaskReq, opts ...grpc.CallOption) (*TaskResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindTask(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) FindOneTask(ctx context.Context, in *FindOneTaskReq, opts ...grpc.CallOption) (*Task, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneTask(ctx, in, opts...)
}

// 修改任务完成状态
func (m *defaultUserZrpcClient) UpdateTask(ctx context.Context, in *UpdateTaskReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateTask(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) SaveTask(ctx context.Context, in *SaveTaskReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveTask(ctx, in, opts...)
}

// 保存分享资源
func (m *defaultUserZrpcClient) SaveShare(ctx context.Context, in *SaveShareReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveShare(ctx, in, opts...)
}

// 修改分享资源密码和有效期
func (m *defaultUserZrpcClient) UpdateShare(ctx context.Context, in *UpdateShareReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateShare(ctx, in, opts...)
}

// 查看分享资源
func (m *defaultUserZrpcClient) FindShare(ctx context.Context, in *FindShareReq, opts ...grpc.CallOption) (*FindShareResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindShare(ctx, in, opts...)
}

// 查看分享资源
func (m *defaultUserZrpcClient) FindOneShare(ctx context.Context, in *FindOneShareReq, opts ...grpc.CallOption) (*Share, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneShare(ctx, in, opts...)
}

// 删除分享资源
func (m *defaultUserZrpcClient) DelShare(ctx context.Context, in *FindOneShareReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelShare(ctx, in, opts...)
}

// 保存资源
func (m *defaultUserZrpcClient) SaveAssets(ctx context.Context, in *SaveAssetsReq, opts ...grpc.CallOption) (*SaveAssetsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveAssets(ctx, in, opts...)
}

// 修改资源name,*tag,status,cid,link
func (m *defaultUserZrpcClient) UpdateAssetsName(ctx context.Context, in *UpdateAssetsNameReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateAssetsName(ctx, in, opts...)
}

// 资源移动
func (m *defaultUserZrpcClient) UpdateAssetsMove(ctx context.Context, in *UpdateAssetsMoveReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateAssetsMove(ctx, in, opts...)
}

// 资源复制
func (m *defaultUserZrpcClient) UpdateAssetsCopy(ctx context.Context, in *UpdateAssetsCopyReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateAssetsCopy(ctx, in, opts...)
}

// 资源删除
func (m *defaultUserZrpcClient) DelAssets(ctx context.Context, in *DelAssetsReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelAssets(ctx, in, opts...)
}

// 查看资源
func (m *defaultUserZrpcClient) FindAssets(ctx context.Context, in *FindAssetsReq, opts ...grpc.CallOption) (*FindAssetsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindAssets(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) FindAssetsNoDel(ctx context.Context, in *FindAssetsReq, opts ...grpc.CallOption) (*FindAssetsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindAssetsNoDel(ctx, in, opts...)
}

// 查看资源
func (m *defaultUserZrpcClient) FindOneAssets(ctx context.Context, in *FindOneAssetsReq, opts ...grpc.CallOption) (*Assets, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneAssets(ctx, in, opts...)
}

// 还原资源
func (m *defaultUserZrpcClient) RestoreAssets(ctx context.Context, in *DelAssetsReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.RestoreAssets(ctx, in, opts...)
}

// 清理资源
func (m *defaultUserZrpcClient) ClearAssets(ctx context.Context, in *DelAssetsReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.ClearAssets(ctx, in, opts...)
}

// 查看资源数据量
func (m *defaultUserZrpcClient) CountAssets(ctx context.Context, in *CountAssetsReq, opts ...grpc.CallOption) (*CountAssetsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.CountAssets(ctx, in, opts...)
}

// 获取资源目录下的所有子资源id
func (m *defaultUserZrpcClient) GetAllAssetIds(ctx context.Context, in *GetAllAssetIDsReq, opts ...grpc.CallOption) (*AllAssetIDsRes, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.GetAllAssetIds(ctx, in, opts...)
}

// 查看邀请奖励
func (m *defaultUserZrpcClient) FindUserRewardIntegral(ctx context.Context, in *FindUserRewardIntegralReq, opts ...grpc.CallOption) (*FindUserRewardIntegralResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindUserRewardIntegral(ctx, in, opts...)
}

// 保存邀请奖励
func (m *defaultUserZrpcClient) SaveUserReward(ctx context.Context, in *SaveUserRewardReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveUserReward(ctx, in, opts...)
}

// 领取邀请奖励
func (m *defaultUserZrpcClient) ClaimInviteReward(ctx context.Context, in *ClaimInviteRewardReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.ClaimInviteReward(ctx, in, opts...)
}

// 朋友列表
func (m *defaultUserZrpcClient) FindUserInvite(ctx context.Context, in *FindUserInviteReq, opts ...grpc.CallOption) (*FindUserInviteResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindUserInvite(ctx, in, opts...)
}

// 管理端
func (m *defaultUserZrpcClient) FindOneByAccountDeletedTime(ctx context.Context, in *AccountReq, opts ...grpc.CallOption) (*Admin, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneByAccountDeletedTime(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) FindOneByIdAccountDeletedTime(ctx context.Context, in *AccountReq, opts ...grpc.CallOption) (*Admin, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneByIdAccountDeletedTime(ctx, in, opts...)
}

// 管理端列表
func (m *defaultUserZrpcClient) FindAdminByIdAccountIsDisableDeletedTime(ctx context.Context, in *AdminReq, opts ...grpc.CallOption) (*QueryAdminRes, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindAdminByIdAccountIsDisableDeletedTime(ctx, in, opts...)
}

// 保存管理端用户
func (m *defaultUserZrpcClient) SaveAdmin(ctx context.Context, in *Admin, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveAdmin(ctx, in, opts...)
}

// 保存管理端用户
func (m *defaultUserZrpcClient) DelAdmin(ctx context.Context, in *DelTaskPoolByIdReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelAdmin(ctx, in, opts...)
}

// 查询用户信息 返回剩余积分、总空间、剩余空间
func (m *defaultUserZrpcClient) FindUserByNameIsDisable(ctx context.Context, in *QueryUserReq, opts ...grpc.CallOption) (*QueryUserRes, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindUserByNameIsDisable(ctx, in, opts...)
}

// 禁用用户
func (m *defaultUserZrpcClient) DisableUser(ctx context.Context, in *DisableUserReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DisableUser(ctx, in, opts...)
}

// 保存任务池
func (m *defaultUserZrpcClient) SaveTaskPool(ctx context.Context, in *SaveTaskPoolReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveTaskPool(ctx, in, opts...)
}

// 查询任务池
func (m *defaultUserZrpcClient) FindTaskPoolByTaskType(ctx context.Context, in *FindTaskPoolByTaskTypeReq, opts ...grpc.CallOption) (*TaskPoolListResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindTaskPoolByTaskType(ctx, in, opts...)
}

// 删除任务池
func (m *defaultUserZrpcClient) DelTaskPoolById(ctx context.Context, in *DelTaskPoolByIdReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelTaskPoolById(ctx, in, opts...)
}

// 保存数据字典
func (m *defaultUserZrpcClient) SaveDict(ctx context.Context, in *SaveDictReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveDict(ctx, in, opts...)
}

// 查询数据字典 name And code
func (m *defaultUserZrpcClient) FindDictByName(ctx context.Context, in *FindDictByNameReq, opts ...grpc.CallOption) (*FindDictByNameResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindDictByName(ctx, in, opts...)
}

// 添加pin消息
func (m *defaultUserZrpcClient) SaveBotPinMsg(ctx context.Context, in *SaveBotPinMsgReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveBotPinMsg(ctx, in, opts...)
}

// 删除pin消息
func (m *defaultUserZrpcClient) DelBotPinMsg(ctx context.Context, in *DelBotPinMsgReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelBotPinMsg(ctx, in, opts...)
}

// 查询pin消息
func (m *defaultUserZrpcClient) FindBotPinMsg(ctx context.Context, in *FindBotPinMsgReq, opts ...grpc.CallOption) (*FindBotPinMsgResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindBotPinMsg(ctx, in, opts...)
}

// 保存bot命令
func (m *defaultUserZrpcClient) SaveBotCommand(ctx context.Context, in *SaveBotCommandReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveBotCommand(ctx, in, opts...)
}

// 查询bot命令
func (m *defaultUserZrpcClient) FindBotCommand(ctx context.Context, in *FindBotCommandReq, opts ...grpc.CallOption) (*FindBotCommandResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindBotCommand(ctx, in, opts...)
}

// 删除bot命令
func (m *defaultUserZrpcClient) DelBotCommand(ctx context.Context, in *DelBotCommandReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelBotCommand(ctx, in, opts...)
}

// 查询bot命令
func (m *defaultUserZrpcClient) FindOneBotCommand(ctx context.Context, in *FindOneBotCommandReq, opts ...grpc.CallOption) (*BotCommandItem, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneBotCommand(ctx, in, opts...)
}

// 保存传输中资源
func (m *defaultUserZrpcClient) SaveAssetFile(ctx context.Context, in *SaveAssetFileReq, opts ...grpc.CallOption) (*SaveAssetFileResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveAssetFile(ctx, in, opts...)
}

// 修改传输中资源tag,status,cid,link,assetId
func (m *defaultUserZrpcClient) UpdateAssetFile(ctx context.Context, in *UpdateAssetFileReq, opts ...grpc.CallOption) (*UpdateAssetResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateAssetFile(ctx, in, opts...)
}

// 资源传输中删除
func (m *defaultUserZrpcClient) DelAssetFile(ctx context.Context, in *DelAssetFileReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelAssetFile(ctx, in, opts...)
}

// 查看传输中资源
func (m *defaultUserZrpcClient) FindAssetFile(ctx context.Context, in *FindAssetFileReq, opts ...grpc.CallOption) (*FindAssetFileResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindAssetFile(ctx, in, opts...)
}

// 查看传输中资源
func (m *defaultUserZrpcClient) FindOneAssetFile(ctx context.Context, in *FindOneAssetFileReq, opts ...grpc.CallOption) (*AssetFile, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneAssetFile(ctx, in, opts...)
}

// 保存用户titan token
func (m *defaultUserZrpcClient) SaveUserTitanToken(ctx context.Context, in *SaveUserTitanTokenReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SaveUserTitanToken(ctx, in, opts...)
}

// 查询用户titan token
func (m *defaultUserZrpcClient) FindOneUserTitanToken(ctx context.Context, in *FindOneUserTitanTokenReq, opts ...grpc.CallOption) (*UserTitanToken, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.FindOneUserTitanToken(ctx, in, opts...)
}

// 修改用户titan token
func (m *defaultUserZrpcClient) UpdateUserTitanToken(ctx context.Context, in *UpdateUserTitanTokenReq, opts ...grpc.CallOption) (*Response, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateUserTitanToken(ctx, in, opts...)
}

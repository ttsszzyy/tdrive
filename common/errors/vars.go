/*
 * @Author: Young
 * @Date: 2022-05-05 14:15:44
 * LastEditors: lihy lihy@zhiannet.com
 * LastEditTime: 2023-10-31 14:57:18
 * @FilePath: /buyday/common/errors/vars.go
 */

package errors

const (
	// ErrUnkonwn 未知错误
	ErrUnkonwn = 999999

	// ErrSystem 系统错误
	ErrSystem = 999998

	// ErrPermission 权限错误
	ErrPermission = 999997

	// ErrUnauth 未授权登录
	ErrUnauth = 999996

	// ErrDB 数据库操作错误
	ErrDB = 999987

	// ErrCodeReqParam 请求参数错误
	ErrCodeReqParam = 999986

	// ErrCodeNotFount 未查询到记录
	ErrCodeNotFount = 999985

	// ErrCodeNotSpace 空间不足
	ErrCodeNotSpace = 999984
	// ErrCodeNotTask 未达标
	ErrCodeNotTask = 999983
	// ErrUserPointsNotEnough 积分不足
	ErrUserPointsNotEnough = 999982
	// ErrUserUploadLimit 用户上传限制
	ErrUserUploadLimit = 999981

	// ErrCodeCustom 自定义业务错误
	ErrCodeCustom = 999900 - iota
	ErrInvalidCallback

	// ErrGetRequestIP 获取请求ip错误
	ErrGetRequestIP
	// ErrNotActionDistance 不在活动范围
	ErrNotActionDistance
	ErrGetRSAPublickey
	ErrReceivedActionPoints
	ErrSameIPForAction
)

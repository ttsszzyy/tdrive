/*
 * Author: lihy lihy@zhiannet.com
 * Date: 2023-12-20 02:55:13
 * LastEditors: lihy lihy@zhiannet.com
 * Note: Need note condition
 */
/*
 * @Author: Young
 * @Date: 2022-08-23 20:24:00
 * LastEditors: lihy lihy@zhiannet.com
 * LastEditTime: 2024-03-02 14:11:09
 * @FilePath: /buyday/app/task/rpc/internal/logic/register.go
 */

package server

import (
	"T-driver/app/user/model"
	"context"

	"github.com/hibiken/asynq"
)

func (t *TaskServer) Register() {
	mux := asynq.NewServeMux()

	// 接受bot信息
	mux.Handle(model.BotProcess, NewTaskBotSendProcessor(t.svcCtx, context.TODO()))
	// Titan回调
	mux.Handle(model.TitanProcess, NewTaskTitanProcess(t.svcCtx, context.TODO()))

	t.mux = mux
}

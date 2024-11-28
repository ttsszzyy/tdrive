package server

import (
	"T-driver/app/tgbot/rpc/internal/svc"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type TaskServer struct {
	svcCtx *svc.ServiceContext
	mux    *asynq.ServeMux
}

func NewTaskServer(svcCtx *svc.ServiceContext) *TaskServer {
	return &TaskServer{
		svcCtx: svcCtx,
	}
}

// Start starts the server.
func (m *TaskServer) Start() {
	m.Register()
	if err := m.svcCtx.AsynqServer.Run(m.mux); err != nil {
		log.Fatal(fmt.Errorf("failed to start mqtask server, error: %v", err))
	}
}

// Stop stops the server.
func (m *TaskServer) Stop() {
	time.Sleep(5 * time.Second)
	m.svcCtx.AsynqServer.Stop()
	m.svcCtx.AsynqServer.Shutdown()
}

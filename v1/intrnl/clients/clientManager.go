package clients

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

type ClientManager struct {
	cfg    *config.Config
	logger logger.Logger
	Grpc   grpc.Manager
}

func NewClientManager(cfg *config.Config, log logger.Logger, grpcHandler grpc.Manager) *ClientManager {
	return &ClientManager{
		cfg:    cfg,
		logger: log,
		Grpc:   grpcHandler,
	}
}

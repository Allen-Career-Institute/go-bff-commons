package middleware

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	m "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

// Manager Middleware manager
type Manager struct {
	cfg     *config.Config
	origins []string
	logger  logger.Logger
	mapper  m.Mapper
	secret  string
	grpc    *grpc.Handler
}

// NewMiddlewareManager Middleware manager constructor
func NewMiddlewareManager(cfg *config.Config, origins []string, log logger.Logger, mapper m.Mapper, secret string, grpc *grpc.Handler) *Manager {
	return &Manager{cfg: cfg, origins: origins, logger: log, mapper: mapper, secret: secret, grpc: grpc}
}

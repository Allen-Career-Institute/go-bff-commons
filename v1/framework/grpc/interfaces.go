package grpc

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/labstack/echo/v4"
	googleGRPC "google.golang.org/grpc"
)

type Manager interface {
	GetConn(ctx echo.Context, log logger.Logger, client string, cnf *config.Config) (googleGRPC.ClientConnInterface, error)
}

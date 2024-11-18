package grpc

import (
	"sync"

	"go.opentelemetry.io/otel/metric"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

type Handler struct {
	Manager
	logger              logger.Logger
	meter               metric.Meter
	safeGrpcConnections ISafeGrpcConnections
}

func NewGRPC(l logger.Logger, meter metric.Meter) *Handler {
	return &Handler{logger: l, meter: meter, safeGrpcConnections: &SafeGrpcConnections{clientConnections: make(map[string]*grpc.ClientConn)}}
}

//nolint:gochecknoglobals // cannot be changed further
var (
	poolMutexes = make(map[string]*sync.Mutex)
	globalMutex sync.Mutex
)

// GetConn returns the grpc connection if already present for a particular client or will initialize a new connection
func (handler *Handler) GetConn(ctx echo.Context, log logger.Logger, client string, cnf *config.Config) (grpc.ClientConnInterface, error) {
	var (
		err  error
		conn *grpc.ClientConn
	)

	// Get or create the pool for the specified client
	conn, exists := handler.safeGrpcConnections.GetConnectionForClient(client)
	if exists {
		log.WithContext(ctx).Infof("connection already exists for client: %s", client)
		return conn, nil
	}

	mutex := getOrCreateMutex(client)
	mutex.Lock()
	defer mutex.Unlock()

	conn, err = createClientConnection(ctx, log, client, cnf, handler)
	if err != nil {
		log.WithContext(ctx).Errorf("error while creating connection for client:%v, err: %v", client, err)
		return nil, err
	}
	log.WithContext(ctx).Infof("successfully created connection for client: %s", client)

	return conn, nil
}

func getOrCreateMutex(client string) *sync.Mutex {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if mutex, ok := poolMutexes[client]; ok {
		return mutex
	}

	poolMutexes[client] = &sync.Mutex{}

	return poolMutexes[client]
}

func createClientConnection(ctx echo.Context, log logger.Logger, client string, cnf *config.Config, handler *Handler) (*grpc.ClientConn, error) {
	conn, err := handler.safeGrpcConnections.CreateConnectionForClient(ctx, log, client, cnf, handler.meter)
	if err != nil {
		return nil, err
	}

	handler.safeGrpcConnections.SetConnectionForClient(client, conn)

	return conn, nil
}

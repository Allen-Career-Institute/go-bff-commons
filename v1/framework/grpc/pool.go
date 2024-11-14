package grpc

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/otel/metrics"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/failsafegrpc"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

const (
	BackoffInitialDelay = 500 * time.Millisecond
	BackoffMaxDelay     = 5 * time.Second
	BackoffMultiplier   = 2
)

type ISafeGrpcConnections interface {
	GetConnectionForClient(client string) (*grpc.ClientConn, bool)
	SetConnectionForClient(client string, conn *grpc.ClientConn)
	CreateConnectionForClient(ctx echo.Context, log logger.Logger, client string, cfg *config.CommonConfig, meter metric.Meter) (*grpc.ClientConn, error)
}

type SafeGrpcConnections struct {
	rwMutex           sync.RWMutex
	clientConnections map[string]*grpc.ClientConn
}

func (safeGrpcConnections *SafeGrpcConnections) GetConnectionForClient(client string) (*grpc.ClientConn, bool) {
	safeGrpcConnections.rwMutex.RLock()
	defer safeGrpcConnections.rwMutex.RUnlock()
	val, exists := safeGrpcConnections.clientConnections[client]

	return val, exists
}

func (safeGrpcConnections *SafeGrpcConnections) SetConnectionForClient(client string, conn *grpc.ClientConn) {
	safeGrpcConnections.rwMutex.Lock()
	defer safeGrpcConnections.rwMutex.Unlock()
	safeGrpcConnections.clientConnections[client] = conn
}

func (*SafeGrpcConnections) CreateConnectionForClient(ctx echo.Context, log logger.Logger, client string, cfg *config.CommonConfig, meter metric.Meter) (*grpc.ClientConn, error) {
	conf := config.GetClientConfigs(client, cfg)

	svcCircuitBreakerConfig := config.GetCircuitBreakerClientConfigs(client, cfg)
	cb := initCircuitBreaker(ctx, log, svcCircuitBreakerConfig, meter, client)

	// note: not adding retry for now as it may cause idempotency issues for some apis
	// svcRetryConfig := config.GetRetryClientConfigs(client, cfg)
	// retry := initRetryer(svcRetryConfig)

	// Create gRPC client interceptor with retry and circuit breaker
	interceptor := failsafegrpc.NewUnaryClientInterceptor[any](cb)

	credential := getCredentials(conf.Endpoint)
	conn, err := grpc.NewClient(conf.Endpoint,
		grpc.WithTransportCredentials(credential),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  BackoffInitialDelay,
				MaxDelay:   BackoffMaxDelay,
				Multiplier: BackoffMultiplier,
			},
		}),
		grpc.WithUnaryInterceptor(interceptor),
	)
	if err != nil {
		log.WithContext(ctx).Errorf("error trying to dial grpc connection, err: %v", err)
		return nil, err
	}

	return conn, nil
}

func initRetryer(svcRetryConfig config.RetryClientConfig) retrypolicy.RetryPolicy[any] {
	// Define the set of gRPC codes to be checked
	grpcCodeSet := map[codes.Code]struct{}{
		codes.Internal:          {},
		codes.Unavailable:       {},
		codes.DeadlineExceeded:  {},
		codes.ResourceExhausted: {},
	}
	retry := retrypolicy.Builder[any]().
		HandleIf(func(_ any, err error) bool {
			if err == nil {
				return false
			}
			grpcCode := status.Code(err)
			// Check if the gRPC code is in the set
			_, found := grpcCodeSet[grpcCode]
			return found
		}).WithMaxRetries(svcRetryConfig.MaxRetries).
		WithDelay(svcRetryConfig.Delay).
		Build()

	return retry
}

func onStateChangeWrapper(currentStateStartTime time.Time, cbMetrics *metrics.CircuitBreakerMetrics,
	ctxx echo.Context, client string, log logger.Logger) func(circuitbreaker.StateChangedEvent) {
	return func(event circuitbreaker.StateChangedEvent) {
		log.WithContext(ctxx).Infof("circuit breaker state changed for client-%s from %s to %s", client, event.OldState, event.NewState)

		clientName := "client-name"
		state := "state"
		paramsMap := map[string]string{
			clientName: client,
			state:      event.NewState.String(),
		}
		circuitBreaker := "CircuitBreaker"
		ctx := ctxx.Request().Context()
		metrics.AddCounter(cbMetrics.StateChangeCount, ctx, circuitBreaker, paramsMap)
		paramsMap[state] = event.OldState.String()
		duration := time.Since(currentStateStartTime)
		metrics.HistogramRecord(cbMetrics.StateDuration, ctx, circuitBreaker, paramsMap, duration)

		currentStateStartTime = time.Now()
	}
}

func initCircuitBreaker(ctx echo.Context, log logger.Logger, svcCircuitBreakerConfig config.CircuitBreakerClientConfig,
	meter metric.Meter, client string) circuitbreaker.CircuitBreaker[any] {
	// Define the set of gRPC codes to be checked
	grpcCodeSet := map[codes.Code]struct{}{
		codes.Internal:          {},
		codes.Unavailable:       {},
		codes.DeadlineExceeded:  {},
		codes.ResourceExhausted: {},
	}

	currentStateStartTime := time.Now()
	circuitBreakerMetric, err := metrics.NewCircuitBreakerMetrics(meter)

	if err != nil {
		log.WithContext(ctx).Errorf("error creating circuit breaker metrics: %v", err)
	}

	cb := circuitbreaker.Builder[any]().
		HandleIf(func(_ any, err error) bool {
			if err == nil {
				return false
			}
			grpcCode := status.Code(err)
			// Check if the gRPC code is in the set
			_, found := grpcCodeSet[grpcCode]
			return found
		}).
		WithFailureRateThreshold(svcCircuitBreakerConfig.FailurePercentageThresholdWithinTimePeriod,
			svcCircuitBreakerConfig.FailureMinExecutionThresholdWithinTimePeriod,
			svcCircuitBreakerConfig.FailurePeriodThreshold).
		WithDelay(svcCircuitBreakerConfig.Delay).
		WithSuccessThreshold(svcCircuitBreakerConfig.SuccessThreshold).
		OnStateChanged(onStateChangeWrapper(currentStateStartTime, circuitBreakerMetric, ctx, client, log)).Build()

	return cb
}

func getCredentials(endpoint string) credentials.TransportCredentials {
	env := os.Getenv("ENV")

	if strings.Contains(endpoint, "localhost:") {
		return insecure.NewCredentials()
	}

	if env != "" {
		return insecure.NewCredentials()
	}
	// TO run on local machine
	return credentials.NewClientTLSFromCert(nil, "")
}

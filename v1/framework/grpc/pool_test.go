package grpc

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/otel/metrics"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func getTestingParams(t *testing.T) (*gomock.Controller, *config.Config, *echo.Echo, logger.Logger,
	metric.Meter, echo.Context, *http.Request) {
	ctrl := gomock.NewController(t)

	e := echo.New()

	c := &config.Config{Logger: config.Logger{Level: "info"}}
	log := logger.NewAPILogger(c)
	log.InitLogger()
	m := otel.Meter("pool-test")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mockContext := echo.New().NewContext(req, rec)
	return ctrl, c, e, log, m, mockContext, req
}

func TestCreateClientConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, cfg, _, log, m, mockContext, _ := getTestingParams(t)

	type args struct {
		ctx    echo.Context
		log    logger.Logger
		client string
		cfg    *config.Config
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "Successful pool creation",
			args: args{
				ctx:    mockContext,
				log:    log,
				client: "test-client",
				cfg:    cfg,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call CreateConnectionForClient function
			safeGrpcPools := &SafeGrpcConnections{}
			_, err := safeGrpcPools.CreateConnectionForClient(tt.args.ctx, tt.args.log, tt.args.client, tt.args.cfg, m)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateConnectionForClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCredentials(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		env      string
		want     credentials.TransportCredentials
	}{
		{
			name:     "Localhost endpoint",
			endpoint: "localhost:50051",
			env:      "",
			want:     insecure.NewCredentials(),
		},
		{
			name:     "Non-localhost endpoint with ENV set",
			endpoint: "example.com:50051",
			env:      "development",
			want:     insecure.NewCredentials(),
		},
		{
			name:     "Non-localhost endpoint without ENV set",
			endpoint: "example.com:50051",
			env:      "",
			want:     credentials.NewClientTLSFromCert(nil, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the environment variable
			os.Setenv("ENV", tt.env)
			defer os.Unsetenv("ENV")

			if got := getCredentials(tt.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initCircuitBreaker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := echo.New().NewContext(nil, nil)
	mockLogger := logger.NewAPILogger(&config.Config{Logger: config.Logger{Level: "info"}})
	mockLogger.InitLogger()
	_, _, _, _, m, _, _ := getTestingParams(t)

	type args struct {
		ctx                     echo.Context
		log                     logger.Logger
		svcCircuitBreakerConfig config.CircuitBreakerClientConfig
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Circuit breaker initialization",
			args: args{
				ctx: mockContext,
				log: mockLogger,
				svcCircuitBreakerConfig: config.CircuitBreakerClientConfig{
					FailurePercentageThresholdWithinTimePeriod:   50,
					FailureMinExecutionThresholdWithinTimePeriod: 10,
					FailurePeriodThreshold:                       60 * time.Second,
					Delay:                                        1 * time.Second,
					SuccessThreshold:                             2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := initCircuitBreaker(tt.args.ctx, tt.args.log, tt.args.svcCircuitBreakerConfig, m, "")

			if cb == nil {
				t.Errorf("initInstrumentedCircuitBreaker() returned nil")
			}
		})
	}
}

func Test_initRetryer(t *testing.T) {
	type args struct {
		svcRetryConfig config.RetryClientConfig
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Retryer initialization with valid config",
			args: args{
				svcRetryConfig: config.RetryClientConfig{
					MaxRetries: 3,
					Delay:      1 * time.Second,
				},
			},
		},
		{
			name: "Retryer initialization with zero retries",
			args: args{
				svcRetryConfig: config.RetryClientConfig{
					MaxRetries: 0,
					Delay:      500 * time.Millisecond,
				},
			},
		},
		{
			name: "Retryer initialization with high delay",
			args: args{
				svcRetryConfig: config.RetryClientConfig{
					MaxRetries: 5,
					Delay:      10 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retry := initRetryer(tt.args.svcRetryConfig)

			if retry == nil {
				t.Errorf("initRetryer() returned nil")
			}
		})
	}
}

func Test_onStateChangeWrapper(t *testing.T) {
	_, _, _, log, m, ctx, _ := getTestingParams(t)
	cbMetric, err := metrics.NewCircuitBreakerMetrics(m)
	if err != nil {
		t.Errorf("NewCircuitBreakerMetrics() returned %v", err)
	}
	type args struct {
		currentStateStartTime time.Time
		cbMetrics             *metrics.CircuitBreakerMetrics
		context               echo.Context
		client                string
		log                   logger.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				currentStateStartTime: time.Now(),
				cbMetrics:             cbMetric,
				context:               ctx,
				log:                   log,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			onStateChangeMethod := onStateChangeWrapper(tt.args.currentStateStartTime, tt.args.cbMetrics, tt.args.context, tt.args.client, tt.args.log)
			onStateChangeMethod(circuitbreaker.StateChangedEvent{})
		})
	}
}

func TestSafeGrpcPools_GetConnectionPoolForClient(t *testing.T) {
	client1 := "client1"
	client2 := "client2"
	connectionPools := map[string]*grpc.ClientConn{
		"client1": {},
	}
	tests := []struct {
		name   string
		client string
	}{
		{
			name:   "Pool exists for client",
			client: client1,
		},
		{
			name:   "Pool does not exist for client",
			client: client2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpcPools := &SafeGrpcConnections{
				clientConnections: connectionPools,
			}
			pool, exists := grpcPools.GetConnectionForClient(tt.client)
			if tt.client == client1 {
				assert.True(t, exists)
				assert.NotNil(t, pool)
			} else {
				assert.False(t, exists)
				assert.Nil(t, pool)
			}
		})
	}
}

func TestSafeGrpcPools_SetConnectionPoolForClient(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpcPools := &SafeGrpcConnections{
				clientConnections: map[string]*grpc.ClientConn{},
			}
			grpcPools.SetConnectionForClient("test-client", &grpc.ClientConn{})
		})
	}
}

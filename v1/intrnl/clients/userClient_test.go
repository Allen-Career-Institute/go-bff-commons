package clients

import (
	"context"
	"errors"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/golang/mock/gomock"
	grpc2 "google.golang.org/grpc"
)

func getTestingParams(t *testing.T) (*gomock.Controller, *config.Config, *echo.Echo, logger.Logger) {
	ctrl := gomock.NewController(t)

	e := echo.New()

	c := &config.Config{Logger: config.Logger{Level: "info"}}
	log := logger.NewAPILogger(c)
	log.InitLogger()

	return ctrl, c, e, log
}

func getGrpcTestConn(ctx echo.Context, log logger.Logger) (*grpc2.ClientConn, error) {
	conn, err := grpc2.DialContext(context.Background(), "test", grpc2.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc2.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc2.WithConnectParams(grpc2.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  500 * time.Millisecond,
				MaxDelay:   5 * time.Second,
				Multiplier: 2,
			},
		}))
	if err != nil {
		log.WithContext(ctx).Errorf("error trying to dial grpc in grpc-pool, err: %v", err)
		return nil, err
	}

	return conn, nil
}

func Test_GetUser(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set(utils.UserID, "123")

	conn, err := getGrpcTestConn(ctx, log)
	if err != nil {
		return
	}

	testCases := []struct {
		description   string
		mock          []*gomock.Call
		expectedError error
	}{
		{
			description: "Failure Case: failed grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(
					nil,
					errors.New("nil conn error"),
				),
			},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(
			tc.description, func(t *testing.T) {
				_, err := cm.GetUser(ctx, c, cm.Grpc)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

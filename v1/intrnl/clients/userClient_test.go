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

	user "github.com/Allen-Career-Institute/common-protos/user_management/v1"
	"github.com/Allen-Career-Institute/common-protos/user_management/v1/types"
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

func Test_GetUserByID(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodPatch, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetUserByID(ctx, c, "432", "123")
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetUserServiceClient(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetUserServiceClient(ctx)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetAdminUserServiceClient(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetAdminUserServiceClient(ctx)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetAddressClient(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	conn, err := getGrpcTestConn(ctx, log)
	if err != nil {
		return
	}

	testCases := []struct {
		description    string
		mock           []*gomock.Call
		expectedOutput interface{}
	}{
		{
			description: "Success Case: failed grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(conn, nil),
			},
			expectedOutput: user.NewAddressClient(conn),
		},
	}
	for i, tc := range testCases {
		t.Run(
			tc.description, func(t *testing.T) {
				out := cm.GetAddressClient(ctx)
				if out != nil && !reflect.DeepEqual(tc.expectedOutput, out) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedOutput, out)
				}
			},
		)
	}
}

func Test_GetAddressByID(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
	}
	for i, tc := range testCases {
		t.Run(
			tc.description, func(t *testing.T) {
				_, err := cm.GetAddressByID(ctx, c, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetAllAddresses(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
	}
	for i, tc := range testCases {
		t.Run(
			tc.description, func(t *testing.T) {
				_, err := cm.GetAllAddresses(ctx, c, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetLocationClient(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	conn, err := getGrpcTestConn(ctx, log)
	if err != nil {
		return
	}

	testCases := []struct {
		description    string
		mock           []*gomock.Call
		expectedOutput interface{}
	}{
		{
			description: "Failure Case: failed grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(conn, nil),
			},
			expectedOutput: user.NewLocationClient(conn),
		},
	}
	for i, tc := range testCases {
		t.Run(
			tc.description, func(t *testing.T) {
				out := cm.GetLocationClient(ctx)
				if out != nil && !reflect.DeepEqual(tc.expectedOutput, out) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedOutput, out)
				}
			},
		)
	}
}

func TestClientManager_GetUserIDToUserMap(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, gomock.Any()).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, gomock.Any()).Return(
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
				_, err := cm.GetUserIDToUserMap(ctx, []string{"test"})
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetUsers(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetUsers(ctx, c, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_UpdateUserPersonaType(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				err := cm.UpdateUserPersonaType(ctx, nil, types.PersonaType_STUDENT, types.PersonaType_PARENT)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetUserIDsByEmpIDs(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetUserIDsByEmpIDs(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetUserMinimal(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetUserMinimal(ctx, c, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_RegisterUser(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.RegisterUser(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_EditUserProfile(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.EditUserProfile(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetUserIDFromEmpID(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetUserIDFromEmpID(ctx, "123")
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_CreateUser(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.CreateUser(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetUserIdentitiesByIdentity(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetUserIdentitiesByIdentity(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetIdentitiesByUser(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetIdentitiesByUser(ctx, c, "123", "456")
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func Test_GetUserByIdentity(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.GetUserByIdentity(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func TestClientManager_CheckCredentials(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.CheckCredentials(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func TestClientManager_DeleteCredentials(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.DeleteCredentials(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func TestClientManager_BulkCardStatusUpdate(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)
	mockClientConn := NewMockClientConnInterface(ctrl)
	mockStream := NewMockClientStream(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	testCases := []struct {
		description   string
		mock          []*gomock.Call
		expectedError error
	}{
		{
			description: "Success case",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(mockClientConn, nil),
				mockClientConn.EXPECT().NewStream(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockStream, nil),
				mockStream.EXPECT().SendMsg(gomock.Any()).Return(nil),
				mockStream.EXPECT().CloseSend().Return(nil),
				mockStream.EXPECT().RecvMsg(gomock.Any()).Return(nil),
			},
			expectedError: errors.New("nil conn error"),
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
		{
			description: "Failure Case: ClientConn NewStream error",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(mockClientConn, nil),
				mockClientConn.EXPECT().NewStream(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New("NewStream error"),
				),
			},
			expectedError: errors.New("NewStream error"),
		},
		{
			description: "Failure Case: Stream SendMsg error",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(mockClientConn, nil),
				mockClientConn.EXPECT().NewStream(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockStream, nil),
				mockStream.EXPECT().SendMsg(gomock.Any()).Return(errors.New("SendMsg error")),
			},
			expectedError: errors.New("SendMsg error"),
		},
		{
			description: "Failure Case: Stream CloseSend error",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(mockClientConn, nil),
				mockClientConn.EXPECT().NewStream(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockStream, nil),
				mockStream.EXPECT().SendMsg(gomock.Any()).Return(nil),
				mockStream.EXPECT().CloseSend().Return(errors.New("CloseSend error")),
			},
			expectedError: errors.New("CloseSend error"),
		},
		{
			description: "Failure Case: Stream RecvMsg error",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.UserServiceClient, c).Return(mockClientConn, nil),
				mockClientConn.EXPECT().NewStream(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockStream, nil),
				mockStream.EXPECT().SendMsg(gomock.Any()).Return(nil),
				mockStream.EXPECT().CloseSend().Return(nil),
				mockStream.EXPECT().RecvMsg(gomock.Any()).Return(errors.New("RecvMsg error")),
			},
			expectedError: errors.New("RecvMsg error"),
		},
	}
	for i, tc := range testCases {
		t.Run(
			tc.description, func(t *testing.T) {
				_, err := cm.BulkCardStatusUpdate(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func TestClientManager_UpdateUser(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodPut, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.UpdateUser(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func TestClientManager_EvaluateStudentEligibility(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.EvaluateStudentEligibility(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func TestClientManager_SubmitDocuments(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.SubmitDocuments(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

func TestClientManager_UpdateVerificationResourceByUserID(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

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
				_, err := cm.UpdateVerificationResourceByUserID(ctx, nil)
				if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
					t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
				}
			},
		)
	}
}

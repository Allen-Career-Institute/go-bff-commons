package clients

import (
	"errors"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	dc "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Allen-Career-Institute/common-protos/cal/v1/request"
	"github.com/Allen-Career-Institute/common-protos/cal/v1/response"
	calTypes "github.com/Allen-Career-Institute/common-protos/learning_material/v1/types/enums"
	"github.com/golang/mock/gomock"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
)

func Test_GetPlaylistFileRequest(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)
	cm := NewClientManager(c, log, grpcMock)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	testCases := []struct {
		description    string
		mock           []*gomock.Call
		deviceType     string
		expectedOutput *request.GetVideoPlaylistFileRequest
	}{
		{
			description: "Success Case",
			expectedOutput: &request.GetVideoPlaylistFileRequest{
				Id:    "123",
				Codec: calTypes.Codec_x264,
			},
		},
		{
			description: "Success Case",
			deviceType:  utils.DeviceTypeWeb,
			expectedOutput: &request.GetVideoPlaylistFileRequest{
				Id:    "123",
				Codec: calTypes.Codec_x264,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx.Set("userID", "123")
			ctx.Request().Header.Set(utils.DeviceType, tc.deviceType)

			out := cm.GetPlaylistFileRequest(ctx, c, "123")

			if out != nil && !reflect.DeepEqual(tc.expectedOutput, out) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedOutput, out)
			}
		})
	}
}

func Test_GetLearningMaterial(t *testing.T) {
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
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(nil, errors.New("nil conn error")),
			},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx.Set("userID", "123")

			_, err := cm.GetLearningMaterial(ctx, c, nil)

			if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
			}
		})
	}
}

func Test_FetchLearningContent(t *testing.T) {
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
		expectedError  error
		contentIDs     []string
		expectedOutput []*response.GetLearningMaterialResponse
	}{
		{
			description:    "Failure Case: empty contentIds",
			expectedOutput: []*response.GetLearningMaterialResponse{},
		},
		{
			description: "Failure Case: failed grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(conn, nil),
			},
			contentIDs:    []string{"123"},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(nil, errors.New("nil conn error")),
			},
			contentIDs:    []string{"123"},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx.Set("userID", "123")

			out, err := cm.FetchLearningContent(ctx, c, tc.contentIDs)

			if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
			}

			if out != nil && reflect.DeepEqual(tc.expectedOutput, out) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedOutput, out)
			}
		})
	}
}

func Test_GetNAC(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	cm := NewClientManager(c, log, grpcMock)
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
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(nil, errors.New("nil conn error")),
			},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			_, err := cm.GetNAC(ctx, c, nil)
			if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
			}

		})
	}
}

func Test_CreateFlashcardsSession(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	cm := NewClientManager(c, log, grpcMock)
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
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(nil, errors.New("nil conn error")),
			},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			_, err := cm.CreateFlashcardsSession(ctx, c, nil)
			if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
			}

		})
	}
}

func Test_GetFlashcards(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	cm := NewClientManager(c, log, grpcMock)
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
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(nil, errors.New("nil conn error")),
			},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			_, err := cm.GetFlashcards(ctx, c, nil)
			if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
			}

		})
	}
}

func Test_FlashcardsCount(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	cm := NewClientManager(c, log, grpcMock)
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
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(nil, errors.New("nil conn error")),
			},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			_, err := cm.GetFlashcardsCount(ctx, c, nil)
			if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
			}

		})
	}
}

func Test_GetFlashcardsStats(t *testing.T) {
	ctrl, c, e, log := getTestingParams(t)

	grpcMock := grpc.NewMockManager(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	cm := NewClientManager(c, log, grpcMock)
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
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.CalServiceClient, c).Return(nil, errors.New("nil conn error")),
			},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			_, err := cm.GetFlashcardsSessionStats(ctx, c, nil)
			if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
			}

		})
	}
}

func TestClientManager_GetPlayList(t *testing.T) {
	_, c, e, log := getTestingParams(t)
	grpcMock := &grpc.MockManager1{}
	dcMock := &dc.MockDynamicConfig{}
	c.DynamicConfig = dcMock
	cm := NewClientManager(c, log, grpcMock)
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	conn, err := getGrpcTestConn(ctx, log)
	if err != nil {
		return
	}

	playlistSourceDC := []string{SourceTencent, SourceAllen}
	var playlistSourceInterface interface{} = playlistSourceDC

	type args struct {
		c   echo.Context
		cnf *config.Config
		req *request.GetVideoPlaylistFileRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *response.GetVideoPlaylistFileResponse
		wantErr error
		mock    func()
	}{
		{
			name: "Failure Case: failed grpc Connection",
			args: args{
				c:   ctx,
				cnf: c,
				req: &request.GetVideoPlaylistFileRequest{},
			},
			want:    nil,
			wantErr: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
			mock: func() {
				dcMock.On("Get", mock.Anything).Return("500", nil).Times(3)
				grpcMock.On("GetConn", mock.Anything, log, clients.CalServiceClient, c).Return(conn, nil).Once()
				dcMock.On("GetAsInterface", mock.Anything).Return(playlistSourceInterface, nil).Once()
			},
		},
		{
			name: "Failure Case: nil grpc Connection",
			args: args{
				c:   ctx,
				cnf: c,
				req: &request.GetVideoPlaylistFileRequest{},
			},
			want:    nil,
			wantErr: errors.New("nil conn error"),
			mock: func() {
				dcMock.On("Get", mock.Anything).Return("500", nil).Times(3)
				grpcMock.On("GetConn", mock.Anything, log, clients.CalServiceClient, c).Return(nil, errors.New("nil conn error")).Once()
			},
		},
		{
			name: "Failure Case: GetAsInterface error",
			args: args{
				c:   ctx,
				cnf: c,
				req: &request.GetVideoPlaylistFileRequest{},
			},
			want:    nil,
			wantErr: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
			mock: func() {
				dcMock.On("Get", mock.Anything).Return("500", nil).Times(3)
				grpcMock.On("GetConn", mock.Anything, log, clients.CalServiceClient, c).Return(conn, nil).Once()
				dcMock.On("GetAsInterface", mock.Anything).Return(nil, errors.New("nil dc response")).Once()
			},
		},
		{
			name: "Failure Case: GetAsInterface wrong data type",
			args: args{
				c:   ctx,
				cnf: c,
				req: &request.GetVideoPlaylistFileRequest{},
			},
			want:    nil,
			wantErr: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
			mock: func() {
				dcMock.On("Get", mock.Anything).Return("500", nil).Times(3)
				grpcMock.On("GetConn", mock.Anything, log, clients.CalServiceClient, c).Return(conn, nil).Once()
				dcMock.On("GetAsInterface", mock.Anything).Return(SourceAllen, nil).Once()
			},
		},
	}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := cm.GetPlayList(tt.args.c, tt.args.cnf, tt.args.req)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		})
	}
}

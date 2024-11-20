package clients

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Allen-Career-Institute/common-protos/cal/v1/request"
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

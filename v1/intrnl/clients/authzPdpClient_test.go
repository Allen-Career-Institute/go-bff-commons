package clients

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Allen-Career-Institute/common-protos/authorization/v1/types"
	"github.com/golang/mock/gomock"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
)

func Test_EnforceRbac(t *testing.T) {
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
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.AuthorizationPdpServiceClient, c).Return(conn, nil),
			},
			expectedError: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp: address test: missing port in address\""),
		},
		{
			description: "Failure Case: nil grpc Connection",
			mock: []*gomock.Call{
				grpcMock.EXPECT().GetConn(gomock.Any(), log, clients.AuthorizationPdpServiceClient, c).Return(nil, errors.New("nil conn error")),
			},
			expectedError: errors.New("nil conn error"),
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx.Set("userID", "123")

			_, err := cm.EnforceRbac(ctx, c, types.ResourceTypes_RESOURCE_UNSPECIFIED, types.Action_ACTION_UNSPECIFIED)

			if err != nil && !reflect.DeepEqual(tc.expectedError.Error(), err.Error()) {
				t.Errorf("Testcase: %v FAILED, Expected: %v, Got: %v", i+1, tc.expectedError, err.Error())
			}
		})
	}
}

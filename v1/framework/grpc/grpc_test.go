package grpc

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"testing"
)

func TestHandler_GetConn(t *testing.T) {
	ctrl, c, _, log, _, ctx, _ := getTestingParams(t)
	safeGrpcPoolMock := NewMockISafeGrpcPool(ctrl)
	client1 := "client-present"
	client2 := "client-not-present"
	tests := []struct {
		name   string
		client string
		mock   func()
	}{
		{
			name:   "Grpc connection exists for client",
			client: client1,
			mock: func() {
				safeGrpcPoolMock.EXPECT().GetConnectionForClient(client1).Return(&grpc.ClientConn{}, true)
			},
		},
		{
			name:   "Grpc connection does not exist for client",
			client: client2,
			mock: func() {
				safeGrpcPoolMock.EXPECT().GetConnectionForClient(client2).Return(nil, false)
				safeGrpcPoolMock.EXPECT().CreateConnectionForClient(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&grpc.ClientConn{}, nil)
				safeGrpcPoolMock.EXPECT().SetConnectionForClient(client2, &grpc.ClientConn{}).Return()
			},
		},
		{
			name:   "Error while creating connection for client",
			client: client2,
			mock: func() {
				safeGrpcPoolMock.EXPECT().GetConnectionForClient(client2).Return(nil, false)
				safeGrpcPoolMock.EXPECT().CreateConnectionForClient(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error in creating client connection"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{
				safeGrpcConnections: safeGrpcPoolMock,
			}
			tt.mock()
			conn, err := handler.GetConn(ctx, log, tt.client, c)
			if tt.client == client1 {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
			} else if tt.client == client2 {
				if conn == nil {
					assert.Error(t, err)
				} else {
					assert.NotNil(t, conn)
				}
			}
		})
	}
}

func TestNewGRPC(t *testing.T) {
	type args struct {
		l     logger.Logger
		meter metric.Meter
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				l:     nil,
				meter: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGRPC(tt.args.l, tt.args.meter)
			assert.NotNil(t, got)
		})
	}
}

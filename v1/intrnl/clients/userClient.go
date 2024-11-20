package clients

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	internal "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"

	googlegrpc "google.golang.org/grpc"

	"github.com/labstack/echo/v4"

	"github.com/Allen-Career-Institute/common-protos/learning_material/v1/types/errors"
	user "github.com/Allen-Career-Institute/common-protos/user_management/v1"
	userReq "github.com/Allen-Career-Institute/common-protos/user_management/v1/request"
	userRes "github.com/Allen-Career-Institute/common-protos/user_management/v1/response"
)

func (cm *ClientManager) getUserServiceGrpcConn(c echo.Context) (googlegrpc.ClientConnInterface, error) {
	conn, err := cm.Grpc.GetConn(c, cm.logger, clients.UserServiceClient, cm.cfg)
	if err != nil {
		cm.logger.WithContext(c).Errorf("failed to get user service client conn, err: %v", err)
		return nil, err
	}

	return conn, nil
}

func (cm *ClientManager) GetUserServiceClient(c echo.Context) (user.UserClient, error) {
	conn, err := cm.getUserServiceGrpcConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error("Error while creating grpc connection with user service: %v", err.Error())
		return nil, err
	}

	userClient := user.NewUserClient(conn)

	return userClient, nil
}

func (cm *ClientManager) GetUser(c echo.Context, cnf *config.Config, grpcHandler grpc.Manager) (*userRes.GetUserResponse, error) {
	uid, err := internal.GetUserID(c)
	if err != nil {
		return nil, err
	}

	if uid == "" {
		cm.logger.WithContext(c).Error("GetUser - UserID is missing")
		return nil, errors.ErrorBadRequest("UserID is missing")
	}

	tenant, err := internal.GetTenantID(c)
	if err != nil {
		return nil, err
	}

	cm.logger.WithContext(c).Infof("calling user service with user %s  and tenant id %s", uid, tenant)
	request := &userReq.GetUserRequest{
		TenantId: tenant,
		UserId:   uid,
	}

	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := userClient.GetUser(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("error while calling user service for user: %v and tenant: %v", uid, tenant)
		return nil, err
	}

	return response, nil
}

package clients

import (
	"fmt"

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
	userTypes "github.com/Allen-Career-Institute/common-protos/user_management/v1/types"
)

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

func (cm *ClientManager) GetUserByID(c echo.Context, _ *config.Config, uid, tenant string) (*userRes.GetUserResponse, error) {
	if uid == "" {
		cm.logger.WithContext(c).Error("GetUserByID - UserID is missing")
		return nil, errors.ErrorBadRequest("UserID is missing")
	}

	cm.logger.WithContext(c).Infof("calling user service with user %s and tenant id %s", uid, tenant)

	request := &userReq.GetUserRequest{
		TenantId: tenant,
		UserId:   uid,
	}

	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := userClient.GetUser(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while Calling user service for user: %v and tenant: %v", uid, tenant)
		return nil, err
	}

	return response, nil
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

func (cm *ClientManager) GetAdminUserServiceClient(c echo.Context) (user.UserAdminClient, error) {
	conn, err := cm.getUserServiceGrpcConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error("Error while creating grpc connection with user service: %v", err.Error())
		return nil, err
	}

	userAdminClient := user.NewUserAdminClient(conn)

	return userAdminClient, nil
}

func (cm *ClientManager) getUserServiceGrpcConn(c echo.Context) (googlegrpc.ClientConnInterface, error) {
	conn, err := cm.Grpc.GetConn(c, cm.logger, clients.UserServiceClient, cm.cfg)
	if err != nil {
		cm.logger.WithContext(c).Errorf("failed to get user service client conn, err: %v", err)
		return nil, err
	}

	return conn, nil
}

func (cm *ClientManager) GetAddressClient(c echo.Context) user.AddressClient {
	conn, err := cm.getUserServiceGrpcConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error("Error while creating grpc connection with user service: %v", err)
		return nil
	}

	addressClient := user.NewAddressClient(conn)

	return addressClient
}

func (cm *ClientManager) GetAddressByID(c echo.Context, cnf *config.Config, request *userReq.GetAddressByIdRequest) (*userRes.GetAddressByIdResponse, error) {
	addressClient := cm.GetAddressClient(c)

	conf := config.GetClientConfigs(clients.UserServiceClient, cnf)

	ctx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout*4)
	defer apiCancel()

	sores, err := addressClient.GetAddressById(ctx, request)

	return sores, err
}

func (cm *ClientManager) GetAllAddresses(c echo.Context, cnf *config.Config, request *userReq.GetAllAddressesRequest) (*userRes.GetAllAddressesResponse, error) {
	addressClient := cm.GetAddressClient(c)

	conf := config.GetClientConfigs(clients.UserServiceClient, cnf)

	ctx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout*4)
	defer apiCancel()

	resp, err := addressClient.GetAllAddresses(ctx, request)

	return resp, err
}

func (cm *ClientManager) GetLocationClient(c echo.Context) user.LocationClient {
	conn, err := cm.getUserServiceGrpcConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error("Error while creating grpc connection with user service: %v", err)
		return nil
	}

	locationClient := user.NewLocationClient(conn)

	return locationClient
}

func (cm *ClientManager) GetUserIDToUserMap(c echo.Context, userIDs []string) (map[string]*userTypes.UserInfo, error) {
	tenantID, err2 := internal.GetTenantID(c)
	if err2 != nil {
		return nil, err2
	}

	cm.logger.WithContext(c).Infof("calling user service for get userIDToUserMap, userIDs :: %v", userIDs)

	request := &userReq.GetUsersRequest{
		TenantId: tenantID,
		UserIds:  userIDs,
	}

	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := userClient.GetUsers(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while calling user service for get users, request ::%v, Err: %v", request, err)
		return nil, err
	}

	userIDToUserMap := make(map[string]*userTypes.UserInfo)
	for _, info := range response.GetData() {
		userIDToUserMap[info.GetId()] = info
	}

	return userIDToUserMap, nil
}

func (cm *ClientManager) GetUsers(c echo.Context, _ *config.Config, request *userReq.GetUsersRequest) (*userRes.GetUsersResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service for get users, request ::%v", request)

	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := userClient.GetUsers(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while calling user service for get users, request ::%v, Err: %v", request, err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Response for Get Users %s ", response.String())

	return response, nil
}

func (cm *ClientManager) UpdateUserPersonaType(c echo.Context, userIDs []string, currentPersonaType, newPersonaType userTypes.PersonaType) error {
	cm.logger.WithContext(c).Infof("calling user service for update user persona type, request ::%v", userIDs)

	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	tenantID, err := internal.GetTenantID(c)
	if err != nil {
		return err
	}

	request := &userReq.UpdatePersonaTypeRequest{
		UserIds:            userIDs,
		TenantId:           tenantID,
		CurrentPersonaType: currentPersonaType,
		NewPersonaType:     newPersonaType,
	}

	_, err = userClient.UpdatePersonaType(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Infof("Error while calling user service for update user persona type, request ::%v", request)
		return err
	}

	return nil
}

func (cm *ClientManager) GetUserIDsByEmpIDs(c echo.Context, empIDs []string) ([]*userTypes.UserIDEmpIDMap, error) {
	cm.logger.WithContext(c).Infof("calling user service for get user ids by emp ids, request ::%v", empIDs)

	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	tenantID, err := internal.GetTenantID(c)
	if err != nil {
		return nil, err
	}

	request := &userReq.GetUserIDsByEmployeeIDsRequest{
		EmployeeIds: empIDs,
		TenantId:    tenantID,
	}

	response, err := userClient.GetUserIDsByEmployeeIDs(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while calling user service for get user ids by emp ids, request ::%v", request)
		return nil, err
	}

	return response.GetData(), nil
}

func (cm *ClientManager) GetUserMinimal(c echo.Context, _ *config.Config, request *userReq.GetUserMinimalRequest) (*userRes.GetUserMinimalResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service for get minimal  users, request ::%v", request)

	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := userClient.GetUserMinimal(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Infof("Error while calling user service for get users, request ::%v", request)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) RegisterUser(c echo.Context, req *userReq.RegisterUserRequest) (*userRes.RegisterUserResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service for register user, request ::%v", req)

	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := userClient.RegisterUser(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while calling user service for register user, request ::%v", req)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) EditUserProfile(c echo.Context, req *userReq.UpdateUserRequest) (*userRes.UpdateUserResponse, error) {
	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	updatedUser, errUpdate := userClient.UpdateUser(apiCtx, req)
	if errUpdate != nil {
		return nil, errUpdate
	}

	return updatedUser, nil
}

func (cm *ClientManager) GetUserIDFromEmpID(c echo.Context, employeeID string) (string, error) {
	if employeeID == "" {
		cm.logger.WithContext(c).Errorf("GetUserIDFromEmpID - EmployeeID is missing")
		return "", errors.ErrorBadRequest("EmployeeID is missing")
	}

	tenantID, err := internal.GetTenantID(c)
	if err != nil {
		return "", err
	}

	getUserRequest := &userReq.GetUsersRequest{
		PageSize:    10,
		CurrentPage: 0,
		EmployeeIds: []string{employeeID},
		TenantId:    tenantID,
	}

	usersResponse, err := cm.GetUsers(c, nil, getUserRequest)
	if err != nil {
		return "", err
	}

	userInfos := usersResponse.GetData()
	if len(userInfos) == 0 {
		cm.logger.WithContext(c).Errorf("employee id not found %v", employeeID)
		return "", errors.ErrorBadRequest("employee id not found")
	}

	cm.logger.WithContext(c).Infof("user info found for employee id %v", userInfos)

	return userInfos[0].Id, nil
}

func (cm *ClientManager) CreateUser(c echo.Context, req *userReq.CreateUserRequest) (*userRes.CreateUserResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service for create user, request :: %v", req)
	userClient, err := cm.GetUserServiceClient(c)

	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client : %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	response, err := userClient.CreateUser(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while calling user service for create user, with error: %s and request :: %v", err.Error(), req)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetUserIdentitiesByIdentity(c echo.Context, req *userReq.GetUserIdentitiesByIdentityRequest) (*userRes.GetUserIdentitiesByIdentityResponse, error) {
	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	userIdentity, err := userClient.GetUserIdentitiesByIdentity(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while GetUserIdentitiesByIdentity: %v", err)
		return nil, err
	}

	return userIdentity, nil
}

func (cm *ClientManager) GetIdentitiesByUser(c echo.Context, cnf *config.Config, tenantID, userID string) (*userRes.GetIdentitiesByUserResponse, error) {
	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cnf)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	userIdentities, err := userClient.GetIdentitiesByUser(apiCtx, &userReq.GetIdentitiesByUserRequest{TenantId: tenantID, UserId: userID})
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while GetIdentitiesByUser: %v", err.Error())
		return nil, err
	}

	return userIdentities, nil
}

func (cm *ClientManager) GetUserByIdentity(c echo.Context, req *userReq.GetUserByIdentityRequest) (*userRes.GetUserByIdentityResponse, error) {
	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	userOutput, err := userClient.GetUserByIdentity(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while GetIdentitiesByUser: %v", err.Error())
		return nil, err
	}

	return userOutput, nil
}

func (cm *ClientManager) AdminUpdateUser(c echo.Context, adminUpdateUserReq *userReq.AdminUpdateUserRequest) (*userRes.AdminUpdateUserResponse, error) {
	adminUserClient, err := cm.GetAdminUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service admin client: %v", err)

		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	resp, err := adminUserClient.UpdateUser(apiCtx, adminUpdateUserReq)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while AdminUpdateUser: %v", err)

		return nil, err
	}

	return resp, err
}

func (cm *ClientManager) AdminDeleteUser(c echo.Context, request *userReq.DeleteUserRequest) (*userRes.DeleteUserResponse, error) {
	userAdminClient, err := cm.GetAdminUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service admin client: %v", err)

		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	resp, errDelete := userAdminClient.DeleteUser(apiCtx, request)
	if errDelete != nil {
		cm.logger.WithContext(c).Errorf("Error while AdminDeleteUser: %v", err)

		return resp, err
	}

	return resp, err
}

func (cm *ClientManager) CheckCredentials(c echo.Context, request *userReq.CheckCredentialsRequest) (*userRes.CheckCredentialsResponse, error) {
	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)

		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	resp, err := userClient.CheckCredentials(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while calling CheckCredentials: %v", err)

		return resp, err
	}

	return resp, err
}

func (cm *ClientManager) DeleteCredentials(c echo.Context, request *userReq.DeleteCredentialsRequest) (*userRes.DeleteCredentialsResponse, error) {
	userClient, err := cm.GetUserServiceClient(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client: %v", err)

		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	resp, err := userClient.DeleteCredentials(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while calling DeleteCredentials: %v", err)

		return resp, err
	}

	return resp, err
}

func (cm *ClientManager) BulkCardStatusUpdate(c echo.Context, request *userReq.BulkCardStatusUpdateRequest) (*userRes.BulkCardStatusUpdateResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service for bulk card status update, request :: %v", request)

	conn, err := cm.getUserServiceGrpcConn(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf(clients.UserClientConnError, err)
		return nil, err
	}

	cardClient := user.NewCardClient(conn)
	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	uploadClient, err := cardClient.BulkCardStatusUpdate(apiCtx)
	if err != nil {
		return nil, cm.logAndReturnError(c, fmt.Errorf("error while calling user service for bulk card status update, request :: %v", request))
	}

	if sendErr := uploadClient.Send(request); sendErr != nil {
		return nil, cm.logAndReturnError(c, sendErr)
	}

	if closeErr := uploadClient.CloseSend(); closeErr != nil {
		return nil, cm.logAndReturnError(c, closeErr)
	}

	incomingResponse, finalErr := uploadClient.Recv()
	if finalErr != nil {
		return nil, cm.logAndReturnError(c, finalErr)
	}

	return incomingResponse, nil
}

func (cm *ClientManager) logAndReturnError(c echo.Context, err error) error {
	cm.logger.WithContext(c).Error(err)
	return err
}

func (cm *ClientManager) UpdateUser(c echo.Context, req *userReq.UpdateUserRequest) (*userRes.UpdateUserResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service for update user, request :: %v", req)
	userClient, err := cm.GetUserServiceClient(c)

	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while fetching user service client : %v", err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	response, err := userClient.UpdateUser(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while calling user service for create user, with error: %s and request :: %v", err.Error(), req)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("update user executed successfully, request :: %v", req)

	return response, nil
}

func (cm *ClientManager) EvaluateStudentEligibility(c echo.Context, request *userReq.EvaluateStudentEligibilityRequest) (*userRes.EvaluateStudentEligibilityResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service :: EvaluateStudentEligibility , request :: %s ", request)

	conn, err := cm.getUserServiceGrpcConn(c)

	if err != nil {
		return nil, err
	}

	userClient := user.NewVerificationResourceClient(conn)
	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := userClient.EvaluateStudentEligibility(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) SubmitDocuments(c echo.Context, request *userReq.SubmitDocumentsRequest) (*userRes.SubmitDocumentsResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service :: SubmitDocuments , request :: %s ", request)

	conn, err := cm.getUserServiceGrpcConn(c)

	if err != nil {
		return nil, err
	}

	userClient := user.NewVerificationResourceClient(conn)
	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := userClient.SubmitDocuments(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdateVerificationResourceByUserID(c echo.Context, req *userReq.UpdateVerificationResourceByUserIDRequest) (*userRes.UpdateVerificationResourceByUserIDResponse, error) {
	cm.logger.WithContext(c).Infof("calling user service :: UpdateVerificationResourceByUserID , request :: %s ", req)

	conn, err := cm.getUserServiceGrpcConn(c)
	if err != nil {
		cm.logger.WithContext(c).Errorf(clients.UserClientConnError, err)
		return nil, err
	}

	verificationResourceClient := user.NewVerificationResourceClient(conn)
	conf := config.GetClientConfigs(clients.UserServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := verificationResourceClient.UpdateVerificationResourceByUserID(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

package clients

import (
	authz "github.com/Allen-Career-Institute/common-protos/authorization/v1"
	"github.com/Allen-Career-Institute/common-protos/authorization/v1/request"
	"github.com/Allen-Career-Institute/common-protos/authorization/v1/response"
	"github.com/Allen-Career-Institute/common-protos/authorization/v1/types"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	internal "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/labstack/echo/v4"
)

func (cm *ClientManager) getAuthzPdpClient(c echo.Context, _ *config.Config) (authz.AuthZEngineClient, error) {
	conn, err := cm.Grpc.GetConn(c, cm.logger, clients.AuthorizationPdpServiceClient, cm.cfg)
	if err != nil {
		return nil, err
	}

	client := authz.NewAuthZEngineClient(conn)

	return client, nil
}

func (cm *ClientManager) EnforceRbac(c echo.Context, cnf *config.Config, resource types.ResourceTypes, action types.Action) (*response.GetDecisionResponse, error) {
	cm.logger.WithContext(c).Infof("EnforceRbac: Request with resource: %v, action: %v", resource.String(), action.String())

	authZClient, err := cm.getAuthzPdpClient(c, cnf)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.AuthorizationPdpServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	tenantID, err := internal.GetTenantID(c)
	if err != nil {
		return nil, err
	}

	userID, err := internal.GetUserID(c)
	if err != nil {
		return nil, err
	}

	cm.logger.WithContext(c).Infof("EnforceRbac: user: %v, tenant: %v", userID, tenantID)
	sores, err := authZClient.GetDecision(apiCtx, &request.GetDecisionRequest{
		TenantId: tenantID,
		UserId:   userID,
		ResourceAttributes: &types.ResourceAttributes{
			Resource: resource,
			Attrs:    nil,
		},
		Action:   action,
		RbacOnly: true,
	})

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("EnforceRbac: Response: %v", sores.Pass)

	return sores, nil
}

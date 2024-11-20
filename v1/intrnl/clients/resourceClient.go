// nolint:staticcheck // deprecated function may be used
package clients

import (
	pb "github.com/Allen-Career-Institute/common-protos/resource/v1"
	pbrq "github.com/Allen-Career-Institute/common-protos/resource/v1/request"
	pbrs "github.com/Allen-Career-Institute/common-protos/resource/v1/response"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/labstack/echo/v4"
	grpcclient "google.golang.org/grpc"
)

func (cm *ClientManager) getGRPCConn(c echo.Context) (grpcclient.ClientConnInterface, error) {
	conn, err := cm.Grpc.GetConn(c, cm.logger, clients.ResourceServiceClient, cm.cfg)
	if err != nil {
		cm.logger.WithContext(c).Errorf("failed to get user conflict resolution client conn, err: %v", err)
		return nil, err
	}

	return conn, nil
}

func (cm *ClientManager) GetStudentBatchDetails(c echo.Context, _ *config.Config, request *pbrq.GetStudentBatchDetailsRequest) (*pbrs.GetStudentBatchDetailsResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service ::GetStudentBatchDetails, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	batchMappingClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchMappingClient.GetStudentBatchDetails(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error occurred in GetStudentBatchDetails , Err: %v", err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Response received from GetStudentBatchDetails is %s ", response.String())

	return response, nil
}

func (cm *ClientManager) GetAncestorsOfAFacility(c echo.Context, request *pbrq.GetAncestorsOfAFacilityRequest) (*pbrs.GetAncestorsOfAFacilityResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get ancestors of a facilities, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	facilityClient := pb.NewFacilityClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	response, err := facilityClient.GetAncestorsOfAFacility(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

package clients

import (
	calReq "github.com/Allen-Career-Institute/common-protos/cal/v1/request"
	calRes "github.com/Allen-Career-Institute/common-protos/cal/v1/response"
	resourceReq "github.com/Allen-Career-Institute/common-protos/resource/v1/request"
	resourceRes "github.com/Allen-Career-Institute/common-protos/resource/v1/response"
	userRes "github.com/Allen-Career-Institute/common-protos/user_management/v1/response"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	"github.com/labstack/echo/v4"
)

type Manager interface {
	// Resource Client
	GetAncestorsOfAFacility(c echo.Context, request *resourceReq.GetAncestorsOfAFacilityRequest) (*resourceRes.GetAncestorsOfAFacilityResponse, error)
	GetStudentBatchDetails(c echo.Context, _ *config.Config, request *resourceReq.GetStudentBatchDetailsRequest) (*resourceRes.GetStudentBatchDetailsResponse, error)

	// User Client
	GetUser(c echo.Context, cnf *config.Config, grpcHandler grpc.Manager) (*userRes.GetUserResponse, error)

	// Cal Client
	GetLearningMaterial(c echo.Context, cnf *config.Config, req *calReq.GetLearningMaterialRequest) (*calRes.GetLearningMaterialResponse, error)
	GetNAC(c echo.Context, cnf *config.Config, req *calReq.GetNACRequest) (*calRes.GetNACResponse, error)
}

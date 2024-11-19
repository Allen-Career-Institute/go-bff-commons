// nolint: revive,gocritic // This function is required to be exported
package utils

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"

	resourceRes "github.com/Allen-Career-Institute/common-protos/resource/v1/response"
	userRes "github.com/Allen-Career-Institute/common-protos/user_management/v1/response"
	"github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1/configs"
	"github.com/labstack/echo/v4"
)

type Manager interface {
	LoadCourseSyllabusResponseFromPreloadDS(ctx echo.Context, cnf *config.Config, log logger.Logger, courseID string, dsm *framework.DataSourceMappings,
		cm clients.Manager)
	LoadCourseSyllabusV2ResponseFromPreloadDS(ctx echo.Context, cnf *config.Config, log logger.Logger, courseID string, dsm *framework.DataSourceMappings,
		cm clients.Manager) (*resourceRes.GetCourseSyllabusV2Response, error)
	LoadStudentBatchDetailsFromPreloadDS(ctx echo.Context, log logger.Logger, cnf *config.Config, tenantID, userIDInContext string,
		dsm *framework.DataSourceMappings, cm clients.Manager)
	LoadUserDetailsFromPreloadDS(ctx echo.Context, log logger.Logger, cnf *config.Config, dsm *framework.DataSourceMappings, cm clients.Manager, grpc grpc.Manager) (*userRes.GetUserResponse, error)
	LoadOLTSCourseIDsFromPreloadDS(ctx echo.Context, log logger.Logger, cnf *config.Config, tenantID, userIDInContext string, dsm *framework.DataSourceMappings, cm clients.Manager) ([]string, error)
}

// for making mocks for UT
type DynamicConfig interface {
	Init(config *configs.Configuration) error
	Get(key string) (string, error)
	GetAsInterface(key string) (interface{}, error)
}

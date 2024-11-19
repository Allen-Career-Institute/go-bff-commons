// nolint:staticcheck // deprecated function may be used
package clients

import (
	pb "github.com/Allen-Career-Institute/common-protos/resource/v1"
	pbrq "github.com/Allen-Career-Institute/common-protos/resource/v1/request"
	pbrs "github.com/Allen-Career-Institute/common-protos/resource/v1/response"
	pbt "github.com/Allen-Career-Institute/common-protos/resource/v1/types"
	pbte "github.com/Allen-Career-Institute/common-protos/resource/v1/types/enums"
	"github.com/labstack/echo/v4"
	grpcclient "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	internal "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
)

func (cm *ClientManager) getGRPCConn(c echo.Context) (grpcclient.ClientConnInterface, error) {
	conn, err := cm.Grpc.GetConn(c, cm.logger, clients.ResourceServiceClient, cm.cfg)
	if err != nil {
		cm.logger.WithContext(c).Errorf("failed to get user conflict resolution client conn, err: %v", err)
		return nil, err
	}

	return conn, nil
}

func (cm *ClientManager) CreateCourse(c echo.Context, _ *config.Config, request *pbrq.CreateCourseRequest) (*pbrs.CreateCourseResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: create course, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.CreateCourse(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) CourseDetail(c echo.Context, request *pbrq.CourseRequest) (*pbrs.CourseDetailResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: create course, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.CourseDetail(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdateCourse(c echo.Context, _ *config.Config, request *pbrq.UpdateCourseRequest) (*pbrs.UpdateCourseResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update course, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.UpdateCourse(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetCoursesList(c echo.Context, request *pbrq.GetCoursesRequest) (*pbrs.GetCoursesResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get courses, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.GetCourses(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting all courses for request %s, Err: %v", request.String(), err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("GetCourses response is %s", response.String())

	return response, nil
}

func (cm *ClientManager) GetCoursesListWithV2Syllabus(c echo.Context, request *pbrq.GetCoursesRequestV2) (*pbrs.GetCoursesResponseV2, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get courses with v2 syllabus, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.GetCoursesWithV2Syllabus(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting all courses for request %s, Err: %v", request.String(), err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("GetCoursesWithV2Syllabus response is %s", response.String())

	return response, nil
}

func (cm *ClientManager) GetCourseSummary(c echo.Context, request *pbrq.CourseRequest) (*pbrs.CourseSummaryResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get courses summary, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.CourseSummary(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting course summary %s, Err: %v", request.String(), err)
		return nil, err
	}
	cm.logger.WithContext(c).Infof("GetCourseSummary response is %s", response.String())
	return response, nil
}

func (cm *ClientManager) GetCoursesListing(c echo.Context, request *pbrq.GetCoursesRequestV2) (*pbrs.GetCoursesResponseV2, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get courses listing, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.GetCoursesListing(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting all courses for request %s, Err: %v", request.String(), err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Get Courses Listing Successful")

	return response, nil
}

func (cm *ClientManager) GetCoursesFilter(c echo.Context, request *pbrq.GetCoursesRequestV2) (*pbrs.ResourceFilterResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get courses filter, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.GetCoursesFilter(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting all courses for request %s, Err: %v", request.String(), err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Get Courses Filter Successful")

	return response, nil
}

func (cm *ClientManager) GetPhasesListing(c echo.Context, request *pbrq.GetPhasesRequestV2) (*pbrs.GetPhasesResponseV2, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get phases listing, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	phaseClient := pb.NewPhaseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := phaseClient.GetPhasesListing(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting all phases for request %s, Err: %v", request.String(), err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Get Phases Listing Successful")

	return response, nil
}

func (cm *ClientManager) GetPhasesFilter(c echo.Context, request *pbrq.GetPhasesRequestV2) (*pbrs.ResourceFilterResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get phases filter, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	phaseClient := pb.NewPhaseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := phaseClient.GetPhasesFilter(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting all phases for request %s, Err: %v", request.String(), err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Get Phases Filter Successful")

	return response, nil
}

func (cm *ClientManager) GetBatchesListing(c echo.Context, request *pbrq.GetBatchesRequestV2) (*pbrs.GetBatchesResponseV2, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get batches listing, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.GetBatchesListing(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting all batches for request %s, Err: %v", request.String(), err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Get Batches Listing Successful")

	return response, nil
}

func (cm *ClientManager) GetBatchesFilter(c echo.Context, request *pbrq.GetBatchesRequestV2) (*pbrs.ResourceFilterResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get batches filter, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.GetBatchesFilter(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting all batches for request %s, Err: %v", request.String(), err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Get Batches Filter Successful")

	return response, nil
}

func (cm *ClientManager) GetClassScheduleSummary(c echo.Context, request *pbrq.ClassSchedulesSummaryRequest) (*pbrs.ClassSchedulesSummaryResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get class schedule summary, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	scheduleClient := pb.NewClassScheduleClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := scheduleClient.GetClassScheduleSummary(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error in getting class schedule summary. Err: %v", err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) AddCourseSyllabus(c echo.Context, _ *config.Config, request *pbrq.AddCourseSyllabusRequest) (*pbrs.AddCourseSyllabusResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: add course syllabus, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseSyllabusClient := pb.NewCourseSyllabusClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseSyllabusClient.AddCourseSyllabus(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) AddCourseSyllabusV2(c echo.Context, _ *config.Config, request *pbrq.AddCourseSyllabusV2Request) (*pbrs.AddCourseSyllabusV2Response, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: add course syllabus v2, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseSyllabusClient := pb.NewCourseSyllabusClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseSyllabusClient.AddCourseSyllabusV2(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DeleteCourseSyllabus(c echo.Context, _ *config.Config, request *pbrq.DeleteCourseSyllabusRequest) (*pbrs.DeleteCourseSyllabusResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: delete course syllabus, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseSyllabusClient := pb.NewCourseSyllabusClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseSyllabusClient.DeleteCourseSyllabus(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetCourseSyllabus(c echo.Context, _ *config.Config, request *pbrq.GetCourseSyllabusRequest) (*pbrs.GetCourseSyllabusResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get course syllabuses, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseSyllabusClient := pb.NewCourseSyllabusClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseSyllabusClient.GetCourseSyllabus(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("GetCourseSyllabus| Error occurred in getting syllabus, Err: %v", err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("GetCourseSyllabus| Response received for courseSyllabus %v", response)

	return response, nil
}

func (cm *ClientManager) GetCourseSyllabusV2(c echo.Context, _ *config.Config, request *pbrq.GetCourseSyllabusV2Request) (*pbrs.GetCourseSyllabusV2Response, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get course syllabus v2 , request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseSyllabusClient := pb.NewCourseSyllabusClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseSyllabusClient.GetCourseSyllabusV2(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("GetCourseSyllabusV2 | Error occurred in getting syllabus, Err: %v", err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("GetCourseSyllabusV2 | Response received for courseSyllabus %v", response)

	return response, nil
}

func (cm *ClientManager) GetBatchSyllabus(c echo.Context, _ *config.Config, request *pbrq.GetBatchSyllabusRequest) (*pbrs.GetCourseSyllabusV2Response, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get batch syllabus , request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.GetBatchSyllabus(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("GetBatchSyllabus | Error occurred in getting syllabus, Err: %v", err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("GetBatchSyllabus | Response received for batchSyllabus %v", response)

	return response, nil
}

func (cm *ClientManager) EditCourseSyllabus(c echo.Context, _ *config.Config, request *pbrq.EditCourseSyllabusRequest) (*pbrs.EditCourseSyllabusResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: edit course syllabus , request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseSyllabusClient := pb.NewCourseSyllabusClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseSyllabusClient.EditCourseSyllabus(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("EditCourseSyllabus | Error occurred in edit syllabus, Err: %v", err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("EditCourseSyllabus | Response received for edit syllabus, response %v",
		response)

	return response, nil
}

func (cm *ClientManager) ValidateSyllabusNodeDeletion(c echo.Context, _ *config.Config, request *pbrq.ValidateSyllabusNodeDeletionRequest) (*pbrs.ValidateSyllabusNodeDeletionResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: validate course syllabus node deletion, request :: %s ",
		request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseSyllabusClient := pb.NewCourseSyllabusClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseSyllabusClient.ValidateSyllabusNodeDeletion(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("ValidateSyllabusNodeDeletion | Error occurred in validate syllabus node deletion, Err: %v",
			err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("ValidateSyllabusNodeDeletion | Response received for validate syllabus node deletion, response %v",
		response)

	return response, nil
}

func (cm *ClientManager) CreateCourseSyllabusFromExisting(c echo.Context, _ *config.Config, courseID, existingCourseID string) (bool, error) {
	cm.logger.WithContext(c).Infof("CreateCourseSyllabusFromExisting :: CourseID: %v, existingCourseID: %v",
		courseID, existingCourseID)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return false, err
	}

	courseSyllabusClient := pb.NewCourseSyllabusClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseSyllabusClient.CreateCourseSyllabusFromExisting(apiCtx, &pbrq.CreateCourseSyllabusFromExistingRequest{
		CourseId:         courseID,
		ExistingCourseId: existingCourseID,
	})
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return false, err
	}

	return response.Result, nil
}

func (cm *ClientManager) GetResourceMetaEntities(c echo.Context, _ *config.Config, request *pbrq.GetResourceMetaEntitiesRequest) (*pbrs.GetResourceMetaEntitiesResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get resource meta entities, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	resourceClient := pb.NewResourceClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := resourceClient.GetResourceMetaEntities(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
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

func (cm *ClientManager) GetStudentCourseDetails(c echo.Context, request *pbrq.StudentCourseDetailsRequest) (*pbrs.StudentCourseDetailsResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service ::GetStudentCourseDetails, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	batchMappingClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchMappingClient.StudentCourseDetails(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error occurred in GetStudentCourseDetails , Err: %v", err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Response received from GetStudentCourseDetails is %s ", response.String())

	return response, nil
}

func (cm *ClientManager) GetStudentBatchListing(c echo.Context, _ *config.Config, request *pbrq.GetBatchesRequestV2) (*pbrs.GetStudentBatchDetailsResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: GetStudentBatchListing, request :: %s ", request.String())

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	batchMappingClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchMappingClient.GetBatchStudents(apiCtx, request)
	if err != nil {
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Response received from GetStudentBatchListing is %s ", response.String())
	return response, nil
}

func (cm *ClientManager) GetStudentBatchChangeHistory(c echo.Context, _ *config.Config, request *pbrq.GetStudentBatchHistoryRequest) (*pbrs.GetStudentBatchHistoryResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: GetStudentBatchChangeHistory, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	batchMappingClient := pb.NewStudentBatchMappingClient(conn)
	apiCtx := utils.GetRequestCtx(c)

	response, err := batchMappingClient.GetStudentBatchHistory(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}
	cm.logger.WithContext(c).Infof("Response received from GetStudentBatchChangeHistory is %s ", response.String())
	return response, nil
}

func (cm *ClientManager) GetFacilities(c echo.Context, request *pbrq.GetFacilitiesRequest) (*pbrs.GetFacilitiesResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get facilities, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	facilityClient := pb.NewFacilityClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	response, err := facilityClient.GetFacilities(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)

		return nil, err
	}

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

func (cm *ClientManager) GetDescendantsOfAFacility(c echo.Context, request *pbrq.GetDescendantsOfAFacilityRequest) (*pbrs.GetDescendantsOfAFacilityResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get descendants of a facilities, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	facilityClient := pb.NewFacilityClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	response, err := facilityClient.GetDescendantsOfAFacility(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetFacilitiesFilter(c echo.Context, request *pbrq.GetFacilitiesRequestV2) (*pbrs.ResourceFilterResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get facilities, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	facilityClient := pb.NewFacilityClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	response, err := facilityClient.GetFacilitiesFilter(apiCtx, request)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetFacilitiesListing(c echo.Context, request *pbrq.GetFacilitiesRequestV2) (*pbrs.GetFacilitiesResponseV2, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get facilities, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	facilityClient := pb.NewFacilityClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	response, err := facilityClient.GetFacilitiesListing(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) AddFacility(c echo.Context, _ *config.Config, request *pbrq.AddFacilityRequest) (*pbrs.AddFacilityResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: add facility, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	facilityClient := pb.NewFacilityClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := facilityClient.AddFacility(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdateFacilities(c echo.Context, _ *config.Config, request *pbrq.UpdateFacilityRequest) (*pbrs.UpdateFacilityResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update facilities, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	facilityClient := pb.NewFacilityClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := facilityClient.UpdateFacility(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DeleteFacilities(c echo.Context, _ *config.Config, request *pbrq.DeleteFacilityRequest) (*pbrs.DeleteFacilityResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: delete facility, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	facilityClient := pb.NewFacilityClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := facilityClient.DeleteFacility(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) CreatePhase(c echo.Context, _ *config.Config, request *pbrq.CreatePhaseRequest) (*pbrs.CreatePhaseResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: create phase, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	phaseClient := pb.NewPhaseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := phaseClient.CreatePhase(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdatePhase(c echo.Context, _ *config.Config, request *pbrq.UpdatePhaseRequest) (*pbrs.UpdatePhaseResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update phase, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	phaseClient := pb.NewPhaseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := phaseClient.UpdatePhase(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DeletePhase(c echo.Context, _ *config.Config, request *pbrq.DeletePhaseRequest) (*pbrs.DeletePhaseResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update phase, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	phaseClient := pb.NewPhaseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := phaseClient.DeletePhase(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetPhases(c echo.Context, _ *config.Config, request *pbrq.GetPhasesRequest) (*pbrs.GetPhasesResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: create phase, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	phaseClient := pb.NewPhaseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := phaseClient.GetPhases(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetPhaseDetail(c echo.Context, _ *config.Config, request *pbrq.GetPhaseDetailRequest) (*pbrs.GetPhaseDetailResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get phase detail, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	phaseClient := pb.NewPhaseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := phaseClient.GetPhaseDetail(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) CreateBatch(c echo.Context, _ *config.Config, request *pbrq.CreateBatchRequest) (*pbrs.CreateBatchResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: create batch, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.CreateBatch(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DownloadBulkBatchCreateTemplate(c echo.Context, request *pbrq.BulkBatchCreateTemplateDownloadRequest) (*pbrs.BulkBatchCreateTemplateDownloadResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: download bulk batch create template, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.DownloadBulkBatchCreateTemplate(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	incomingResponse, err := response.Recv()
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return incomingResponse, nil
}

func (cm *ClientManager) UploadBulkBatchCreate(c echo.Context, request *pbrq.BulkBatchCreateUploadRequest) (*pbrs.BulkBatchCreateUploadResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: upload bulk batch create")

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	uploadClient, err := batchClient.UploadBulkBatchCreate(apiCtx)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	if sendErr := uploadClient.Send(request); sendErr != nil {
		cm.logger.WithContext(c).Error(sendErr)
		return nil, sendErr
	}

	if closeErr := uploadClient.CloseSend(); closeErr != nil {
		cm.logger.WithContext(c).Error(closeErr)
		return nil, closeErr
	}

	incomingResponse, finalErr := uploadClient.Recv()
	if finalErr != nil {
		cm.logger.WithContext(c).Error(finalErr)
		return nil, finalErr
	}

	return incomingResponse, nil
}

func (cm *ClientManager) UpdateBatch(c echo.Context, _ *config.Config, request *pbrq.UpdateBatchRequest) (*pbrs.UpdateBatchResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update batch, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.UpdateBatch(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DeleteBatch(c echo.Context, _ *config.Config, request *pbrq.DeleteBatchRequest) (*pbrs.DeleteBatchResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: delete batch, request :: %s ", request)
	conn, err := cm.getGRPCConn(c)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	apiCtx := utils.GetRequestCtx(c)

	response, err := batchClient.DeleteBatch(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GenerateBatchCodePrefixRequest(c echo.Context, _ *config.Config, request *pbrq.GenerateBatchCodePrefixRequest) (*pbrs.GenerateBatchCodePrefixResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: generate batch code prefix, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.GenerateBatchCodePrefix(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)

		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetBatches(c echo.Context, request *pbrq.GetBatchesRequest) (*pbrs.GetBatchesResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get batches, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.GetBatches(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetBatchesInBulk(c echo.Context, request *pbrq.BulkGetBatchRequest) (*pbrs.GetBatchesResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: GetBatchesInBulk, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.BulkGetBatches(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) ValidateBatchCode(c echo.Context, request *pbrq.ValidateBatchCodeRequest) (bool, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: ValidateBatchCode, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return false, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.ValidateBatchCode(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return false, err
	}

	return response.Value, nil
}

func (cm *ClientManager) GetBatchDetail(c echo.Context, request *pbrq.GetBatchDetailRequest) (*pbrs.GetBatchDetailResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: GetBatchDetail, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.GetBatchDetail(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("response is %v", response)

	return response, nil
}

func (cm *ClientManager) GetBatch(c echo.Context, request *pbrq.GetBatchRequest) (*pbrs.GetBatchResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: GetBatch, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchClient := pb.NewBatchClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := batchClient.GetBatch(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("response is %v", response)

	return response, nil
}

func (cm *ClientManager) AddCourseContent(c echo.Context, _ *config.Config, request *pbrq.AddCourseContentRequest) (*pbrs.AddCourseContentResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: add course content, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseContentClient := pb.NewCourseContentClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseContentClient.AddCourseContent(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) RemoveCourseContent(c echo.Context, _ *config.Config, request *pbrq.RemoveCourseContentRequest) (*pbrs.RemoveCourseContentResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: remove course content, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseContentClient := pb.NewCourseContentClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseContentClient.RemoveCourseContent(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdateCourseContent(c echo.Context, _ *config.Config, request *pbrq.UpdateCourseContentRequest) (*pbrs.UpdateCourseContentResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update course content, request :: %v ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseContentClient := pb.NewCourseContentClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseContentClient.UpdateCourseContent(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error("Error while updating course content %v", err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) CopyCourseContent(c echo.Context, _ *config.Config, request *pbrq.CopyCourseContentRequest) (*pbrs.CopyCourseContentResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: copy course content, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseContentClient := pb.NewCourseContentClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseContentClient.CopyCourseContent(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetCourseContent(c echo.Context, _ *config.Config, request *pbrq.GetCourseContentRequest) (*pbrs.GetCourseContentResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get course content, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseContentClient := pb.NewCourseContentClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseContentClient.GetCourseContent(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) StudentBatchMovement(c echo.Context, _ *config.Config, request *pbrq.StudentBatchMovementRequest) (*pbrs.StudentBatchMovementResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: student batch movement, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	studentBatchMappingClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := studentBatchMappingClient.StudentBatchMovement(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetStudentDetailsForBatch(c echo.Context, _ *config.Config, request *pbrq.GetStudentsForBatchRequest) (*pbrs.GetStudentsForBatchInResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get student details for batch, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	studentBatchMappingClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := studentBatchMappingClient.GetStudentDetailsForBatch(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetStudentsForBatchIn(c echo.Context, _ *config.Config, request *pbrq.GetStudentsForBatchInRequest) (*pbrs.GetStudentsForBatchInResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get students for batch in, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	studentBatchMappingClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := studentBatchMappingClient.GetStudentsForBatchIn(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetLecturePlanClient(c echo.Context) (pb.LecturePlanClient, error) {
	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	lecturePlanClient := pb.NewLecturePlanClient(conn)

	return lecturePlanClient, nil
}

func (cm *ClientManager) CreateLecturePlansWithMeta(c echo.Context, request *pbrq.LecturePlansWithMetaCreateRequest) (*pbrs.LecturePlansWithMetaCreateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: create lecture plans with meta, request :: %s ", request)

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.CreateLecturePlansWithMeta(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdateLecturePlanMeta(c echo.Context, request *pbrq.LecturePlansMetaUpdateRequest) (*pbrs.LecturePlansMetaUpdateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: create lecture plans with meta, request :: %s ", request)

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.UpdateLecturePlanMeta(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) AddLecturePlansToMeta(c echo.Context, request *pbrq.LecturePlansAddRequest) (*pbrs.LecturePlansToMetaAddResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: add lecture plans to meta, request :: %s ", request)

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.AddLecturePlansToMeta(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdateLecturePlans(c echo.Context, request *pbrq.LecturePlansUpdateRequest) (*pbrs.LecturePlansUpdateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update lecture plans, request :: %s ", request)

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.UpdateLecturePlans(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) ValidateLecturePlans(c echo.Context, request *pbrq.LecturePlansValidateRequest) (*pbrs.LecturePlansValidateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: validate lecture plans, request :: %s ", request)

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.ValidateLecturePlans(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetLecturePlanDetails(c echo.Context, request *pbrq.LecturePlansGetRequest) (*pbrs.LecturePlansGetResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get lecture plans, request :: %s ", request)

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.GetLecturePlanDetails(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetLecturePlanMetas(c echo.Context, request *pbrq.LecturePlanMetasGetRequest) (*pbrs.LecturePlanMetasGetResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get lecture plan metas, request :: %s ", request)

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.GetLecturePlanMetas(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DownloadLecturePlanTemplate(c echo.Context, request *pbrq.LecturePlansTemplateDownloadRequest) (*pbrs.LecturePlansTemplateDownloadResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: download lecture plan template, request :: %s ", request)

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.DownloadLecturePlanTemplate(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	incomingResponse, err := response.Recv()
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return incomingResponse, nil
}

func (cm *ClientManager) UploadLecturePlan(c echo.Context, md metadata.MD, request *pbrq.LecturePlansUploadRequest) (*pbrs.LecturePlansUploadResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: upload lecture plans")

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	uploadCtx := metadata.NewOutgoingContext(apiCtx, md)
	uploadClient, err := client.UploadLecturePlan(uploadCtx)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	if sendErr := uploadClient.Send(request); sendErr != nil {
		cm.logger.WithContext(c).Error(sendErr)
		return nil, sendErr
	}

	if closeErr := uploadClient.CloseSend(); closeErr != nil {
		cm.logger.WithContext(c).Error(closeErr)
		return nil, closeErr
	}

	incomingResponse, finalErr := uploadClient.CloseAndRecv()
	if finalErr != nil {
		cm.logger.WithContext(c).Error(finalErr)
		return nil, finalErr
	}

	return incomingResponse, nil
}

// GetCourseTopicNodes currently it's in lecture plan context but we'll separate it into two different APIs.
func (cm *ClientManager) GetCourseTopicNodes(c echo.Context, request *pbrq.LecturePlanTopicNodesRequest) (*pbrs.LecturePlanTopicNodesResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get course topic nodes")

	client, err := cm.GetLecturePlanClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	return client.GetLecturePlanTopicNodes(apiCtx, request)
}

func (cm *ClientManager) GetClassScheduleClient(c echo.Context) (pb.ClassScheduleClient, error) {
	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	classScheduleClient := pb.NewClassScheduleClient(conn)

	return classScheduleClient, nil
}

func (cm *ClientManager) GetClassSchedule(c echo.Context, request *pbrq.ClassScheduleGetRequest) (*pbrs.ClassScheduleGetResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.GetClassSchedule(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetClassSchedules(c echo.Context, request *pbrq.ClassSchedulesGetRequest) (*pbrs.ClassSchedulesGetResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.GetClassSchedules(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) GetClassSchedulesV2(c echo.Context, request *pbrq.ClassSchedulesGetRequestV2) (*pbrs.ClassSchedulesGetResponseV2, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get class schedule V2, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.GetClassScheduleListing(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DeleteClassSchedules(c echo.Context, request *pbrq.ClassSchedulesDeleteRequest) (*pbrs.ClassSchedulesDeleteResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	response, err := client.DeleteClassSchedules(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) ValidateClassSchedules(c echo.Context, request *pbrq.ClassScheduleValidateRequest) (*pbrs.ClassScheduleValidateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: validate class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.ValidateClassSchedules(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) BulkCreateClassSchedule(c echo.Context, request *pbrq.ClassScheduleBulkCreateRequest) (*pbrs.ClassScheduleBulkCreateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: bulk create class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.BulkCreateClassSchedules(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) BulkCreateClassScheduleV2(c echo.Context, request *pbrq.ClassScheduleBulkCreateRequest) (*pbrs.ClassScheduleBulkCreateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: bulk create class schedule V2, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.BulkCreateClassSchedulesV2(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) CreateClassSchedule(c echo.Context, request *pbrq.ClassScheduleCreateRequest) (*pbrs.ClassScheduleCreateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: bulk create class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.CreateClassSchedule(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdateClassSchedule(c echo.Context, request *pbrq.ClassScheduleUpdateRequest) (*pbrs.ClassScheduleUpdateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.UpdateClassSchedule(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UpdateClassScheduleStatusV2(c echo.Context, request *pbrq.ClassSchedulesStatusUpdateRequestV2) (*pbrs.ClassSchedulesStatusUpdateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: update class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.UpdateClassScheduleStatusV2(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) ReconcileClassSchedules(c echo.Context, request *pbrq.ClassSchedulesReconcileRequest) (*pbrs.ClassSchedulesReconcileResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: class schedule reconcile, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.ReconcileClassSchedules(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) BulkUpdateClassSchedule(c echo.Context, request *pbrq.ClassScheduleBulkUpdateRequest) (*pbrs.ClassScheduleBulkUpdateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: bulk update class schedule, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.BulkUpdateClassSchedules(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DownloadClassScheduleTemplate(c echo.Context, request *pbrq.ClassSchedulesTemplateDownloadRequest) (*pbrs.ClassSchedulesTemplateDownloadResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: download lecture plan template, request :: %s ", request)

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.DownloadClassSchedulesTemplate(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	incomingResponse, err := response.Recv()
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return incomingResponse, nil
}

func (cm *ClientManager) UploadClassSchedule(c echo.Context, md metadata.MD, request *pbrq.ClassSchedulesUploadRequest) (*pbrs.ClassSchedulesUploadResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: upload lecture plans")

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	uploadCtx := metadata.NewOutgoingContext(apiCtx, md)
	uploadClient, err := client.UploadClassSchedules(uploadCtx)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	if sendErr := uploadClient.Send(request); sendErr != nil {
		cm.logger.WithContext(c).Error(sendErr)
		return nil, sendErr
	}

	if closeErr := uploadClient.CloseSend(); closeErr != nil {
		cm.logger.WithContext(c).Error(closeErr)
		return nil, closeErr
	}

	incomingResponse, finalErr := uploadClient.CloseAndRecv()
	if finalErr != nil {
		cm.logger.WithContext(c).Error(finalErr)
		return nil, finalErr
	}

	return incomingResponse, nil
}

func (cm *ClientManager) GetClassScheduleColumns(c echo.Context, request *pbrq.ClassScheduleGetColumnsRequest) (*pbrs.ClassScheduleGetColumnsResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get class schedule columns")

	client, err := cm.GetClassScheduleClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	incomingResponse, err := client.GetClassScheduleColumns(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return incomingResponse, nil
}

func (cm *ClientManager) ReplaceTeacherMappings(c echo.Context, request *pbrq.ReplaceTeacherRequest) (*pbrs.ReplaceTeacherResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: ReplaceTeacherMappings, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	client := pb.NewDoubtTeacherMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.ReplaceTeacher(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DownloadDoubtMappingTemplateSpecialBatch(c echo.Context, request *pbrq.DownloadDoubtMappingTemplateRequest) (*pbrs.DownloadDoubtMappingSpecialBatchTemplateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: DownloadDoubtMappingTemplateSpecialBatch, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	client := pb.NewDoubtTeacherMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.DownloadDoubtMappingSpecialBatchTemplate(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	incomingResponse, err := response.Recv()
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return incomingResponse, nil
}

func (cm *ClientManager) DownloadDoubtMappingTemplate(c echo.Context, request *pbrq.DownloadDoubtMappingTemplateRequest) (*pbrs.DoubtMappingTemplateDownloadResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: DownloadDoubtMappingTemplate, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	client := pb.NewDoubtTeacherMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.DownloadDoubtMappingTemplate(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	incomingResponse, err := response.Recv()
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return incomingResponse, nil
}

func (cm *ClientManager) GetDoubtTeacherMapping(c echo.Context, request *pbrq.DoubtMappingRequest) (*pbrs.DoubtMappingResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: GetDoubtTeacherMapping request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	client := pb.NewDoubtTeacherMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := client.GetDoubtMapping(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) DoubtTeacherMappingUploadSpecialBatch(c echo.Context, md metadata.MD, request *pbrq.UploadDoubtMappingTemplateRequest) (*pbrs.UploadDoubtMappingTemplateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: DoubtTeacherMappingUploadSpecialBatch, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	client := pb.NewDoubtTeacherMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	uploadCtx := metadata.NewOutgoingContext(apiCtx, md)

	uploadClient, err := client.UploadDoubtMappingTemplateSpecialBatch(uploadCtx)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	if sendErr := uploadClient.Send(request); sendErr != nil {
		cm.logger.WithContext(c).Error(sendErr)
		return nil, sendErr
	}

	if closeErr := uploadClient.CloseSend(); closeErr != nil {
		cm.logger.WithContext(c).Error(closeErr)
		return nil, closeErr
	}

	incomingResponse, finalErr := uploadClient.CloseAndRecv()
	if finalErr != nil {
		cm.logger.WithContext(c).Error(finalErr)
		return nil, finalErr
	}

	return incomingResponse, nil
}

func (cm *ClientManager) DoubtTeacherMappingUpload(c echo.Context, md metadata.MD, request *pbrq.UploadDoubtMappingTemplateRequest) (*pbrs.UploadDoubtMappingTemplateResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: DoubtTeacherMappingUpload, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	client := pb.NewDoubtTeacherMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	uploadCtx := metadata.NewOutgoingContext(apiCtx, md)

	uploadClient, err := client.UploadDoubtMappingTemplate(uploadCtx)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	if sendErr := uploadClient.Send(request); sendErr != nil {
		cm.logger.WithContext(c).Error(sendErr)
		return nil, sendErr
	}

	if closeErr := uploadClient.CloseSend(); closeErr != nil {
		cm.logger.WithContext(c).Error(closeErr)
		return nil, closeErr
	}

	incomingResponse, finalErr := uploadClient.CloseAndRecv()
	if finalErr != nil {
		cm.logger.WithContext(c).Error(finalErr)
		return nil, finalErr
	}

	return incomingResponse, nil
}

func (cm *ClientManager) GetBatchIDToBatchMap(c echo.Context, batchIDs []string) (map[string]*pbrs.ResourceFilterResponse_Result, error) {
	tenantID, err2 := internal.GetTenantID(c)
	if err2 != nil {
		return nil, err2
	}

	if len(batchIDs) == 0 {
		return make(map[string]*pbrs.ResourceFilterResponse_Result), nil
	}

	request := &pbrq.GetBatchesRequestV2{
		TenantId: tenantID,
		Search:   make([]*pbt.Search, 0),
		Filter: &pbrq.GetBatchesRequestV2_Filter{
			Batches: &pbt.ListKey{
				Values: batchIDs,
				Op:     pbte.Operation_IN,
			}},
		Pagination: &pbt.Pagination{
			PageNumber: 1,
			PageSize:   int64(len(batchIDs)),
			Sort: &pbt.Pagination_Sort{
				By:    "createdAt",
				Order: pbte.SortOrder_ASC,
			},
		},
	}
	batches, err := cm.GetBatchesFilter(c, request)

	if err != nil {
		return nil, err
	}

	batchIDToBatchMap := make(map[string]*pbrs.ResourceFilterResponse_Result)

	for _, batch := range batches.GetResults() {
		batchIDToBatchMap[batch.Id] = batch
	}

	return batchIDToBatchMap, nil
}

func (cm *ClientManager) GetFacilityIDToFacilityMap(c echo.Context, facilityIDs []string) (map[string]*pbt.FacilityInfo, error) {
	tenantID, err := internal.GetTenantID(c)
	if err != nil {
		return nil, err
	}

	facilities, err := cm.GetFacilities(c, &pbrq.GetFacilitiesRequest{TenantId: tenantID, Facilities: facilityIDs})
	if err != nil {
		return nil, err
	}

	facilityIDToBatchMap := make(map[string]*pbt.FacilityInfo)

	for _, facility := range facilities.GetData() {
		facilityIDToBatchMap[facility.Id] = facility
	}

	return facilityIDToBatchMap, nil
}

func (cm *ClientManager) EnrollStudentToCourse(c echo.Context, request *pbrq.StudentEnrollRequestV2) (*pbrs.StudentEnrollResponseV2, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: enroll student to course")

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	batchMappingClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	student, err := batchMappingClient.EnrollStudentV2(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return student, nil
}

// UnenrollTestStudent NOTE: To be used only for test stage student users
func (cm *ClientManager) UnenrollTestStudent(c echo.Context, request *pbrq.StudentUnenrollRequest) error {
	cm.logger.WithContext(c).Info("calling resource service :: un-enroll test student from course")
	conn, err := cm.getGRPCConn(c)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return err
	}

	batchMappingClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	_, err = batchMappingClient.UnenrollStudent(apiCtx, request)

	return err
}

func (cm *ClientManager) AddUserSkillMapping(c echo.Context, request *pbrq.AddUserSkillRequest) (*pbrs.AddUserSkillMappingResponse, error) {
	cm.logger.WithContext(c).Info("calling resource service :: add user skill mapping")
	conn, err := cm.getGRPCConn(c)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	userSkillMappingClient := pb.NewUserSkillMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	userSkillMapping, err := userSkillMappingClient.AddUserSkillMapping(apiCtx, request)
	if err != nil {
		return nil, err
	}

	return userSkillMapping, nil
}

func (cm *ClientManager) GetUserSkillMappings(c echo.Context, request *pbrq.GetUserSkillRequest) (*pbrs.GetUserSkillMappingResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: get user skill mappings, request : %v", request)
	conn, err := cm.getGRPCConn(c)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	userSkillMappingClient := pb.NewUserSkillMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	userSkillMapping, err := userSkillMappingClient.GetUserSkillMapping(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return userSkillMapping, nil
}

func (cm *ClientManager) GetFreemiumCourse(c echo.Context, request *pbrq.GetFreemiumCourseRequest) (*pbrs.GetFreemiumCourseResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service ::GetFreemiumCourse, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	courseClient := pb.NewCourseClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := courseClient.GetFreemiumCourse(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error occurred in GetFreemiumCourse , Err: %v", err)
		return nil, err
	}

	cm.logger.WithContext(c).Infof("Response received from GetFreemiumCourse is %s ", response.String())

	return response, nil
}

func (cm *ClientManager) EnrollStudentsToSpecialBatch(c echo.Context, request *pbrq.SpecialBatchEnrollRequest) (*pbrs.SuccessCountResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: EnrollToSpecialBatch, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	sbmClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := sbmClient.EnrollStudentsToSpecialBatch(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

func (cm *ClientManager) UnEnrollStudentsFromSpecialBatch(c echo.Context, request *pbrq.SpecialBatchUnEnrollRequest) (*pbrs.SuccessCountResponse, error) {
	cm.logger.WithContext(c).Infof("calling resource service :: UnEnrollStudentsFromSpecialBatch, request :: %s ", request)

	conn, err := cm.getGRPCConn(c)
	if err != nil {
		return nil, err
	}

	sbmClient := pb.NewStudentBatchMappingClient(conn)
	conf := config.GetClientConfigs(clients.ResourceServiceClient, cm.cfg)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	response, err := sbmClient.UnEnrollStudentsFromSpecialBatch(apiCtx, request)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return response, nil
}

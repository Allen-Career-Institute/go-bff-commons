package pagehandler

import (
	"fmt"
	ps "github.com/Allen-Career-Institute/common-protos/page_service/v1"
	pbReq "github.com/Allen-Career-Institute/common-protos/page_service/v1/request"
	"github.com/Allen-Career-Institute/common-protos/page_service/v1/response"
	pbTypes "github.com/Allen-Career-Institute/common-protos/page_service/v1/types"
	resReq "github.com/Allen-Career-Institute/common-protos/resource/v1/request"
	resRes "github.com/Allen-Career-Institute/common-protos/resource/v1/response"
	resourceEnums "github.com/Allen-Career-Institute/common-protos/resource/v1/types/enums"
	userRes "github.com/Allen-Career-Institute/common-protos/user_management/v1/response"
	userTypes "github.com/Allen-Career-Institute/common-protos/user_management/v1/types"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients"
	pageds2 "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/datasources/pageds"
	models "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/models/commons"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/otel"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/metric"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/datasource"
	intrnl "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl"
	grpcClients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/models/page"
	log "github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
)

type pageDataHandler struct {
	cnf    config.Config
	dsm    framework.DatasourceMappingsManager
	logger log.Logger
	prm    pageds2.ResponseMapper
	meter  metric.Meter
	m      intrnl.Mapper
	eutil  utils.EchoUtil
	cm     clients.Manager
	grpc   grpc.Manager
}

func NewPageDataHandler(cfg *config.Config, dsm framework.DatasourceMappingsManager, logger *log.Logger, meter metric.Meter, m intrnl.Mapper, grpc grpc.Manager) pageds2.Handlers {
	return &pageDataHandler{
		cnf:    *cfg,
		dsm:    dsm,
		logger: *logger,
		prm:    *pageds2.NewResponseMapper(logger),
		meter:  meter,
		m:      m,
		eutil:  utils.NewEchoUtil(*logger),
		cm:     clients.NewClientManager(cfg, *logger, grpc),
		grpc:   grpc,
	}
}

func (pdh *pageDataHandler) getPageServiceClient(c echo.Context, cnf *config.Config) (ps.PageClient, error) {
	conn, err := pdh.grpc.GetConn(c, pdh.logger, grpcClients.PageServiceClient, cnf)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("error while getting page service client conn, err: %v", err)
		return nil, err
	}

	client := ps.NewPageClient(conn)
	return client, nil
}

func (pdh *pageDataHandler) GetDSList(c echo.Context, dsNames []string) (dsl []*datasource.DataSource) {
	for _, dsn := range dsNames {
		ds := pdh.dsm.GetDataSourceByName(dsn)
		if ds == nil {
			pdh.logger.WithContext(c).Errorf("error while getting data source for ds : %s", dsn)
			continue
		}
		dsl = append(dsl, ds)
	}
	return dsl
}

// GetPage TODO: Limit the number of widgets to 20 widgets per page,
// and include multiple pages in resp if no. of widgets is > 20
func (pdh *pageDataHandler) GetPage() datasource.HandlerFunc {
	return func(c echo.Context, cnf *config.Config) (models.DSResponse, error) {
		c, span := otel.Trace(c, "DataSource.GetPage")
		defer span.End()

		userDetails, err := pdh.cm.GetUser(c, cnf, pdh.grpc)
		if err != nil {
			pdh.logger.WithContext(c).Errorf("error while fetching user details, err: %v", err.Error())
		}
		userContext := map[string]string{"onboarded": "false", "enrolled": "false", "stream": "", "class": ""}
		if userDetails != nil {
			pdh.populateUserRelatedContext(c, userDetails, userContext)
		}
		pdh.populateAppVersionRelatedContext(c, userContext)
		pdh.populateCourseModuleRelatedContext(c, userContext)

		pdh.logger.WithContext(c).Infof("Fetching Page Data..")
		gpr := &page.GetPageRequest{UserContext: userContext}
		// validate request
		// TODO: add validations in request
		if err = utils.ReadRequest(c, gpr); err != nil {
			pdh.logger.WithContext(c).Errorf("Error while parsing request, err: %v", err)
			return intrnl.PopulateResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error()), err
		}
		// setting pageURL in context, so that it can be used in datasource (redirection use case)
		c.Set(utils.PageURL, gpr.PageURL)

		pdh.handleQueryParams(c, gpr.PageURL, userContext)

		// create page service req
		pageClient, err := pdh.getPageServiceClient(c, cnf)
		if err != nil {
			return models.DSResponse{}, err
		}

		request := getPageRequest(c, gpr)
		conf := config.GetClientConfigs(grpcClients.PageServiceClient, cnf)
		apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
		defer apiCancel()

		gpResp, err := pageClient.GetPageFromCache(apiCtx, request)
		if err != nil {
			pdh.logger.WithContext(c).Errorf("Error while fetching page from page service, URL: %s, err: %v", gpr.PageURL, err)
			code, msg := utils.HandleError(c, err, pdh.logger)
			if code >= http.StatusInternalServerError && code < http.StatusNetworkAuthenticationRequired {
				msg = "Error while fetching page from page service, URL: " + gpr.PageURL
				return intrnl.PopulateResponse(code, UserFacingMessage(code), msg), nil
			} else {
				return intrnl.PopulateResponse(code, UserFacingMessage(code), msg), nil
			}

		}

		setUrlMetaInContext(c, gpResp)
		pInfo := gpResp.PageInfo

		switch pInfo.PageMeta.PageType {
		case pbTypes.PageMeta_LIST:
			return pdh.handleListPage(c, pInfo)
		case pbTypes.PageMeta_TAB:
			return pdh.handleTabPage(c, pInfo)
		default:
			pdh.logger.WithContext(c).Errorf("unsupported pageType received in pageInfo for URL: %s", gpr.PageURL)
			return intrnl.PopulateResponse(http.StatusInternalServerError, utils.GenericError, ErrorUnsupportedPageType.Error()), err
		}
	}
}

func setUrlMetaInContext(c echo.Context, resp *response.GetPageReply) {
	urlMeta := resp.PageInfo.UrlMeta
	if urlMeta != nil {
		c.Set(utils.URLMeta, urlMeta)
	} else {
		c.Set(utils.URLMeta, nil)
	}
}

func (pdh *pageDataHandler) populateAppVersionRelatedContext(c echo.Context, userContext map[string]string) {
	appVersionCode := c.Request().Header.Get(utils.AppVersionCodeHeader)
	if appVersionCode != "" {
		userContext[utils.AppVersionCodeHeader] = appVersionCode
	}
}

func (pdh *pageDataHandler) populateUserRelatedContext(c echo.Context, userDetails *userRes.GetUserResponse, userContext map[string]string) {
	additionalInfo, isValid := (userDetails.Data.GetAdditionalInfo()).(*userTypes.UserInfo_StudentInfo)
	if isValid {
		if userDetails.Data.Status == userTypes.UserStatus_ONBOARDED {
			userContext[OnboardedContextCriteriaParam] = "true"
		}

		if userDetails.Data.Status == userTypes.UserStatus_ENROLLED {
			userContext[EnrolledContextCriteriaParam] = "true"
			// adding enrolled status to context
			c.Set("enrolled", true)

		}

		if userDetails.Data.Stage == userTypes.UserStage_TEST {
			userContext[InternalUserContextCriteriaParam] = "true"
		}

		if additionalInfo.StudentInfo.Stream.String() != "" {
			pdh.logger.WithContext(c).Infof("Stream in user context using GetUserResponse, stream: %s", additionalInfo.StudentInfo.Stream.String())
			userContext[StreamContextCriteriaParam] = additionalInfo.StudentInfo.Stream.String()
		}

		if additionalInfo.StudentInfo.CurrentClass.String() != "" {
			pdh.logger.WithContext(c).Infof("Class in user context using GetUserResponse, class: %s", additionalInfo.StudentInfo.CurrentClass.String())
			userContext[ClassContextCriteriaParam] = additionalInfo.StudentInfo.CurrentClass.String()
		}
	}
}

func (pdh *pageDataHandler) getStudentsInfo(c echo.Context) (*resRes.GetStudentBatchDetailsResponse, error) {

	uid, err := intrnl.GetUserID(c)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("error in GetUserID: %v", err)
		return nil, err
	}
	tenantID, err := intrnl.GetTenantID(c)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("error in GetTenantID: %v", err)
		return nil, err
	}
	selectedCourseID := c.QueryParams().Get(utils.SelectedCourseID)
	studentBatchDetailsReq := &resReq.GetStudentBatchDetailsRequest{
		TenantId:  tenantID,
		StudentId: uid,
	}

	if selectedCourseID != utils.EmptyString {
		studentBatchDetailsReq.CourseId = &selectedCourseID
	} else {
		pdh.logger.WithContext(c).Warnf("selectedCourseID is empty for userID: %s", uid)
	}

	studentBatchDetailsResponse, err := pdh.cm.GetStudentBatchDetails(c, &pdh.cnf, studentBatchDetailsReq)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("error in GetStudentBatchDetails: %v, student info's can not be embedded", err)
		return nil, err
	}
	return studentBatchDetailsResponse, nil
}

func (pdh *pageDataHandler) handleListPage(c echo.Context, pInfo *pbTypes.PageInfo) (models.DSResponse, error) {
	pageResp, err := pdh.processPageDetailsAndWidgetData(c, pInfo)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("Error while processing page details and widget data, err: %v", err)
		return intrnl.PopulateResponse(http.StatusInternalServerError, utils.GenericError, err.Error()), nil
	}
	return intrnl.PopulateResponse(http.StatusOK, http.StatusText(http.StatusOK), pageResp), nil
}

func (pdh *pageDataHandler) handleTabPage(c echo.Context, pInfo *pbTypes.PageInfo) (models.DSResponse, error) {
	pageResp, err := pdh.processPageDetailsAndWidgetData(c, pInfo)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("Error while processing page details and widget data, err: %v", err)
		return intrnl.PopulateResponse(http.StatusInternalServerError, utils.GenericError, err.Error()), nil
	}
	err = pdh.processTabData(c, pInfo, pageResp)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("Error while processing tab data, err: %v", err)
		return intrnl.PopulateResponse(http.StatusInternalServerError, utils.GenericError, err.Error()), nil
	}
	return intrnl.PopulateResponse(http.StatusOK, http.StatusText(http.StatusOK), pageResp), nil
}

func (pdh *pageDataHandler) processTabData(c echo.Context, pInfo *pbTypes.PageInfo, pageResp *page.CommonPageResponse) error {
	pageResp.TabData = pdh.prm.HandleTabPageResponse(c, pInfo)
	if len(pInfo.TabData) == 0 || len(pageResp.TabData) == 0 {
		return ErrorNoTabsToShow
	}
	// Now resolving page inside tabs for all selected tabs, resolving here to avoid cyclic dependency
	pdh.resolvePageWithinTabs(c, pInfo.TabData, pageResp)
	return nil
}

func (pdh *pageDataHandler) processPageDetailsAndWidgetData(c echo.Context, pageInfo *pbTypes.PageInfo) (*page.CommonPageResponse, error) {
	resolvedWidgetsMap := make(map[string]bool)
	pageResp, dsNames, wPosList, err := pdh.prm.MapResponse(c, pageInfo, resolvedWidgetsMap)
	if err != nil {
		return nil, fmt.Errorf("error in mapping pageinfo : %v, pageID: %d", err, pageInfo.Id)
	}

	preloadDS := pdh.getAllPreloadDS(c, dsNames, pageInfo.PageMeta.PreloadDataSources)

	pdh.processPreloadDataSources(&c, preloadDS)

	pageResp, err = pdh.processDataSources(c, pageResp, dsNames, wPosList, resolvedWidgetsMap)
	if err != nil {
		return nil, fmt.Errorf("error in processing data sources : %v, pageID: %d", err, pageInfo.Id)
	}
	// resolving visibility rules for widgets
	err = pdh.resolveVisibilityForWidgets(c, pageResp, resolvedWidgetsMap)
	if err != nil {
		return nil, fmt.Errorf(" error in resolving visibility for widgets : %v, pageID: %d", err, pageInfo.Id)
	}
	return pageResp, nil
}

func (pdh *pageDataHandler) getAllPreloadDS(c echo.Context, dsNames, pds []string) []string {
	// list to contain unique page level and ds level of preload datasource
	var dsList []string

	dsl := pdh.GetDSList(c, dsNames)

	// append page level preload datasource
	for _, dataSource := range dsl {
		dsList = append(dsList, dataSource.GetPreloadDS()...)
	}

	// append ds level preload datasource
	for _, dsName := range pds {
		if !utils.Contains(dsList, dsName) {
			dsList = append(dsList, dsName)
		}
	}

	return dsList
}

// getPageRequest TODO@Himanshu: Get User Context from User Service
func getPageRequest(c echo.Context, gpr *page.GetPageRequest) *pbReq.GetPageFromCacheRequest {
	r := &pbReq.GetPageFromCacheRequest{}
	urlString, _, _ := strings.Cut(gpr.PageURL, utils.QuestionString)
	r.Url = urlString
	clientType := c.Request().Header.Get("x-client-type")
	r.UserContext = make(map[string]string)
	r.UserContext["x-client-type"] = clientType
	if intrnl.IsUserLoggedIn(c) {
		r.UserContext["user-logged-in"] = strconv.FormatBool(intrnl.IsUserLoggedIn(c))
	}
	for k, v := range gpr.UserContext {
		r.UserContext[k] = v
	}
	return r
}

func (pdh *pageDataHandler) handleQueryParams(c echo.Context, urlString string, userContext map[string]string) {
	_, queryString, _ := strings.Cut(urlString, utils.QuestionString)
	values, err := url.ParseQuery(queryString)
	if err != nil {
		pdh.logger.WithContext(c).Error(err)
	} else {
		var queryParams = c.Request().URL.Query()
		for key := range values {
			for _, s := range values[key] {
				queryParams.Add(key, s)
				// adding query params to user context to customise response as per request
				_, ok := userContext[key]
				if !ok {
					userContext[key] = s
				}
			}
		}
		c.Request().URL.RawQuery = queryParams.Encode()
	}
}

func (pdh *pageDataHandler) resolvePageWithinTabs(c echo.Context, tabData []*pbTypes.TabContent, pageResp *page.CommonPageResponse) {
	for i, tab := range tabData {
		if tab.Selected {
			if tab.TabInfo != nil && tab.TabInfo.PageData != nil {
				listPageRes, err := pdh.processPageDetailsAndWidgetData(c, tab.TabInfo.PageData)
				if err != nil {
					pdh.logger.WithContext(c).Errorf("could not resolve data for pageID: %v within Tab, err: %v, skipping this page", tab.TabInfo.PageId, err)
					continue
				}
				pageResp.TabData[i].TabInfo.PageData = listPageRes
			}
		}
	}
}

func (pdh *pageDataHandler) resolveVisibilityForWidgets(c echo.Context, resp *page.CommonPageResponse, resolvedWidgetsMap map[string]bool) error {
	if resp == nil || resp.PageInfo.PageMeta.VisibilityRules == nil {
		return nil
	}
	var err error
	widgetVisibilityRules := resp.PageInfo.PageMeta.VisibilityRules.WidgetVisibilityRules
	resp.PageContent.HeaderWidgets, err = pdh.filterWidgetsBasedOnVisibilityRules(c, widgetVisibilityRules, resolvedWidgetsMap, resp.PageContent.HeaderWidgets)
	if err != nil {
		return err
	}
	resp.PageContent.Widgets, err = pdh.filterWidgetsBasedOnVisibilityRules(c, widgetVisibilityRules, resolvedWidgetsMap, resp.PageContent.Widgets)
	if err != nil {
		return err
	}
	resp.PageContent.FooterWidgets, err = pdh.filterWidgetsBasedOnVisibilityRules(c, widgetVisibilityRules, resolvedWidgetsMap, resp.PageContent.FooterWidgets)
	if err != nil {
		return err
	}
	resp.PageContent.OnloadWidgets, err = pdh.filterWidgetsBasedOnVisibilityRules(c, widgetVisibilityRules, resolvedWidgetsMap, resp.PageContent.OnloadWidgets)
	if err != nil {
		return err
	}
	return nil
}

func (pdh *pageDataHandler) filterWidgetsBasedOnVisibilityRules(c echo.Context, widgetVisibilityRules []*pbTypes.WidgetVisibilityRule, resolvedWidgetsMap map[string]bool, widgets []*page.WidgetData) ([]*page.WidgetData, error) {
	var widgetsToBeVisible []*page.WidgetData
	widgetsToBeVisible = make([]*page.WidgetData, 0)
	// traverse all the configured widgets
	for _, configuredWidget := range widgets {
		var visibilityRuleExists bool

		// traverse through all the visibility rules
		for _, rules := range widgetVisibilityRules {
			constWidgetIDFromVisibilityRule := rules.WidgetId

			// check if the visibility rule exists for the current configured widget
			if constWidgetIDFromVisibilityRule == configuredWidget.ConstWidgetID {
				visibilityRuleExists = true
				visibilityStatus := pdh.evaluateVisibilityConditions(c, rules, resolvedWidgetsMap)
				if visibilityStatus {
					widgetsToBeVisible = append(widgetsToBeVisible, configuredWidget)
				}
			}

		}
		// if the visibility rule does not exist for the current configured widget, consider it visible
		if !visibilityRuleExists {
			widgetsToBeVisible = append(widgetsToBeVisible, configuredWidget)
		}
	}
	return widgetsToBeVisible, nil
}

func (pdh *pageDataHandler) evaluateVisibilityConditions(c echo.Context, rules *pbTypes.WidgetVisibilityRule, resolvedWidgetsMap map[string]bool) bool {
	for _, condition := range rules.Rules {
		switch condition.Operator {
		case pbTypes.VisibilityRule_AND:
			res := pdh.evaluateANDConditions(c, condition, rules, resolvedWidgetsMap)
			if !res {
				return false
			}
		case pbTypes.VisibilityRule_OR:
			res := pdh.evaluateORConditions(c, condition, rules, resolvedWidgetsMap)
			if !res {
				return false
			}
		default:
			pdh.logger.WithContext(c).Warn("Unsupported Operator found in visibility condition")
		}
	}
	return true
}

func (pdh *pageDataHandler) evaluateORConditions(c echo.Context, rules *pbTypes.VisibilityRule, wvr *pbTypes.WidgetVisibilityRule, resolvedWidgetsMap map[string]bool) bool {
	for _, widgetRule := range rules.Conditions {
		wID := widgetRule.WidgetId

		widgetResolvedStatus := resolvedWidgetsMap[wID]
		if widgetRule.IsResolved == widgetResolvedStatus {
			return true // If any condition matches, return true
		}
	}
	pdh.logger.WithContext(c).Errorf("OR conditions did not satisfy for widget, dropping widgetID: %s", wvr.WidgetId)
	return false // If no condition matches, return false
}

func (pdh *pageDataHandler) evaluateANDConditions(c echo.Context, rules *pbTypes.VisibilityRule, wvr *pbTypes.WidgetVisibilityRule, resolvedWidgetsMap map[string]bool) bool {
	for _, widgetRule := range rules.Conditions {
		wID := widgetRule.WidgetId

		widgetResolvedStatus := resolvedWidgetsMap[wID]
		if widgetRule.IsResolved != widgetResolvedStatus {
			pdh.logger.WithContext(c).Errorf("AND conditions did not satisfy for widget, dropping widgetID: %s", wvr.WidgetId)
			return false // If any condition doesn't match, return false
		}
	}
	return true
}

func (pdh *pageDataHandler) populateCourseModuleRelatedContext(c echo.Context, userContext map[string]string) {
	uID, err := intrnl.GetUserID(c)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("error in GetUserID: %v, course module + center criteria can not be embedded", err)
		return
	}

	tenantID, err := intrnl.GetTenantID(c)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("error in GetTenantID: %v, course module + center criteria can not be embedded", err)
		return
	}

	studentBatchDetailsResponse, err := pdh.getStudentsInfo(c)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("error in GetStudentBatchDetails: %v, course module + center criteria can not be embedded", err)
		return
	}
	if studentBatchDetailsResponse == nil || len(studentBatchDetailsResponse.GetStudentBatchDetails()) == 0 {
		pdh.logger.WithContext(c).Errorf("studentBatchDetailsResponse resp is nil, course module + center criteria can not be embedded")
		return
	}
	// NOTE: A student can be part of multiple batches(REGULAR, SPECIAL etc), currently we are only considering Data from REGULAR batch
	for _, studentBatchDetail := range studentBatchDetailsResponse.GetStudentBatchDetails() {
		var centreName string

		if studentBatchDetail.GetBatchTypeEnum() == resourceEnums.BatchType_BATCH_REGULAR {
			getAncestorsOfAFacilityRequest := &resReq.GetAncestorsOfAFacilityRequest{
				TenantId:   tenantID,
				FacilityId: studentBatchDetail.GetCenterId(),
			}
			facilityAncestors, getAncestorsErr := pdh.cm.GetAncestorsOfAFacility(c, getAncestorsOfAFacilityRequest)
			if getAncestorsErr != nil {
				pdh.logger.WithContext(c).Errorf("GetAncestorsOfAFacility err: %v", getAncestorsErr)
				continue
			}
			for _, facilityAncestor := range facilityAncestors.Results {
				if facilityAncestor.GetType() == resourceEnums.FacilityType_FACILITY_TYPE_CENTER {
					centreName = facilityAncestor.GetName()
					break
				}
			}
		} else if studentBatchDetail.GetBatchTypeEnum() == resourceEnums.BatchType_BATCH_DEFAULT {
			centreName = studentBatchDetail.GetCenterName()
		} else {
			continue
		}

		c.Set(CenterNameContextCriteriaParam, centreName)
		// embedding centerName as a criteria
		if centreName == utils.EmptyString {
			pdh.logger.WithContext(c).Infof("centre name is empty for userID: %s, batchType: %s", uID, studentBatchDetail.GetBatchTypeEnum().String())
		} else {
			userContext[CenterNameContextCriteriaParam] = centreName
		}

		userContext[CourseModeContextCriteriaParam] = studentBatchDetail.GetCourseModeEnum().String()
		userContext[CourseTypeContextCriteriaParam] = studentBatchDetail.GetCourseType()
		userContext[CourseIDContextCriteriaParam] = studentBatchDetail.GetCourseId()
		userContext[PhaseIDContextCriteriaParam] = studentBatchDetail.GetPhaseId()
		userContext[PhaseNumberContextCriteriaParam] = studentBatchDetail.GetPhaseNumber()
		studentSession := studentBatchDetail.GetSession()
		session := ConvertDateRange(studentSession)
		if session == studentSession {
			pdh.logger.WithContext(c).Errorf("invalid session format: %v", session)
		}
		userContext[SessionContextCriteriaParam] = session
		userContext[BatchTypeEnumContextCriteriaParam] = studentBatchDetail.GetBatchTypeEnum().String()
		userContext[BatchCodeContextCriteriaParam] = studentBatchDetail.GetBatchCode()
		userContext[FacilityCodeContextCriteriaParam] = studentBatchDetail.GetFacilityCode()
		// TODO: check if we need to handle courseModule considering it as list of course Modules or single course module
		// Note: Currently we are handling considering only TEST_SERIES to be part of course_modules list currently
		for _, courseModule := range studentBatchDetail.GetCourseModules() {
			if courseModule.Type == resourceEnums.CourseModuleType_COURSE_MODULE_TYPE_TEST_SERIES {
				userContext[CourseModuleModeContextCriteriaParam] = courseModule.GetMode().String()
				userContext[CourseModuleTypeContextCriteriaParam] = courseModule.GetType().String()
				return
			}
		}
	}
	val, err := pdh.cnf.DynamicConfig.Get(utils.MultiCourseEnabled)
	if err == nil && val == "true" {
		for _, studentBatchDetail := range studentBatchDetailsResponse.GetStudentBatchDetails() {
			if studentBatchDetail.GetBatchTypeEnum() == resourceEnums.BatchType_BATCH_REGULAR {
				userContext[StreamContextCriteriaParam] = studentBatchDetail.Stream.String()
				userContext[ClassContextCriteriaParam] = studentBatchDetail.ClassEnum.String()
				pdh.logger.WithContext(c).Infof("Stream & Class in user context using GetStudentBatchDetails , stream: %s, class: %s", studentBatchDetail.Stream.String(), studentBatchDetail.ClassEnum.String())
				return
			}
		}
	}

}

// ConvertDateRange converts "04-2023 - 03-2024" to "04_2023__03_2024"
func ConvertDateRange(input string) string {
	cleanInput := strings.ReplaceAll(input, " ", "")

	parts := strings.Split(cleanInput, "-")
	if len(parts) == 4 {
		return parts[0] + "_" + parts[1] + "__" + parts[2] + "_" + parts[3]
	} else {
		return input
	}
}

package utils

import (
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	DeviceID                             = "X-Device-ID"
	VisitorID                            = "X-Visitor-Id"
	DeviceType                           = "X-Client-Type"
	IPAddress                            = "X-Forwarded-For"
	ServiceHeader                        = "x-md-service-name"
	AccessTokenHeader                    = "X-ACCESS-TOKEN"
	RefreshTokenHeader                   = "X-REFRESH-TOKEN"
	ReferrerHeader                       = "x-referrer"
	AppVersionCodeHeader                 = "X-Client-App-Version-Code"
	DeviceTypeWeb                        = "web"
	DeviceTypeiOS                        = "iOS"
	DeviceTypeAndroid                    = "android"
	DeviceTypeMweb                       = "mweb"
	AndroidAppVersionForDownloadsKey     = "android_app_version_for_downloads"
	IOSAppVersionForDownloadsKey         = "ios_app_version_for_downloads"
	DefaultAndroidAppVersionForDownloads = "67"
	DefaultIOSAppVersionForDownloads     = "37"
	IOSAppVersionForDownloads            = 37
	AndroidAppVersionForDownloads        = 67
)

const (
	EnvProd    = "prod"
	EnvStage   = "stage"
	EnvSandbox = "sandbox"
	Env        = "ENV"
)
const (
	ExternalNavigation = "EXTERNAL_NAVIGATION"
	SharedDataSource   = "shared_datasource"
	ResolveLMMWidget   = "ResolveLMMWigdets"
	Polymorphic        = "POLYMORPHIC"
)

// Error messages
const (
	GenericError                  = "something went wrong. Please try again"
	NotEligibileForPurchase       = "Course Not Eligible for Purchase"
	ContentIDMissing              = "content id is missing"
	DeviceTypeMissing             = "device type is missing"
	UserPrivilegesNotFoundError   = "Error while fetching user privileges"
	InvalidCartStatus             = "Invalid cart status"
	UserPrivilegesSheetNotFound   = "'User Privileges' Sheet is missing or sheet name is incorrect."
	InvalidUserPrivilegeExcelData = "'Employee ID' or 'Role Name' columns are missing / in incorrect order or column names are incorrect. " +
		"Please ensure that column names and order match the sample file."
	UserNotFoundError           = "error while fetching user details"
	ListingsNotFoundError       = "Error while fetching Listing titles"
	MissingUserID               = "UserID is missing"
	MissingPhoneNumber          = "PhoneNumber is missing"
	MissingTenantID             = "TenantID is missing"
	NotFreemiumUser             = "user is not a freemium user"
	TestResultsUnAvailable      = "Tests results unavailable"
	TestResultsAlreadyDeclared  = "Test results already declared"
	SheetNotFoundError          = "Sheet not found in excel file"
	FileNotFoundInRequestError  = "File not found in request"
	SheetStructureInvalidError  = "Sheet structure is not valid"
	ValidationError             = "Validation Error"
	CourseSyllabusDataNotFound  = "Course Syllabus data is missing"
	BatchSyllabusDataNotFound   = "batch syllabus data is missing"
	NoFlashcardsFound           = "no flashcards found"
	SelectedBatchNotEligible    = "selected batch is not eligible"
	NoEligibleBatchFound        = "no eligible batch found"
	InternalServerError         = "Internal Server Error"
	DifficultyLevelMissing      = "difficulty level is missing"
	ScoreMissing                = "score is missing"
	ParsingError                = "unable to parse provided input"
	AccessDenied                = "Access Denied"
	InvalidEmployeeIDs          = "Invalid Employee IDs"
	UserSkillMappingUploadError = "Error while uploading user skill mapping"
	NoQuestionsFoundError       = "No questions available"
	TestCorrectionError         = "Invalid state to Initiate Test Correction"
	TestPaperMappingError       = "No tests for paper ids"
	BatchNotFoundError          = "Batch not found for the given batches"
	PageAccessDeniedError       = "PAGE_ACCESS_DENIED"
	InvalidRequestType          = "Invalid request type"
	BadRequest                  = "bad request"
	NoQuestionsForCriteriaError = "No questions available for selected criteria"
	FailedToFetchOLTSCourses    = "Failed to fetch the OLTS courses list. Please try again later."
	FailedToRetrieveOLTSContent = "Failed to retrieve the OLTS course content. Please try again later."
	NoOLTSContentFound          = "No content found for the specified OLTS course"
	NotEligibleForPurchase      = "Course Not Eligible for Purchase"
)

// BatchCreateTests file constants
const (
	DefaultExcelSheetName                  = "Sheet1"
	BulkTestCreationFailedExcelFileName    = "bulk_tests_failed.xlsx"
	BulkTestCreationSucceededExcelFileName = "bulk_tests_succeeded.xlsx"
)

const (
	Tracer              = "tracer-bff-service"
	SpanTracer          = "bff-span-context"
	TracerServiceName   = "bff-layer"
	Tag                 = "tag"
	OtelGrpcTemporality = "otelGrpcTemporality"
	Cumulative          = "cumulative"
	Delta               = "delta"
)

const (
	Count             = "_count"
	Duration          = "_duration"
	BffDsMetricPrefix = "bff_ds_request"
	Meter             = "meter-"
	ServiceName       = "service_name"
	ServiceEnv        = "service_env"
	URI               = "uri"
	StatusCode        = "status_code"
	MetricPrefix      = "bff_request"
	DataSourceName    = "datasource"
)

const (
	AuthServTenantID           = "aUSsW8GI03dyQ0AIFVn92"
	AuthUserID                 = "uid"
	AuthTenantID               = "tid"
	AuthPersonaType            = "pt"
	AuthSessionID              = "sid"
	ExternalID                 = "e_id"
	InternalStudentUserInToken = "isu"
	ID                         = "id"
	UserID                     = "userID"
	InternalStudentUser        = "internalUser"
	TenantID                   = "tenantID"
	PersonaType                = "personaType"
	SessionID                  = "sessionID"
	UserExternalID             = "userExternalId"
	LoggedIn                   = "loggedIn"
	XAccessToken               = "X-Access-Token"
	XRefreshToken              = "X-Refresh-Token"
	ClientRequestID            = "client_request_id"
	UserEnrolledStatus         = "enrolled"
	PageURL                    = "page_url"
	CenterName                 = "center-name"
	StudentID                  = "student_id"
	ReportID                   = "report_id"
	SelectedCourseID           = "selected_course_id"
	SelectedBatchList          = "selected_batch_list"
	MultiCourseEnabled         = "isMultiCourseEnabled"
	KotaCentreName             = "KOTA"
)

const (
	MeetingIDParam          = "meeting_id"
	StreamModeQueryParam    = "stream_mode"
	UserPersonaTypeParam    = "persona_type"
	PersonaTypeStudent      = "STUDENT"
	PersonaTypeTeacher      = "TEACHER"
	PersonaTypeInternalUser = "INTERNAL_USER"
	Enrolled                = "ENROLLED"
	NAVIGATION              = "NAVIGATION"
	ShowBottomSheet         = "SHOW_BOTTOMSHEET"
	StartTime               = "start_time"
	EndTime                 = "end_time"
	PptConvertorClientId    = "client_id"
)

const (
	CommaString     = ","
	SpaceString     = " "
	EmptyString     = ""
	QuestionString  = "?"
	HyphenString    = "-"
	ConfigureString = " configured for "
)

const (
	PageIDParam               = "page_id"
	WidgetIDParam             = "widget_id"
	URLPageLinkIDParam        = "url_page_link_id"
	OffSet                    = "offset"
	Limit                     = "limit"
	Version                   = "version"
	URL                       = "url"
	WidgetDefinitionIDParam   = "widget_definition_id"
	WidgetDefinitionNameParam = "name"
	TabIDParam                = "tab_id"
	TestStatusParam           = "test_status"
	ViewAll                   = "View All"
)

const (
	WidgetClientNotFoundErrorMsg           = "widget client is nil, cannot proceed"
	WidgetDefinitionClientNotFoundErrorMsg = "widget definition client is nil, cannot proceed"
	PageClientNotFoundErrorMsg             = "grpc connection to page service failed, cannot proceed"
	TabClientNotFoundErrorMsg              = "tab client is nil, cannot proceed"
	URLClientNotFoundErrorMsg              = "URL client is nil, cannot proceed"
	TaxonomyNodeNotFoundError              = "rpc error: code = NotFound desc = Data Not Found"
)

const (
	JuspayOrderPaymentStatusEndpoint = "/checkout/juspay/order-payment-status"
	AppStorePurchaseStatusEndpoint   = "/checkout/in-app-purchase/signal-workflow"
	StudentDashboardSection          = "section"
	EnquiryFormEndpoint              = "/marketing/enquiry"
	SFAdsatRegistrationEndpoint      = "/adsat/registration"
	EnquiryFormOtpVerified           = "Enquiry form otp verified"
	PostPurchase                     = "POST_PURCHASE"
	WebviewPostpurchaseURI           = "/postpurchase"
)

const (
	GotErrorLogConstant = "got error: %v"
	TextPositive        = "text-positive"
	Success             = "success"
	Warning             = "warning"
)

const (
	URLMeta                    = "url_meta"
	WidgetData                 = "widget_data"
	WidgetIndexToWidgetDataMap = "widget_index_to_widget_data_map"
)

const (
	IsEnrolledTrue     = "1"
	IsEnrolledFalse    = "0"
	CareCornerPageName = "Care Corner Screen"
	CareCornerTabName  = "Break"
)

const (
	CardName      = "card_name"
	CardType      = "card_type"
	CurrentScreen = "current_screen"
	CardV2Name    = "button_name"
	CardV2Type    = "button_type"
	EventType     = "event_type"
)

const (
	ImageMediaType                                                  = "images"
	AudioMediaType                                                  = "audios"
	VideoMediaType                                                  = "videos"
	Text                                                            = "text"
	Timestamp                                                       = "timestamp"
	IsSender                                                        = "is_sender"
	ExpandText                                                      = "expand_text"
	ReportCta                                                       = "report_cta"
	MediaList                                                       = "media_list"
	TransferredToTeacherDoubtReplyText                              = "Doubt transferred to Teacher"
	WebPrimaryMessage                                               = "web_primary_message"
	TransferToTeacherChipMessage                                    = "Now ask follow-up questions with the teacher"
	BackgroundColor                                                 = "background_color"
	HexCode                                                         = "hex_code"
	Name                                                            = "name"
	TransferredToTeacherChipWebPrimaryMessageBackgroundColorName    = "bg-quaternary-light"
	TextColor                                                       = "text_color"
	TransferredToTeacherChipWebPrimaryMessageTextColorName          = "text-primary"
	PrimaryMessage                                                  = "primary_message"
	TransferredToTeacherChipAppPrimaryMessageBackgroundColorHexCode = "#3D444E"
	TransferredToTeacherChipAppPrimaryMessageTextColorHexCode       = "#FFFFFF"
	TransferredToTeacherDottedLineTitle                             = "Transferred to Teacher"
	TransferredToTeacherDottedLineIconName                          = "user_octagon_bold"
)

const (
	PreviousAcademicSession = 2
	NextAcademicSession     = 1
)

const (
	OnlineTestSeries = "Online Test Series"
)

func GetStudentDashboardMap() map[string]string {
	return map[string]string{
		"profile":      "StudentDashboardProfileSection",
		"courses":      "StudentDashboardBatchSection",
		"purchase":     "StudentDashboardPurchaseSection",
		"courseChange": "StudentDashboardCourseChangeListingSection",
	}
}

func GetListOfCustomHeadersAsCommaSeparated() string {
	headerList := []string{AccessTokenHeader, RefreshTokenHeader}
	cs := strings.Join(headerList, CommaString)

	return cs
}

func UserLocationAsiaKolkata() *time.Location {
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Errorf("error %v while loading IST location, switching to time.Now().Location()", err)
		return time.Now().Location()
	}

	return location
}

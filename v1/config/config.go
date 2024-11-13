package config

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"

	dc "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1"
)

// Config App config struct
type Config struct {
	Server                          ServerConfig
	AppConfig                       AppConfig
	DynamicConfig                   dc.DynamicConfig
	Logger                          Logger
	GoPool                          GoPool
	DataSource                      DataSource
	PlaylistFilenameConfig          PlaylistConfig
	LmmRedisSecretLocation          string
	AkamaiConfig                    AkamaiConfig
	WhiteListSubTypesForContentAuth WhiteListSubTypesForContentAuth
}

// DataSource consists of all datasource configs
type DataSource struct {
	FetchReportsDS                                   models.DataSourceConfig
	SupersetLoginDS                                  models.DataSourceConfig
	VerifyOtpDS                                      models.DataSourceConfig
	LoginWithUsernamePasswordDS                      models.DataSourceConfig
	SendOtpDS                                        models.DataSourceConfig
	LogoutDS                                         models.DataSourceConfig
	OnboardUserDS                                    models.DataSourceConfig
	GetUserIdentitiesDS                              models.DataSourceConfig
	GetUsersByIdentityDS                             models.DataSourceConfig
	UpdateUserPIIDS                                  models.DataSourceConfig
	AdminAddParentDS                                 models.DataSourceConfig
	DeleteCredentialsDS                              models.DataSourceConfig
	ResetPasswordDS                                  models.DataSourceConfig
	GetStudentHealthAppInfoDS                        models.DataSourceConfig
	EditProfileDS                                    models.DataSourceConfig
	DeleteProfileDS                                  models.DataSourceConfig
	PreCheckoutUserDS                                models.DataSourceConfig
	PostCheckoutUserDS                               models.DataSourceConfig
	UpdateContextDS                                  models.DataSourceConfig
	GetContextDS                                     models.DataSourceConfig
	LoginWithGoogleDS                                models.DataSourceConfig
	GetGoogleRedirectLinkDS                          models.DataSourceConfig
	RefreshTokenDS                                   models.DataSourceConfig
	AddAddressDS                                     models.DataSourceConfig
	GetAddressByIDDS                                 models.DataSourceConfig
	GetAllAddressesDS                                models.DataSourceConfig
	UpdateAddressDS                                  models.DataSourceConfig
	DeleteAddressDS                                  models.DataSourceConfig
	StudentDocumentUploadDS                          models.DataSourceConfig
	StudentDocumentDownloadDS                        models.DataSourceConfig
	GetLocationByPincodeDS                           models.DataSourceConfig
	GetLocationsByPincodeDS                          models.DataSourceConfig
	GetStudentInfoDS                                 models.DataSourceConfig
	GetUserDetailsDS                                 models.DataSourceConfig
	GetUserAdditionalDetailsDS                       models.DataSourceConfig
	StartTestDS                                      models.DataSourceConfig
	GetQuestionPaperInstructionsDS                   models.DataSourceConfig
	GetQuestionPaperDS                               models.DataSourceConfig
	StartTestTimerDS                                 models.DataSourceConfig
	MarkResponseDS                                   models.DataSourceConfig
	PreviewAnswersDS                                 models.DataSourceConfig
	SubmitTestDS                                     models.DataSourceConfig
	StudentResultSummaryWidgetDS                     models.DataSourceConfig
	StudentResultSummaryHeaderWidgetDS               models.DataSourceConfig
	StudentTestScoreCardWidgetDS                     models.DataSourceConfig
	StudentTestInsightCardWidgetDS                   models.DataSourceConfig
	StudentTestResultCardWidgetDS                    models.DataSourceConfig
	ProvisionalTestResultCardsWidgetDS               models.DataSourceConfig
	GetStudentTestInsightsDS                         models.DataSourceConfig
	GetAcademicHomeworkInsightDS                     models.DataSourceConfig
	GetStudentTestOverviewDS                         models.DataSourceConfig
	GetStudentTestReportDS                           models.DataSourceConfig
	ProcessSubjectiveResultsFileDS                   models.DataSourceConfig
	BatchGetTestDS                                   models.DataSourceConfig
	GetStudentPastTestsDS                            models.DataSourceConfig
	GetTestLeaderBoardStudentsDS                     models.DataSourceConfig
	GetTestPastActivitiesDS                          models.DataSourceConfig
	GenerateResultForStudentsDS                      models.DataSourceConfig
	PublishStudentTestResultsDS                      models.DataSourceConfig
	GetAdminTestCorrectionDS                         models.DataSourceConfig
	GetTestSolutionDS                                models.DataSourceConfig
	GetAdminTestSolutionDS                           models.DataSourceConfig
	GetTestQuestionSetDS                             models.DataSourceConfig
	GetQuestionPaperSetsAdminDS                      models.DataSourceConfig
	UpdateAdminTestCorrectionDS                      models.DataSourceConfig
	UpdateTestCorrectionDS                           models.DataSourceConfig
	GetTestInsightDS                                 models.DataSourceConfig
	GetQuestionTestInsightByTestDS                   models.DataSourceConfig
	GetTestSubmissionStatisticsDS                    models.DataSourceConfig
	GetStudentTestInsightsByFiltersDS                models.DataSourceConfig
	GetStudentTestOverviewByFilterDS                 models.DataSourceConfig
	ProcessBulkOrderDs                               models.DataSourceConfig
	GetActivities                                    models.DataSourceConfig
	GetCalendar                                      models.DataSourceConfig
	GetCalendarFreemium                              models.DataSourceConfig
	CourseContentPresentation                        models.DataSourceConfig
	CourseContentRecording                           models.DataSourceConfig
	RecordedContent                                  models.DataSourceConfig
	LiveContentRecording                             models.DataSourceConfig
	EmptyContent                                     models.DataSourceConfig
	EmptyLiveContent                                 models.DataSourceConfig
	EmptyReviseContent                               models.DataSourceConfig
	CourseContentSpecialBooklets                     models.DataSourceConfig
	CourseContentRace                                models.DataSourceConfig
	CourseContentDpp                                 models.DataSourceConfig
	CourseContentStudyModule                         models.DataSourceConfig
	CourseContentConceptCard                         models.DataSourceConfig
	CourseContentQuiz                                models.DataSourceConfig
	CourseContentNcr                                 models.DataSourceConfig
	CourseContentExtraEdge                           models.DataSourceConfig
	SubjectDetails                                   models.DataSourceConfig
	SubjectDetailsV2                                 models.DataSourceConfig
	SubjectDetailsV3                                 models.DataSourceConfig
	CourseDetails                                    models.DataSourceConfig
	CourseDetailsV2                                  models.DataSourceConfig
	SubjectHeader                                    models.DataSourceConfig
	SubjectHeaderV2                                  models.DataSourceConfig
	TopicHeader                                      models.DataSourceConfig
	TopicHeaderV2                                    models.DataSourceConfig
	GetPageDS                                        models.DataSourceConfig
	SubjectBreadcrumb                                models.DataSourceConfig
	SubjectBreadcrumbV2                              models.DataSourceConfig
	GetTestActivities                                models.DataSourceConfig
	GetAdsatResultDS                                 models.DataSourceConfig
	GetImageBannerDS                                 models.DataSourceConfig
	DummyDS                                          models.DataSourceConfig
	JoinMeetingDS                                    models.DataSourceConfig
	EndMeetingDS                                     models.DataSourceConfig
	PreClassStudentDs                                models.DataSourceConfig
	PreClassTeacherDs                                models.DataSourceConfig
	PreClassCommonDs                                 models.DataSourceConfig
	CreateClassMaterialDS                            models.DataSourceConfig
	CompleteUploadForClassMaterialDS                 models.DataSourceConfig
	LeaveMeetingDS                                   models.DataSourceConfig
	StreamingTokenDs                                 models.DataSourceConfig
	CreateClassroomRecordingDS                       models.DataSourceConfig
	CourseStartDateHeaderDS                          models.DataSourceConfig
	UserInfoDS                                       models.DataSourceConfig
	GetUserInfoIDS                                   models.DataSourceConfig
	YourProfileCardWidgetDS                          models.DataSourceConfig
	EditableProfileWidgetDS                          models.DataSourceConfig
	EditablePersonDetailsWidgetDS                    models.DataSourceConfig
	BatchDetailsWidgetDS                             models.DataSourceConfig
	CalDS                                            models.DataSourceConfig
	GetChildrenDS                                    models.DataSourceConfig
	GetRootNodesDS                                   models.DataSourceConfig
	GetTopicsDS                                      models.DataSourceConfig
	GetNodeIDDS                                      models.DataSourceConfig
	GetNodePathDS                                    models.DataSourceConfig
	GetChildrenForNodesDS                            models.DataSourceConfig
	GetNestedChildrenDS                              models.DataSourceConfig
	GetTaxonomyLevelsDS                              models.DataSourceConfig
	UploadTaxonomyDS                                 models.DataSourceConfig
	BorrowTaxonomyDS                                 models.DataSourceConfig
	MapTaxonomyNodesDS                               models.DataSourceConfig
	UpdateTaxonomyDS                                 models.DataSourceConfig
	UpdateTaxonomyNodeDS                             models.DataSourceConfig
	DeactivateTaxonomyNodeDS                         models.DataSourceConfig
	AddNodeInTaxonomyDS                              models.DataSourceConfig
	AddNodeForQBDS                                   models.DataSourceConfig
	ViewMappingsDS                                   models.DataSourceConfig
	UpdateMappingDS                                  models.DataSourceConfig
	DeleteMappingsDS                                 models.DataSourceConfig
	DownloadMappingsDS                               models.DataSourceConfig
	DownloadTaxonomyDS                               models.DataSourceConfig
	CreateTaxonomyBulkDS                             models.DataSourceConfig
	ClearRelatedNodesDS                              models.DataSourceConfig
	GetStaticClassesDS                               models.DataSourceConfig
	CreatePageDS                                     models.DataSourceConfig
	FetchPageDetailsDS                               models.DataSourceConfig
	UpdatePageDS                                     models.DataSourceConfig
	DeletePageDS                                     models.DataSourceConfig
	FilterPageDS                                     models.DataSourceConfig
	MovePageVersionDS                                models.DataSourceConfig
	FetchMaxPageVersionDS                            models.DataSourceConfig
	GetLatestTestCorrectionDS                        models.DataSourceConfig
	CreateWidgetDS                                   models.DataSourceConfig
	FetchWidgetDetailsDS                             models.DataSourceConfig
	UpdateWidgetDS                                   models.DataSourceConfig
	DeleteWidgetDS                                   models.DataSourceConfig
	FilterWidgetDS                                   models.DataSourceConfig
	MoveWidgetVersionDS                              models.DataSourceConfig
	FetchMaxWidgetVersionDS                          models.DataSourceConfig
	CreateWidgetDefinitionDS                         models.DataSourceConfig
	FetchWidgetDefinitionDetailsDS                   models.DataSourceConfig
	UpdateWidgetDefinitionDS                         models.DataSourceConfig
	DeleteWidgetDefinitionDS                         models.DataSourceConfig
	CreateTabDS                                      models.DataSourceConfig
	FetchTabDetailsDS                                models.DataSourceConfig
	UpdateTabDS                                      models.DataSourceConfig
	DeleteTabDS                                      models.DataSourceConfig
	FilterTabDS                                      models.DataSourceConfig
	MoveTabVersionDS                                 models.DataSourceConfig
	FetchMaxTabVersionDS                             models.DataSourceConfig
	LinkWidgetsToPageDS                              models.DataSourceConfig
	LinkTabsToPageDS                                 models.DataSourceConfig
	ListPagesDS                                      models.DataSourceConfig
	ListWidgetsDS                                    models.DataSourceConfig
	ListWidgetDefinitionsDS                          models.DataSourceConfig
	ListTabsDS                                       models.DataSourceConfig
	LinkEntitiesToPageDS                             models.DataSourceConfig
	FetchURLByIDDS                                   models.DataSourceConfig
	AttachURLDS                                      models.DataSourceConfig
	UpdateURLPageMapDS                               models.DataSourceConfig
	DetachURLDS                                      models.DataSourceConfig
	ListURLMappingsDS                                models.DataSourceConfig
	ListURLMappingsV2DS                              models.DataSourceConfig
	EnableURLDS                                      models.DataSourceConfig
	ResolveLMMWigdets                                models.DataSourceConfig
	PreviewPageDS                                    models.DataSourceConfig
	NavOptionsDS                                     models.DataSourceConfig
	PersonaNavOptionsDS                              models.DataSourceConfig
	CourseCreationDS                                 models.DataSourceConfig
	CourseListDS                                     models.DataSourceConfig
	CourseListWithV2SyllabusDS                       models.DataSourceConfig
	CourseFilterDS                                   models.DataSourceConfig
	CourseListingDS                                  models.DataSourceConfig
	CourseUpdateDS                                   models.DataSourceConfig
	ResourceMetaEntitiesDS                           models.DataSourceConfig
	ResourceMappingEntitiesDS                        models.DataSourceConfig
	FacilityListDS                                   models.DataSourceConfig
	FacilityFilterDS                                 models.DataSourceConfig
	CenterBasedAdminFilterListDS                     models.DataSourceConfig
	FacilityListingDS                                models.DataSourceConfig
	FacilityCreationDS                               models.DataSourceConfig
	FacilityUpdateDS                                 models.DataSourceConfig
	FacilityDeleteDS                                 models.DataSourceConfig
	CourseSyllabusCreationDS                         models.DataSourceConfig
	CourseSyllabusCreationV2DS                       models.DataSourceConfig
	CourseSyllabusCreationFromExistingDS             models.DataSourceConfig
	CourseSyllabusCoursesDS                          models.DataSourceConfig
	CourseSyllabusesDS                               models.DataSourceConfig
	CourseSyllabusDS                                 models.DataSourceConfig
	BatchSyllabusDS                                  models.DataSourceConfig
	CourseSyllabusEditDS                             models.DataSourceConfig
	CourseSyllabusValidateNodeDeletionDS             models.DataSourceConfig
	GetCourseSyllabusIDS                             models.DataSourceConfig
	GetCourseSyllabusV2IDS                           models.DataSourceConfig
	PhaseCreationDS                                  models.DataSourceConfig
	PhaseUpdateDS                                    models.DataSourceConfig
	PhaseDeleteDS                                    models.DataSourceConfig
	PhaseListDS                                      models.DataSourceConfig
	PhaseListingDS                                   models.DataSourceConfig
	PhaseFilterDS                                    models.DataSourceConfig
	PhaseDetailDS                                    models.DataSourceConfig
	BatchCreationDS                                  models.DataSourceConfig
	BulkBatchCreateTemplateDownloadDS                models.DataSourceConfig
	BulkBatchCreateUploadDS                          models.DataSourceConfig
	BatchUpdateDS                                    models.DataSourceConfig
	BatchDeleteDS                                    models.DataSourceConfig
	BatchListDS                                      models.DataSourceConfig
	BatchListingDS                                   models.DataSourceConfig
	BatchFilterDS                                    models.DataSourceConfig
	BatchPlansAndSchedulesSummaryDS                  models.DataSourceConfig
	BatchDetailDS                                    models.DataSourceConfig
	BatchBulkGetDS                                   models.DataSourceConfig
	BatchCodePrefixDS                                models.DataSourceConfig
	BatchCodeValidateDS                              models.DataSourceConfig
	BatchStudentListingDS                            models.DataSourceConfig
	CourseContentAdditionDS                          models.DataSourceConfig
	CourseContentRemoveDS                            models.DataSourceConfig
	CourseContentUpdationDS                          models.DataSourceConfig
	CourseContentStaticFilterNamesDS                 models.DataSourceConfig
	CourseContentCopyDS                              models.DataSourceConfig
	CourseContentListDS                              models.DataSourceConfig
	StudentBatchMovementDS                           models.DataSourceConfig
	BatchStudentsDS                                  models.DataSourceConfig
	CourseSyllabusDeleteDS                           models.DataSourceConfig
	VideoPlayerDS                                    models.DataSourceConfig
	LecturePlanWithMetaCreationDS                    models.DataSourceConfig
	LecturePlanAddToMetaDS                           models.DataSourceConfig
	LecturePlanGetAllMetasDS                         models.DataSourceConfig
	LecturePlanGetAllPlansDS                         models.DataSourceConfig
	LecturePlanUpdateAllPlansDS                      models.DataSourceConfig
	LecturePlanUpsertDS                              models.DataSourceConfig
	LecturePlanValidateAllPlansDS                    models.DataSourceConfig
	LecturePlanTemplateDownloadDS                    models.DataSourceConfig
	LecturePlanUploadDS                              models.DataSourceConfig
	LecturePlanTopicsDS                              models.DataSourceConfig
	BatchPlanListingDS                               models.DataSourceConfig
	LecturePlanMetaUpdateDS                          models.DataSourceConfig
	OffersConfig                                     OfferServiceConfig
	LearningMaterialConfig                           LearningMaterialConfig
	LearningJourneyConfig                            LearningJourneyConfig
	QReelsConfig                                     QReelsConfig
	ImprovementBookConfig                            ImprovementBookConfig
	QuestionCollectionConfig                         QuestionCollectionConfig
	ClassScheduleCreateDS                            models.DataSourceConfig
	ClassScheduleCreateAllDS                         models.DataSourceConfig
	ClassScheduleCreateAllV2DS                       models.DataSourceConfig
	ClassScheduleValidateAllDS                       models.DataSourceConfig
	ClassScheduleUpdateDS                            models.DataSourceConfig
	ClassScheduleUpdateAllDS                         models.DataSourceConfig
	ClassScheduleUpdateStatusDS                      models.DataSourceConfig
	ClassScheduleReconcileDS                         models.DataSourceConfig
	ClassScheduleUpsertAllDS                         models.DataSourceConfig
	ClassScheduleDeleteAllDS                         models.DataSourceConfig
	ClassScheduleGetAllDS                            models.DataSourceConfig
	ClassScheduleGetDS                               models.DataSourceConfig
	ClassScheduleUploadDS                            models.DataSourceConfig
	ClassScheduleTemplateDownloadDS                  models.DataSourceConfig
	RegisterDeviceDS                                 models.DataSourceConfig
	DoubtsConfig                                     DoubtsConfig
	DoubtsV2Config                                   DoubtsV2Config
	DoubtsbotConfig                                  DoubtsbotConfig
	DoubtsftueConfig                                 DoubtsftueConfig
	ClassScheduleTeacherViewDS                       models.DataSourceConfig
	ClassScheduleBaseFiltersDS                       models.DataSourceConfig
	ClassScheduleSummaryFiltersDS                    models.DataSourceConfig
	ClassScheduleTeachersDS                          models.DataSourceConfig
	ClassScheduleFacilitiesDS                        models.DataSourceConfig
	DoubtTeacherMappingUploadDS                      models.DataSourceConfig
	DoubtTeacherMappingUploadSpecialBatchDS          models.DataSourceConfig
	DoubtTeacherSpecialBatchDownloadDS               models.DataSourceConfig
	DoubtTeacherMappingDownloadDS                    models.DataSourceConfig
	UserSkillMappingUploadDS                         models.DataSourceConfig
	UserSkillMappingListDS                           models.DataSourceConfig
	UserCoursePageDS                                 models.DataSourceConfig
	CourseSelectWidgetDS                             models.DataSourceConfig
	CourseNavigationWidgetDS                         models.DataSourceConfig
	AddTestUserDS                                    models.DataSourceConfig
	GetTestUserDS                                    models.DataSourceConfig
	EnrollTestUserDS                                 models.DataSourceConfig
	DeleteTestUserDS                                 models.DataSourceConfig
	UserManagementInternalConsoleDS                  models.DataSourceConfig
	ListRolesDS                                      models.DataSourceConfig
	CreateUserPrivilegesDS                           models.DataSourceConfig
	DeleteUserPrivilegesDS                           models.DataSourceConfig
	ListUserIDsByRoleDS                              models.DataSourceConfig
	StudentDashboardDS                               models.DataSourceConfig
	MeetingContentVisibilityDS                       models.DataSourceConfig
	SalesOrchestratorDS                              models.DataSourceConfig
	ScheduleCardDs                                   models.DataSourceConfig
	FreemiumScheduleCardDs                           models.DataSourceConfig
	FreemiumVideoCardDs                              models.DataSourceConfig
	WhyAllenVideoDS                                  models.DataSourceConfig
	ChatServiceDS                                    models.DataSourceConfig
	GetTestRegistrationDS                            models.DataSourceConfig
	SubmitTestRegistrationDS                         models.DataSourceConfig
	ClassroomBcpNotifyDS                             models.DataSourceConfig
	GetTeacherAgendaDS                               models.DataSourceConfig
	CreateTestDS                                     models.DataSourceConfig
	BatchCreateTestsDS                               models.DataSourceConfig
	FilterAdminTestsDs                               models.DataSourceConfig
	DetailedTestsDS                                  models.DataSourceConfig
	GetQuestionPaperAdminPreviewDS                   models.DataSourceConfig
	AttachQuestionPaperDS                            models.DataSourceConfig
	OrderConfig                                      OrderConfig
	CheckoutConfig                                   CheckoutConfig
	UpdateTestStatusDS                               models.DataSourceConfig
	UpdateTestDS                                     models.DataSourceConfig
	UploadTestSolutionDS                             models.DataSourceConfig
	ProvisionalWorkflowAdminDS                       models.DataSourceConfig
	FinalResultWorkflowAdminDS                       models.DataSourceConfig
	DiscoveryUploadUserAttributeDS                   models.DataSourceConfig
	DiscoveryGetUserAttributeDS                      models.DataSourceConfig
	DiscoveryDeleteUserAttributeDS                   models.DataSourceConfig
	DiscoveryGetUserNoticesDS                        models.DataSourceConfig
	ListMasterDataDS                                 models.DataSourceConfig
	SearchAdminTestsDS                               models.DataSourceConfig
	ListOfflineTestResultUploadDS                    models.DataSourceConfig
	TestInsightsWithSupersetDS                       models.DataSourceConfig
	SearchOfflineTestResultFileDS                    models.DataSourceConfig
	UploadOfflineTestResultDS                        models.DataSourceConfig
	UploadOfflineCorrectionTestResultDS              models.DataSourceConfig
	DownloadOfflineTestResultFileDS                  models.DataSourceConfig
	DeleteOfflineResponseFileDS                      models.DataSourceConfig
	UserConflictResolutionConfig                     UserConflictResolutionConfig
	GetS3UploadURLDS                                 models.DataSourceConfig
	GetS3FileURLDS                                   models.DataSourceConfig
	StartBulkProcessingDS                            models.DataSourceConfig
	FilterProcessingRecordDS                         models.DataSourceConfig
	HomeworkServiceConfig                            HomeworkServiceConfig
	GetGroupMentorshipDS                             models.DataSourceConfig
	StudyEssentialsDS                                models.DataSourceConfig
	GetStudentBatchDetails                           models.DataSourceConfig
	SpecialBatchEnrollDS                             models.DataSourceConfig
	SpecialBatchUnenrollDS                           models.DataSourceConfig
	TriggerOMRUploadedWorkflowSignalDS               models.DataSourceConfig
	SyncSubmitTestDS                                 models.DataSourceConfig
	GetTestSyllabusUploadURLDS                       models.DataSourceConfig
	GetFiltersForBatchesDropdownDs                   models.DataSourceConfig
	NoticeBoardConfig                                NoticeBoardConfig
	PDPV2Config                                      PDPV2Config
	GetDynamicURLDS                                  models.DataSourceConfig
	GetMeditationStatusDS                            models.DataSourceConfig
	GetMeditationStatusV2DS                          models.DataSourceConfig
	GetLottieHeaderDS                                models.DataSourceConfig
	CancelTestDS                                     models.DataSourceConfig
	GetStudentRank                                   models.DataSourceConfig
	GetNewGuidanceSessionBannerDS                    models.DataSourceConfig
	GetSolutionAdminPreviewDS                        models.DataSourceConfig
	StudentMergedTestResultWidgetDs                  models.DataSourceConfig
	StudentTestResultInsightsDSV2                    models.DataSourceConfig
	MergeTestsDs                                     models.DataSourceConfig
	FetchMergeAbleTestsDs                            models.DataSourceConfig
	ListTestMinimalDs                                models.DataSourceConfig
	UnMergeTestsDs                                   models.DataSourceConfig
	GetExternalNavigationURLDS                       models.DataSourceConfig
	StudentPaperOneTestResultWidgetDs                models.DataSourceConfig
	StudentPaperTwoTestResultWidgetDs                models.DataSourceConfig
	RankPredictorCsvDs                               models.DataSourceConfig
	AttachSolutionsDs                                models.DataSourceConfig
	MigrateHomeworkDS                                models.DataSourceConfig
	HealPaperDS                                      models.DataSourceConfig
	MigratePaperDS                                   models.DataSourceConfig
	ResumeLearningDS                                 models.DataSourceConfig
	ResumeFlashcardsDS                               models.DataSourceConfig
	GetFlashcardsDataDS                              models.DataSourceConfig
	CreateFlashcardSessionDS                         models.DataSourceConfig
	GetFlashcardsDS                                  models.DataSourceConfig
	GetFlashcardSessionStatsDS                       models.DataSourceConfig
	GetFlashcardsCountDS                             models.DataSourceConfig
	GetFlashcardsSubjectsAndTopicsDS                 models.DataSourceConfig
	GetFlashcardsFtuxDS                              models.DataSourceConfig
	CourseChangeRetryDS                              models.DataSourceConfig
	CDEInsightsConfig                                CDEInsightsConfig
	ValidateQuestionPaperDS                          models.DataSourceConfig
	ExtendTestDS                                     models.DataSourceConfig
	ListingConfig                                    ListingConfig
	CustomTestConfig                                 CustomTestConfig
	ValidateQuestionSetDS                            models.DataSourceConfig
	GetGroupMentorshipForStudentDS                   models.DataSourceConfig
	ColdStartGroupMentorshipDS                       models.DataSourceConfig
	BulkGroupMentorshipRuleCreateUploadDS            models.DataSourceConfig
	ListGroupMentorshipRulesDS                       models.DataSourceConfig
	ReplaceDoubtTeacherMappingsDS                    models.DataSourceConfig
	RemoveDoubtTeacherMappingsDS                     models.DataSourceConfig
	PeakmindGetStudentDetailsDS                      models.DataSourceConfig
	BulkGroupMentorshipRulesCreateTemplateDownloadDS models.DataSourceConfig
	FetchNextPptForConversion                        models.DataSourceConfig
	UpdatePptConversionStatus                        models.DataSourceConfig
	BulkAssignMentorToStudentUploadDS                models.DataSourceConfig
	BulkAssignMentorTemplateDownloadDS               models.DataSourceConfig
	GetTestimonialsDS                                models.DataSourceConfig
	GetNewGuidanceSessionBannerV1DS                  models.DataSourceConfig
	GetSupportWidgetDs                               models.DataSourceConfig
	StudentTestResultWidgetsDs                       models.DataSourceConfig
	CreatePersonalMentorshipScheduleDS               models.DataSourceConfig
	GetTodaysActivityDS                              models.DataSourceConfig
	GetMonthlyActivitiesDS                           models.DataSourceConfig
	GetSurveyDataByIDDS                              models.DataSourceConfig
	SubmitSurveyDataDS                               models.DataSourceConfig
	GetTestPastActivitiesDSV2                        models.DataSourceConfig
	GetTestPastActivitiesDSV3                        models.DataSourceConfig
	AddGroupMentorshipRuleDS                         models.DataSourceConfig
	UpdateGroupMentorshipRuleDS                      models.DataSourceConfig
	DeleteGroupMentorshipRuleDS                      models.DataSourceConfig
	GetMeditationDeeplinkRouteDataDS                 models.DataSourceConfig
	GetMentalHealthCheckDeeplinkRouteDataDS          models.DataSourceConfig
	GetSleepZoneDeeplinkRouteDataDS                  models.DataSourceConfig
	CalServiceConfig                                 CalServiceConfig
	CreateMentorshipBatchExecutionDS                 models.DataSourceConfig
	ListBulkMentorshipCreationExecutionsDS           models.DataSourceConfig
	DownloadMentorshipExecutionDetailsDS             models.DataSourceConfig
	GetCardsStackDS                                  models.DataSourceConfig
	UpdatePersonalMentorshipScheduleDS               models.DataSourceConfig
	GetCounsellorImageBannerDS                       models.DataSourceConfig
	FetchStudentDataDS                               models.DataSourceConfig
	PhaseMergeDS                                     models.DataSourceConfig
	GetMentorshipSessionForStudentDS                 models.DataSourceConfig
	FetchMentorContactDS                             models.DataSourceConfig
	FetchMentorContactForParentDS                    models.DataSourceConfig
	GetUpcomingMentorshipScheduleDS                  models.DataSourceConfig
	FetchPastGroupMentorshipListDS                   models.DataSourceConfig
	GetUpcomingMentorshipScheduleForParentDS         models.DataSourceConfig
	FetchPastGroupMentorshipListForParentDS          models.DataSourceConfig
	GetPolicyDS                                      models.DataSourceConfig
	DeletePolicyDS                                   models.DataSourceConfig
	FilterPolicyDS                                   models.DataSourceConfig
	CreatePolicyDS                                   models.DataSourceConfig
	CreatePolicyVersionDS                            models.DataSourceConfig
	ClonePolicyDS                                    models.DataSourceConfig
	SubmitPolicyReviewDS                             models.DataSourceConfig
	UpdatePolicyDS                                   models.DataSourceConfig
	PolicyMetaEntitiesDS                             models.DataSourceConfig
	FetchMentorStudentNotesDS                        models.DataSourceConfig
	FetchStudentNotesDS                              models.DataSourceConfig
	FetchStudentSliderDataDS                         models.DataSourceConfig
	FetchSchedulesForMentorDS                        models.DataSourceConfig
	FetchPersonalMentorshipScheduleDS                models.DataSourceConfig
	GetFullLottieHeaderDS                            models.DataSourceConfig
	SchedulingServiceConfig                          SchedulingServiceConfig
	BulkCardStatusUpdateDS                           models.DataSourceConfig
	EvaluateStudentEligibilityDS                     models.DataSourceConfig
	SubmitDocumentsDS                                models.DataSourceConfig
	VerificationDocUploadDS                          models.DataSourceConfig
	AssignDefaultMentorToBatchListDS                 models.DataSourceConfig
	AcadOpsBatchListDataDS                           models.DataSourceConfig
	UpdateMentorForBatchDS                           models.DataSourceConfig
	GetPersonalMentorsForBatchDS                     models.DataSourceConfig
	GetStudentOverallPerformanceDataDS               models.DataSourceConfig
	GetStudentSubjectWisePerformanceDataDS           models.DataSourceConfig
	ScoreBoosterConfig                               ScoreBoosterConfig
	OLTSContent                                      models.DataSourceConfig
}

// ClientConfig Internal call config struct
type ClientConfig struct {
	Endpoint  string
	Timeout   time.Duration
	Conn      time.Duration
	Namespace string
}

type CircuitBreakerClientConfig struct {
	FailurePercentageThresholdWithinTimePeriod   uint          // failure percentage threshold before opening the circuit
	FailureMinExecutionThresholdWithinTimePeriod uint          // The number of executions must also exceed the failureExecutionThreshold within the failureThresholdingPeriod
	FailurePeriodThreshold                       time.Duration // failureThresholdingPeriod is the time period in which the failure rate is calculated
	SuccessThreshold                             uint          // number of successive successes before closing the circuit
	Delay                                        time.Duration // delay before retrying after a failure
}

type RetryClientConfig struct {
	MaxRetries int           // Number of retries
	Delay      time.Duration // delay interval between each retry
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// DataSourceConfig Data Source config struct
type GoPool struct {
	MaxConcurrentRoutines uint32
}

type SchedulingServiceConfig struct {
	FetchSchedulesDS        models.DataSourceConfig
	FetchSchedulesSummaryDS models.DataSourceConfig
}

type HomeworkServiceConfig struct {
	CreateHomework                models.DataSourceConfig
	UpdateHomework                models.DataSourceConfig
	HwListingTodayDS              models.DataSourceConfig
	HwListingDueTodayDS           models.DataSourceConfig
	HwListingPendingDS            models.DataSourceConfig
	HwListingCompletedDS          models.DataSourceConfig
	HwListingBacklogDS            models.DataSourceConfig
	GetHomeworkDS                 models.DataSourceConfig
	UpdateHomeworkStatus          models.DataSourceConfig
	GetHomeworkInsights           models.DataSourceConfig
	GetHomeworkOverallInsights    models.DataSourceConfig
	GetHomeworkSubmissionInsights models.DataSourceConfig
	BulkUpsertHomework            models.DataSourceConfig
	ResetHomework                 models.DataSourceConfig
	CreateHomeworks               models.DataSourceConfig
	UpdateHomeworks               models.DataSourceConfig
	HwHomepageListingDS           models.DataSourceConfig
}

type OfferServiceConfig struct {
	GetOfferByID               models.DataSourceConfig
	UpdateOffer                models.DataSourceConfig
	CreateOffer                models.DataSourceConfig
	ActivateOffer              models.DataSourceConfig
	DeactivateOffer            models.DataSourceConfig
	AcknowledgeUpload          models.DataSourceConfig
	PresignedURLOffer          models.DataSourceConfig
	RegisterUserAndCreateOffer models.DataSourceConfig
	DeleteUserFromOffer        models.DataSourceConfig
	FilterOffers               models.DataSourceConfig
	FilterOffersWithoutAuth    models.DataSourceConfig
	FilterUploadRecord         models.DataSourceConfig
	GetOfferFileURL            models.DataSourceConfig
}

// ServerConfig Server config struct
type ServerConfig struct {
	Port                        string
	App                         App
	ReadTimeout                 time.Duration
	WriteTimeout                time.Duration
	JwtSecret                   string
	JwtSecretLocation           string
	AesEncryptionKey            string
	AesSecretIV                 string
	AesEncryptionSecretLocation string
}

type LearningMaterialConfig struct {
	CreateLearningMaterialBulkDS         models.DataSourceConfig
	InitiateMultipartUploadForMaterialDS models.DataSourceConfig
	PresignedURLToUploadMaterialPartDS   models.DataSourceConfig
	CompleteMultipartUploadForMaterialDS models.DataSourceConfig
	GetPresignedURLForMaterialDS         models.DataSourceConfig
	GetMaterialByIDDS                    models.DataSourceConfig
	ValidateBulkCreateRequestDS          models.DataSourceConfig
	FilterMaterialsDS                    models.DataSourceConfig
	FilterMaterialsDSV2                  models.DataSourceConfig
	GetMaterialByMaterialIDDS            models.DataSourceConfig
	UpdateMaterialDS                     models.DataSourceConfig
	ActivateMaterialDS                   models.DataSourceConfig
	DeactivateMaterialDS                 models.DataSourceConfig
	TypeSubTypeMap                       map[string][]string
	CreateThumbnailBulkDS                models.DataSourceConfig
	GetThumbnailByIDDS                   models.DataSourceConfig
	FilterThumbnailsDS                   models.DataSourceConfig
	UpdateThumbnailDS                    models.DataSourceConfig
	ThumbnailInitUploadDS                models.DataSourceConfig
	CompleteThumbnailUploadDS            models.DataSourceConfig
	DownloadThumbnailURLDS               models.DataSourceConfig
	ValidateBulkThumbnailCreateRequestDS models.DataSourceConfig
	GetThumbnailURLByTaxonomyDS          models.DataSourceConfig
	IndexMaterialToLmmSearchDS           models.DataSourceConfig
	DefaultThumbnailURL                  string
	IsAbacEnabled                        bool
	IsLmmRedisEnabled                    bool
	ThumbnailFetchTimeout                int32 // ms
	GetTypeSubTypeData                   models.DataSourceConfig
	GetLanguages                         models.DataSourceConfig
	GetLearningCategories                models.DataSourceConfig
	GetPodcastVideoWidget                models.DataSourceConfig
	InitNACMultipartUploadDS             models.DataSourceConfig
	GetPresignedURLToUploadPartDS        models.DataSourceConfig
	CompleteMultiPartUploadDS            models.DataSourceConfig
	GetNACPresignedURLDS                 models.DataSourceConfig
	CreateBulkNACDS                      models.DataSourceConfig
	GetNACByIDDS                         models.DataSourceConfig
	UpdateNACDS                          models.DataSourceConfig
	FilterNACDS                          models.DataSourceConfig
	IndexNACDS                           models.DataSourceConfig
	GetPodcastCarouselWidget             models.DataSourceConfig
	UploadChaptersDS                     models.DataSourceConfig
}

type CalServiceConfig struct {
	DownloadLearningMaterialDS models.DataSourceConfig
	GetDownloadPresignedURLDS  models.DataSourceConfig
	GetBulkDownloadMaterialDS  models.DataSourceConfig
	GetLearningMaterialDS      models.DataSourceConfig
	GetStaticHighlightsInfoDS  models.DataSourceConfig
}

type LearningJourneyConfig struct {
	GetGoalTemplateDS             models.DataSourceConfig
	GetStartedScreenFreemiumDS    models.DataSourceConfig
	UpdateBookmarkDS              models.DataSourceConfig
	JourneyUpdateBookmarkDS       models.DataSourceConfig
	GetContentActionDS            models.DataSourceConfig
	CreateUserGoalDS              models.DataSourceConfig
	GetStartedWidgetDS            models.DataSourceConfig
	GetNextUserContentDS          models.DataSourceConfig
	RecordUserActionsDS           models.DataSourceConfig
	SuggestNextUserContentDS      models.DataSourceConfig
	CloseAndInitiateUserStepDS    models.DataSourceConfig
	GetUserProgressDS             models.DataSourceConfig
	GetFremiumProgressWidgetDS    models.DataSourceConfig
	GetProgressSummaryWidgetDS    models.DataSourceConfig
	ListTopicFreemiumDS           models.DataSourceConfig
	GetSubjectTabWidgetDS         models.DataSourceConfig
	ListBookmarksWidgetDS         models.DataSourceConfig
	GetBookmarksDS                models.DataSourceConfig
	DownloadJourneyTemplateDS     models.DataSourceConfig
	BulkCreateLearningJourneysDS  models.DataSourceConfig
	DeleteUserGoalByFilterDS      models.DataSourceConfig
	RefreshJourneyDS              models.DataSourceConfig
	GetLearningJourneyStepsListDS models.DataSourceConfig
	GetLearningJourneyStepDS      models.DataSourceConfig
	UpdateLearningJourneyStepsDS  models.DataSourceConfig
	UnlockAllUserStepsDS          models.DataSourceConfig
}

type QReelsConfig struct {
	GetContentDS           models.DataSourceConfig
	RecordActionsDS        models.DataSourceConfig
	QreelsSuggestContentDS models.DataSourceConfig
	GetQuestionsByIDsDS    models.DataSourceConfig
	QreelsBookmarkDS       models.DataSourceConfig
}

type ImprovementBookConfig struct {
	GetIBSummaryDS                models.DataSourceConfig
	GetIBSubjectTopicsDS          models.DataSourceConfig
	GetIBQuestionsByStatusDS      models.DataSourceConfig
	GetIBStatusAggregateDetailsDS models.DataSourceConfig
	IBRecordAnsDS                 models.DataSourceConfig
	IBReattemptMarkResponseDS     models.DataSourceConfig
	GetIBQuestionWidgetDS         models.DataSourceConfig
	GetSolutionDS                 models.DataSourceConfig
	IBReattemptResultDS           models.DataSourceConfig
	GetIBQuestionsForReattemptDS  models.DataSourceConfig
	GetIBLatestByStudentIDDS      models.DataSourceConfig
	GetIBFtuxDS                   models.DataSourceConfig
}

type QuestionCollectionConfig struct {
	CopyQuestionCollectionDS                models.DataSourceConfig
	GetRelevantQuestionCollectionDS         models.DataSourceConfig
	GetQuestionCollectionPreviewDS          models.DataSourceConfig
	GetQuestionCollectionMergedViewDS       models.DataSourceConfig
	UpdateQuestionCollectionSectionsDS      models.DataSourceConfig
	UpdateQuestionCollectionQuestionsDS     models.DataSourceConfig
	UpdateQuestionCollectionStatusDS        models.DataSourceConfig
	GetQuestionCollectionStatsDS            models.DataSourceConfig
	GetTopicsDS                             models.DataSourceConfig
	GetMeetingsByQuestionCollectionFilterDS models.DataSourceConfig
	BatchGetQuestionCollectionStatsDS       models.DataSourceConfig
}

type CustomTestConfig struct {
	GetPracticePageDS               models.DataSourceConfig
	GetPracticeWidgetDS             models.DataSourceConfig
	GetPracticeDialogueDS           models.DataSourceConfig
	CompletedTestsDS                models.DataSourceConfig
	GetFilterSelectionWidgetDS      models.DataSourceConfig
	CreateCustomTestDS              models.DataSourceConfig
	FetchCustomTestQuestionsDS      models.DataSourceConfig
	GetResultPageDS                 models.DataSourceConfig
	RecordTestActionsDS             models.DataSourceConfig
	CustomTestSubmitDS              models.DataSourceConfig
	GetPersonalizedPracticeWidgetDS models.DataSourceConfig
}

type DoubtsConfig struct {
	DoubtsHomepageDS            models.DataSourceConfig
	DoubtsBannerDS              models.DataSourceConfig
	DoubtsOpenDS                models.DataSourceConfig
	DoubtsResolvedDS            models.DataSourceConfig
	DoubtsEmptyDS               models.DataSourceConfig
	GetDoubtDS                  models.DataSourceConfig
	CreateDoubtDS               models.DataSourceConfig
	CreateDoubtReplyDS          models.DataSourceConfig
	GetFileURLDS                models.DataSourceConfig
	GetPresignedURLForUploadDS  models.DataSourceConfig
	GetAllDoubtsForSeekerIDDS   models.DataSourceConfig
	GetAllDoubtsForResolverIDDS models.DataSourceConfig
	GetSubjectTopicDS           models.DataSourceConfig
	GetRepliesOnDoubtDS         models.DataSourceConfig
	GetResolverFiltersDS        models.DataSourceConfig
	GetSeekerFiltersDS          models.DataSourceConfig
	MarkDoubtStatusDS           models.DataSourceConfig
	GetDoubtTeacherMappingDS    models.DataSourceConfig
	MarkDoubtAsSupportDS        models.DataSourceConfig
}

type DoubtsV2Config struct {
	DoubtsHomepageV2DS               models.DataSourceConfig
	DoubtsWebHomepageV2DS            models.DataSourceConfig
	DoubtsBannerV2DS                 models.DataSourceConfig
	DoubtsOpenV2DS                   models.DataSourceConfig
	DoubtsResolvedV2DS               models.DataSourceConfig
	DoubtsEmptyV2DS                  models.DataSourceConfig
	DoubtsTypescreenV2DS             models.DataSourceConfig
	DoubtsWebTypescreenV2DS          models.DataSourceConfig
	DoubtsSupportV2DS                models.DataSourceConfig
	DoubtsPreviewV2DS                models.DataSourceConfig
	GetDoubtV2DS                     models.DataSourceConfig
	CreateDoubtV2DS                  models.DataSourceConfig
	CreateDraftDoubtV2DS             models.DataSourceConfig
	CreateDoubtReplyV2DS             models.DataSourceConfig
	GetFileURLV2DS                   models.DataSourceConfig
	GetPresignedURLForUploadV2DS     models.DataSourceConfig
	GetAllDoubtForSeekerIdV2DS       models.DataSourceConfig
	GetAllDoubtForResolverIdV2DS     models.DataSourceConfig
	GetSubjectTopicV2DS              models.DataSourceConfig
	GetRepliesOnDoubtV2DS            models.DataSourceConfig
	GetResolverFiltersV2DS           models.DataSourceConfig
	GetSeekerFiltersV2DS             models.DataSourceConfig
	MarkDoubtStatusV2DS              models.DataSourceConfig
	GetDoubtTeacherMappingV2DS       models.DataSourceConfig
	MarkDoubtAsSupportV2DS           models.DataSourceConfig
	SupportPopUpDS                   models.DataSourceConfig
	DoubtsSubmitV2DS                 models.DataSourceConfig
	DoubtsIllustrationV2DS           models.DataSourceConfig
	DoubtsTransferV1DS               models.DataSourceConfig
	DoubtsHomepageV3DS               models.DataSourceConfig
	DoubtsBannerV3DS                 models.DataSourceConfig
	DoubtsTypescreenV3DS             models.DataSourceConfig
	DoubtsSubmitV3DS                 models.DataSourceConfig
	GetRepliesOnDoubtV3DS            models.DataSourceConfig
	DoubtsReassignFilterV3DS         models.DataSourceConfig
	DoubtsReassignSpecialBatchDS     models.DataSourceConfig
	DoubtsReassignV3DS               models.DataSourceConfig
	DoubtsListDS                     models.DataSourceConfig
	TransferAllDoubtsOfResolverDS    models.DataSourceConfig
	DoubtsPaginatedRepliesDS         models.DataSourceConfig
	DoubtsPostTimestampRepliesListDS models.DataSourceConfig
	AskDoubtFabDS                    models.DataSourceConfig
	AskDoubtFabV2DS                  models.DataSourceConfig
}

type DoubtsbotConfig struct {
	GetRepliesOnDoubtV4DS models.DataSourceConfig
	CreateDoubtReplyV4DS  models.DataSourceConfig
	DoubtsOpenV4DS        models.DataSourceConfig
	DoubtsResolvedV4DS    models.DataSourceConfig
	DoubtsListV4DS        models.DataSourceConfig
	DoubtsSubmitV4DS      models.DataSourceConfig
	DoubtsTypescreenV4DS  models.DataSourceConfig
	DoubtsBannerV4DS      models.DataSourceConfig
	DoubtsPreviewV4DS     models.DataSourceConfig
	BotIntroBannerDS      models.DataSourceConfig
	BotFteuDS             models.DataSourceConfig
	DoubtsSubmitV5DS      models.DataSourceConfig
	DoubtsPreviewV5DS     models.DataSourceConfig
	DoubtsSubmitV6DS      models.DataSourceConfig
	DoubtsListV5DS        models.DataSourceConfig
}

type DoubtsftueConfig struct {
	DoubtsBotFtueLoaderDS      models.DataSourceConfig
	CreateSampleDoubtBotFtueDS models.DataSourceConfig
	GetAllDoubtForSeekerIdV5DS models.DataSourceConfig
	DoubtsBannerV5DS           models.DataSourceConfig
	DoubtsBannerV6DS           models.DataSourceConfig
	GetSeekerFiltersV3DS       models.DataSourceConfig
	AskDoubtVtagDS             models.DataSourceConfig
	GetAllDoubtForSeekerIdV6DS models.DataSourceConfig
	AskDoubtVtagV2DS           models.DataSourceConfig
	CreateDraftDoubtV3DS       models.DataSourceConfig
}

type CheckoutConfig struct {
	CreateCart                 models.DataSourceConfig
	UpdatePurchaserDetails     models.DataSourceConfig
	GetCart                    models.DataSourceConfig
	ReviewCart                 models.DataSourceConfig
	ApplyCoupons               models.DataSourceConfig
	RemoveCoupons              models.DataSourceConfig
	CreateCheckout             models.DataSourceConfig
	ReInitiateCheckout         models.DataSourceConfig
	GetOrderStatus             models.DataSourceConfig
	JuspayPaymentStatus        models.DataSourceConfig
	AppstorePaymentStatus      models.DataSourceConfig
	GeneratePaymentLinkDS      models.DataSourceConfig
	GenerateRetryPaymentLinkDS models.DataSourceConfig
	CreateCartV2               models.DataSourceConfig
	GetCartV2                  models.DataSourceConfig
	UpdateCartV2               models.DataSourceConfig
}

type OrderConfig struct {
	GetOrderDs                models.DataSourceConfig
	GetInvoicesDs             models.DataSourceConfig
	ListOrdersDs              models.DataSourceConfig
	GetBulkOrderUploadDs      models.DataSourceConfig
	ProcessBulkOrderDs        models.DataSourceConfig
	ListOrderByStatusDS       models.DataSourceConfig
	ListOrderByStatusForMobDS models.DataSourceConfig
	GetDownLoadURLDS          models.DataSourceConfig
}
type NoticeBoardConfig struct {
	CreateNoticeDS   models.DataSourceConfig
	GetNoticeDS      models.DataSourceConfig
	GetNoticesDS     models.DataSourceConfig
	UpdateNoticeDS   models.DataSourceConfig
	SendNoticeDS     models.DataSourceConfig
	DeleteNoticeDS   models.DataSourceConfig
	MediaURLDS       models.DataSourceConfig
	ReadReceiptDS    models.DataSourceConfig
	TeacherViewDS    models.DataSourceConfig
	GetUnreadCountDS models.DataSourceConfig
}

type CDEInsightsConfig struct {
	AttendanceInsightsDS models.DataSourceConfig
}

type PlaylistConfig struct {
	HLSFilename  string
	DASHFilename string
}

type WhiteListSubTypesForContentAuth struct {
	SubTypes []string
}

type AkamaiConfig struct {
	BaseURL string
	Key     string
}

type UserConflictResolutionConfig struct {
	GetUserConflictByID models.DataSourceConfig
	UpdateUserConflict  models.DataSourceConfig
	FilterUserConflicts models.DataSourceConfig
}

type PDPV2Config struct {
	ProductInfoDS           models.DataSourceConfig
	PaymentFooterMobileDS   models.DataSourceConfig
	DiscoveryListingGroupDS models.DataSourceConfig
}

type ListingConfig struct {
	GetListing            models.DataSourceConfig
	GetListingWithFilter  models.DataSourceConfig
	CreateListing         models.DataSourceConfig
	ActivateListing       models.DataSourceConfig
	DeactivateListing     models.DataSourceConfig
	UpdateListing         models.DataSourceConfig
	BrowseListings        models.DataSourceConfig
	ListingHighlights     models.DataSourceConfig
	ListingFaqTnc         models.DataSourceConfig
	ListingDetails        models.DataSourceConfig
	ListingInfo           models.DataSourceConfig
	EnrollDs              models.DataSourceConfig
	ListingPaymentDetails models.DataSourceConfig
	GetAvailableCourse    models.DataSourceConfig
	GetGroupAttributeList models.DataSourceConfig
	CourseGrid            models.DataSourceConfig
}

type ScoreBoosterConfig struct {
	GetGenericEntryPointDS models.DataSourceConfig
}

type App struct {
	Name    string
	Version string
}

type AppConfig struct {
	AppID           string
	ConfigName      string
	PollingInterval int64
}

const (
	CircuitBreakerFailurePercentageThresholdSuffix       = ".cb.failure_percentage_threshold"
	CircuitBreakerMinExecutionThresholdSuffix            = ".cb.min_execution_threshold"
	CircuitBreakerFailurePeriodThresholdSuffix           = ".cb.failure_period_threshold"
	CircuitBreakerDelaySuffix                            = ".cb.delay"
	CircuitBreakerSuccessThresholdSuffix                 = ".cb.success_threshold"
	DefaultCircuitBreakerPercentageThreshold             = 50
	DefaultCircuitMinExecutionThreshold                  = 100
	DefaultCircuitBreakerFailurePeriodThresholdInSeconds = 60
	DefaultCircuitBreakerDelayMs                         = 30000
	DefaultCircuitBreakerSuccessThreshold                = 5
)

const (
	RetryMaxRetriesSuffix  = ".retry.max_retries"
	RetryDelaySuffix       = ".retry.delay"
	DefaultRetryMaxRetries = 3
	DefaultRetryDelayMs    = 1000
)

const (
	Base10    = 10
	BitSize32 = 32
)

const (
	TimeoutSuffix     = ".timeout"
	EndPointSuffix    = ".endpoint"
	ConnTimeoutSuffix = ".conn"
	NamespaceSuffix   = ".namespace"
	TimeInMs          = "ms"
	TimeInSeconds     = "s"
	FileType          = "properties"
	LocalConfigName   = "client"
	JwtSecretKey      = "jwt_secret"
)

const (
	ReadConfigErrorLog      = "error in reading server properties file config client: %v ,Error: %v"
	StringToIntParsingError = "error parsing %s for client: %v ,Error: %v, using defaults"
	DynConfigParsingError   = "error fetching %s aws config for client: %v ,Error: %v"
)

func GetCircuitBreakerClientConfigs(client string, cnf *Config) CircuitBreakerClientConfig {
	if cnf.DynamicConfig == nil {
		return readCircuitBreakerConfigFromLocalConfig(client)
	}

	return readCircuitBreakerConfigFromDynConfig(client, cnf)
}

func getConfigDirectory() string {
	var dir string

	env := os.Getenv("ENV")
	log.Infof("Found Env as %v", env)

	if env != "" {
		dir = "/data/conf/" + env + "/"
		log.Infof("Found dir as %v", dir)
	} else {
		log.Infof("failed to get build environment, using local configs")

		dir = "./config/local/"
	}

	return dir
}

func readCircuitBreakerConfigFromLocalConfig(client string) CircuitBreakerClientConfig {
	v := viper.New()
	v.SetConfigType(FileType)
	v.AddConfigPath(getConfigDirectory())
	v.SetConfigName(LocalConfigName)

	if err := v.ReadInConfig(); err != nil {
		log.Errorf(ReadConfigErrorLog, client, err.Error())
	}

	failurePercentageThresholdConfig := v.GetString(join(client, CircuitBreakerFailurePercentageThresholdSuffix))
	failurePercentageThresholdConfigInt, err := strconv.ParseUint(failurePercentageThresholdConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "failurePercentageThresholdConfig", client, err)

		failurePercentageThresholdConfigInt = DefaultCircuitBreakerPercentageThreshold
	}

	failureMinExecutionThresholdConfig := v.GetString(join(client, CircuitBreakerMinExecutionThresholdSuffix))
	failureMinExecutionThresholdConfigInt, err := strconv.ParseUint(failureMinExecutionThresholdConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "failureMinExecutionThresholdConfig", client, err)

		failureMinExecutionThresholdConfigInt = DefaultCircuitMinExecutionThreshold
	}

	failurePeriodThresholdConfig := v.GetString(join(client, CircuitBreakerFailurePeriodThresholdSuffix))
	failurePeriodThresholdConfigInDuration, err := time.ParseDuration(join(failurePeriodThresholdConfig, TimeInSeconds))

	if failurePeriodThresholdConfigInDuration == 0 || err != nil {
		log.Warnf(StringToIntParsingError, "failurePeriodThresholdConfig", client, err)

		failurePeriodThresholdConfigInDuration = DefaultCircuitBreakerFailurePeriodThresholdInSeconds * time.Second
	}

	successThresholdConfig := v.GetString(join(client, CircuitBreakerSuccessThresholdSuffix))
	successThresholdConfigInt, err := strconv.ParseUint(successThresholdConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "successThresholdConfig", client, err)

		successThresholdConfigInt = DefaultCircuitBreakerSuccessThreshold
	}

	delayConfig := v.GetString(join(client, CircuitBreakerDelaySuffix))
	delayConfigInDuration, err := time.ParseDuration(join(delayConfig, TimeInMs))

	if delayConfigInDuration == 0 || err != nil {
		delayConfigInDuration = DefaultCircuitBreakerDelayMs * time.Millisecond
	}

	return CircuitBreakerClientConfig{
		FailurePercentageThresholdWithinTimePeriod:   uint(failurePercentageThresholdConfigInt),
		FailureMinExecutionThresholdWithinTimePeriod: uint(failureMinExecutionThresholdConfigInt),
		FailurePeriodThreshold:                       failurePeriodThresholdConfigInDuration,
		SuccessThreshold:                             uint(successThresholdConfigInt),
		Delay:                                        delayConfigInDuration,
	}
}

func readCircuitBreakerConfigFromDynConfig(client string, cnf *Config) CircuitBreakerClientConfig {
	// Reading from aws App Config
	failurePercentageThresholdConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerFailurePercentageThresholdSuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "failurePercentageThresholdConfig", client, err)
	}

	failurePercentageThresholdConfigInt, err := strconv.ParseUint(failurePercentageThresholdConfig, 10, 32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "failurePercentageThresholdConfig", client, err)

		failurePercentageThresholdConfigInt = DefaultCircuitBreakerPercentageThreshold
	}

	failureMinExecutionThresholdConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerMinExecutionThresholdSuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "failureMinExecutionThresholdConfig", client, err)
	}

	failureMinExecutionThresholdConfigInt, err := strconv.ParseUint(failureMinExecutionThresholdConfig, 10, 32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "failureMinExecutionThresholdConfig", client, err)

		failureMinExecutionThresholdConfigInt = DefaultCircuitMinExecutionThreshold
	}

	failurePeriodThresholdConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerFailurePeriodThresholdSuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "failurePeriodThresholdConfig", client, err)
	}

	failurePeriodThresholdConfigInDuration, err := time.ParseDuration(join(failurePeriodThresholdConfig, TimeInSeconds))
	if failurePeriodThresholdConfigInDuration == 0 || err != nil {
		log.Warnf(StringToIntParsingError, "failurePeriodThresholdConfig", client, err)

		failurePeriodThresholdConfigInDuration = DefaultCircuitBreakerFailurePeriodThresholdInSeconds * time.Second
	}

	successThresholdConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerSuccessThresholdSuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "successThresholdConfig", client, err)
	}

	successThresholdConfigInt, err := strconv.ParseUint(successThresholdConfig, 10, 32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "successThresholdConfig", client, err)

		successThresholdConfigInt = DefaultCircuitBreakerSuccessThreshold
	}

	delayConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerDelaySuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "successThresholdConfig", client, err)
	}

	delayConfigInDuration, err := time.ParseDuration(join(delayConfig, TimeInMs))
	if delayConfigInDuration == 0 || err != nil {
		log.Warnf(StringToIntParsingError, "delayConfig", client, err)

		delayConfigInDuration = DefaultCircuitBreakerDelayMs * time.Millisecond
	}

	return CircuitBreakerClientConfig{
		FailurePercentageThresholdWithinTimePeriod:   uint(failurePercentageThresholdConfigInt),
		FailureMinExecutionThresholdWithinTimePeriod: uint(failureMinExecutionThresholdConfigInt),
		FailurePeriodThreshold:                       failurePeriodThresholdConfigInDuration,
		SuccessThreshold:                             uint(successThresholdConfigInt),
		Delay:                                        delayConfigInDuration,
	}
}

func GetRetryClientConfigs(client string, cnf *Config) RetryClientConfig {
	if cnf.DynamicConfig == nil {
		return readRetryConfigFromLocalConfig(client)
	}

	return readRetryFromDynConfig(client, cnf)
}

func readRetryConfigFromLocalConfig(client string) RetryClientConfig {
	v := viper.New()
	v.SetConfigType(FileType)
	v.AddConfigPath(getConfigDirectory())
	v.SetConfigName(LocalConfigName)

	if err := v.ReadInConfig(); err != nil {
		log.Errorf(ReadConfigErrorLog, client, err.Error())
	}

	maxRetriesConfig := v.GetString(join(client, RetryMaxRetriesSuffix))
	maxRetriesConfigInt, err := strconv.ParseInt(maxRetriesConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf("error parsing maxRetriesConfig for client: %v ,Error: %v, using defaults", client, err.Error())

		maxRetriesConfigInt = DefaultRetryMaxRetries
	}

	delayConfig := v.GetString(join(client, RetryDelaySuffix))
	delayConfigInDuration, err := time.ParseDuration(join(delayConfig, TimeInMs))

	if delayConfigInDuration == 0 || err != nil {
		log.Warnf("error parsing retry delayConfig for client: %v ,Error: %v, using defaults", client, err)

		delayConfigInDuration = DefaultRetryDelayMs * time.Millisecond
	}

	return RetryClientConfig{
		MaxRetries: int(maxRetriesConfigInt),
		Delay:      delayConfigInDuration,
	}
}

func readRetryFromDynConfig(client string, cnf *Config) RetryClientConfig {
	// Reading from aws App Config
	maxRetriesConfig, err := cnf.DynamicConfig.Get(client + RetryMaxRetriesSuffix)
	if err != nil {
		log.Warnf("error fetching failureThresholdConfig aws config for client: %v ,Error: %v", client, err.Error())
	}

	maxRetriesConfigInt, err := strconv.ParseInt(maxRetriesConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf("error parsing maxRetriesConfig for client: %v ,Error: %v, using defaults", client, err.Error())

		maxRetriesConfigInt = DefaultRetryMaxRetries
	}

	delayConfig, err := cnf.DynamicConfig.Get(client + RetryDelaySuffix)
	if err != nil {
		log.Warnf("error fetching retry delayConfig aws config for client: %v ,Error: %v", client, err.Error())
	}

	delayConfigInDuration, err := time.ParseDuration(join(delayConfig, TimeInMs))
	if delayConfigInDuration == 0 || err != nil {
		log.Warnf("error parsing retry delayConfig for client: %v ,Error: %v, using defaults", client, err)

		delayConfigInDuration = DefaultRetryDelayMs * time.Millisecond
	}

	return RetryClientConfig{
		MaxRetries: int(maxRetriesConfigInt),
		Delay:      delayConfigInDuration,
	}
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}

	return sb.String()
}

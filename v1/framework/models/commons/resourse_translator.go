package commons

import (
	cpbte "github.com/Allen-Career-Institute/common-protos/classroom/v1/types"
	"github.com/Allen-Career-Institute/common-protos/resource/v1/types/enums"
)

func GetModeDisplay(mode string) string {
	return GetEnumDisplayName(enums.Mode(enums.Mode_value[mode]))
}

func GetCourseModuleModeDisplay(mode string) string {
	return GetEnumDisplayName(enums.CourseModuleMode(enums.CourseModuleMode_value[mode]))
}

func GetLangDisplay(language string) string {
	return GetEnumDisplayName(enums.Language(enums.Language_value[language]))
}

func GetClassDisplay(class string) string {
	return GetEnumDisplayName(enums.Class(enums.Class_value[class]))
}

func GetMeetingTypeDisplayName(meetingType string) string {
	return GetEnumExtension[string](cpbte.MeetingType(cpbte.MeetingType_value[meetingType]), cpbte.E_DisplayName)
}

func GetMeetingTypeCode(meetingType string) string {
	return GetEnumExtension[string](cpbte.MeetingType(cpbte.MeetingType_value[meetingType]), cpbte.E_Code)
}

func GetScheduleModeDisplayName(scheduleMode string) string {
	return GetEnumExtension[string](enums.ScheduleMode(enums.ScheduleMode_value[scheduleMode]), cpbte.E_DisplayName)
}

func GetScheduleModeCode(scheduleMode string) string {
	return GetEnumExtension[string](enums.ScheduleMode(enums.ScheduleMode_value[scheduleMode]), cpbte.E_Code)
}

func GetScheduleStatusDisplayName(scheduleStatus string) string {
	return GetEnumExtension[string](enums.ScheduleStatus(enums.ScheduleStatus_value[scheduleStatus]), cpbte.E_DisplayName)
}

func GetScheduleStatusCode(scheduleStatus string) string {
	return GetEnumExtension[string](enums.ScheduleStatus(enums.ScheduleStatus_value[scheduleStatus]), cpbte.E_Code)
}

func GetOfferingDisplay(offering string) string {
	return GetEnumDisplayName(enums.CourseOffering(enums.CourseOffering_value[offering]))
}

func GetFeatureDisplay(feature string) string {
	return GetEnumDisplayName(enums.CourseSpecial(enums.CourseSpecial_value[feature]))
}

func GetStreamDisplay(stream string) string {
	return GetEnumDisplayName(enums.Stream(enums.Stream_value[stream]))
}

func GetMasterCourseDisplay(masterCourse string) string {
	return GetEnumDisplayName(enums.MasterCourse(enums.MasterCourse_value[masterCourse]))
}

func GetBoardDisplay(board string) string {
	return GetEnumDisplayName(enums.Board(enums.Board_value[board]))
}

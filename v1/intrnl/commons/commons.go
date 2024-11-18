package commons

import (
	"google.golang.org/protobuf/types/known/structpb"
)

type Icon struct {
	IconURI    string `json:"icon_uri"`
	ButtonType string `json:"button_type"`
	Action     Action `json:"action"`
}
type Caption struct {
	Text      string `json:"text"`
	TextColor Color  `json:"text_color"`
}

type VideoData struct {
	URL           string       `json:"url"`
	ChaptersURL   string       `json:"chapters_url"`
	ContentID     string       `json:"content_id"`
	Watermark     string       `json:"watermark"`
	Configuration *VideoConfig `json:"configuration,omitempty"`
}

type VideoDataV2 struct {
	URL           string       `json:"url"`
	ContentID     string       `json:"content_id"`
	Watermark     string       `json:"watermark"`
	Configuration *VideoConfig `json:"configuration"`
	ImageData     Image        `json:"image_data"`
}

type AesEncryptionSecrets struct {
	EncryptionSecretKey string `json:"encryption_secret_key"`
	EncryptionIv        string `json:"encryption_iv"`
}

type VideoConfig struct {
	Muted    bool `json:"muted"`
	AutoPlay bool `json:"auto_play"`
}

type Action struct {
	Data           Data             `json:"data"`
	Type           string           `json:"type"`
	TrackingParams *structpb.Struct `json:"tracking_params,omitempty"`
}

type ExternalNavigationURL struct {
	URL string `json:"url"`
}

type Data struct {
	URI     string         `json:"uri"`
	Name    string         `json:"name"`
	Query   map[string]any `json:"query"`
	Type    string         `json:"type"`
	State   string         `json:"state,omitempty"`
	Scale   string         `json:"scale,omitempty"`
	Title   string         `json:"title,omitempty"`
	Action  string         `json:"action,omitempty"`
	Method  string         `json:"method,omitempty"`
	Phone   string         `json:"phone"`
	Email   string         `json:"email"`
	Payload interface{}    `json:"payload,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
	Content interface{}    `json:"content,omitempty"`
}

type GenericSupportData struct {
	ImageData    Image  `json:"image_data"`
	PrimaryCTA   CTA    `json:"primary_cta"`
	SecondaryCTA *CTA   `json:"secondary_cta"`
	Subtitle     string `json:"subtitle"`
	Title        string `json:"title"`
}

type Image struct {
	URL     string `json:"url"`
	Name    string `json:"name,omitempty"`
	AltText string `json:"alt_text,omitempty"`
	ID      string `json:"id,omitempty"`
	Color   string `json:"color,omitempty"`
}

type CTA struct {
	Label     string      `json:"label,omitempty"`
	Icon      string      `json:"icon_uri,omitempty"`
	ImageData Image       `json:"image_data,omitempty"`
	Action    interface{} `json:"action,omitempty"`
}

type NavigationWithForwardAction struct {
	Type          string         `json:"type,omitempty"`
	Data          NavigationData `json:"data,omitempty"`
	ForwardAction *Action        `json:"forward_action,omitempty"`
}

type Navigation struct {
	Type string         `json:"type"`
	Data NavigationData `json:"data"`
}

type DropDown struct {
	Label   string   `json:"label,omitempty"`
	Icon    string   `json:"icon_uri,omitempty"`
	Options []Option `json:"options,omitempty"`
}

type Option struct {
	Id     string      `json:"id,omitempty"`
	Label  string      `json:"label,omitempty"`
	Action interface{} `json:"action,omitempty"`
}

type ExternalNavigation struct {
	Type           string                 `json:"type"`
	Data           ExternalNavigationData `json:"data"`
	TrackingParams *structpb.Struct       `json:"tracking_params,omitempty"`
}

type NavigationData struct {
	URI               string                 `json:"uri"`
	Type              string                 `json:"type,omitempty"`
	Query             map[string]interface{} `json:"query,omitempty"`
	PollCurrentScreen bool                   `json:"poll_current_screen,omitempty"`
	Name              string                 `json:"name,omitempty"`
	Action            string                 `json:"action,omitempty"`
	OpenNewTab        bool                   `json:"open_new_tab,omitempty"`
}

type ExternalNavigationData struct {
	URL  string `json:"url,omitempty"`
	URI  string `json:"uri,omitempty"`
	TYPE string `json:"type,omitempty"`
}

type Address struct {
	Line1    string `json:"address_line_1"`
	Line2    string `json:"address_line_2"`
	Line3    string `json:"address_line_3"`
	City     string `json:"city"`
	District string `json:"district"`
	State    string `json:"state"`
	PinCode  string `json:"pin_code"`
}

type WebLoginDrawer struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Forward Navigation  `json:"forward_action"`
}

type Color struct {
	Name    string `json:"name"`
	HexCode string `json:"hex_code"`
}
type InternalUserList struct {
	Name           string   `json:"name"`
	UserID         string   `json:"userId"`
	PhoneNumber    string   `json:"phone_number"`
	UserType       string   `json:"user_type"`
	EmailID        string   `json:"email_id"`
	Roles          []string `json:"role"`
	Center         string   `json:"center"`
	EmployeeID     string   `json:"employee_id"`
	ExternalUserID string   `json:"external_user_id"`
}

type InternalUserListResponse struct {
	Users        []InternalUserList `json:"users"`
	PageSize     int64              `json:"page_size"`
	CurrentPage  int64              `json:"current_page"`
	TotalPages   int64              `json:"total_pages"`
	TotalRecords int64              `json:"total_records"`
}

type InternalUserProfile struct {
	Name           string `json:"name"`
	EmployeeID     string `json:"employee_id"`
	UserID         string `json:"userId"`
	PhoneNumber    string `json:"phone_number"`
	UserType       string `json:"user_type"`
	WorkEmail      string `json:"work_email"`
	PersonalEmail  string `json:"personal_email"`
	ProfilePicURL  string `json:"profile_pic_url"`
	ActiveFrom     int64  `json:"active_from"`
	TenantID       string `json:"tenant_id"`
	ExternalUserID string `json:"external_user_id"`
}

type RoleInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
}

type PrivilegeWithAttributesInfo struct {
	Resource    string     `json:"resource"`
	AccessLevel string     `json:"access_level"`
	Attributes  Attributes `json:"attributes"`
}

type Attributes struct {
	Center   []string
	Campus   []string
	Stream   []string
	Course   []string
	Phase    []string
	Batch    []string
	Class    []string
	Subject  []string
	Topic    []string
	SubTopic []string
}

type RolePrivilegeWithAttributes struct {
	RoleInfo      RoleInfo                      `json:"role_info"`
	PrivilegeInfo []PrivilegeWithAttributesInfo `json:"privileges"`
}

type InternalUserAndRoleProfile struct {
	Name              string                        `json:"name"`
	EmployeeID        string                        `json:"employee_id"`
	UserID            string                        `json:"userId"`
	PhoneNumber       string                        `json:"phone_number"`
	UserType          string                        `json:"user_type"`
	WorkEmail         string                        `json:"work_email"`
	PersonalEmail     string                        `json:"personal_email"`
	ProfilePicURL     string                        `json:"profile_pic_url"`
	ActiveFrom        int64                         `json:"active_from"`
	TenantID          string                        `json:"tenant_id"`
	ExternalUserID    string                        `json:"external_user_id"`
	RolePrivilegeInfo []RolePrivilegeWithAttributes `json:"user_privileges"`
}

type DashboardCourseList struct {
	OrderID  string `json:"order_id"`
	CourseID string `json:"course_id"`
	PhaseID  string `json:"phase_id"`
	BatchID  string `json:"batch_id"`
	CenterID string `json:"center_id"`
	Class    string `json:"class"`
	Stream   string `json:"stream"`
	Course   string `json:"course"`
	Phase    string `json:"phase"`
	Batch    string `json:"batch"`
	Center   string `json:"center"`
}

type DashboardStudentList struct {
	Name            string                `json:"name"`
	UserID          string                `json:"user_id"`
	PhoneNumber     string                `json:"phone_number"`
	CoursesEnrolled []DashboardCourseList `json:"courses_enrolled"`
}
type DashboardStudentListingResponse struct {
	Students     []DashboardStudentList `json:"students"`
	PageSize     int64                  `json:"page_size"`
	CurrentPage  int64                  `json:"current_page"`
	TotalRecords int64                  `json:"total_records"`
}

type ScheduleStatusTextColor struct {
	Name    string `json:"name"`
	HexCode string `json:"hex_code"`
}

type DynamicString []DynamicText

type DynamicText struct {
	Text  string `json:"text"`
	Style Style  `json:"style"`
}

type Style struct {
	Name    string `json:"name,omitempty"`
	HexCode string `json:"hex_code,omitempty"`
}

type LottieInfo struct {
	Url       string `json:"url"`
	AltText   string `json:"alt_text"`
	Iteration string `json:"iteration,omitempty"`
}

type WidgetList struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type StringKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TextComponent struct {
	Text      string              `json:"text,omitempty"`
	Style     *TextComponentStyle `json:"style,omitempty"`
	ImageData *Image              `json:"image_data,omitempty"`
}

type TextComponentStyle struct {
	Font FontStyle `json:"font"`
}

type FontStyle struct {
	Color  Style  `json:"color"`
	Weight string `json:"weight,omitempty"`
}

type BackgroundColor struct {
	Name    string `json:"name"`
	HexCode string `json:"hex_code"`
}

type Ftue struct {
	Key string `json:"key"`
}

type MultiAction struct {
	Data           *MultiActionData `json:"data"`
	TrackingParams map[string]any   `json:"tracking_params,omitempty"`
	Type           string           `json:"type,omitempty"`
	Name           string           `json:"name,omitempty"`
}

type MultiActionData struct {
	TrackingParams map[string]any     `json:"tracking_params,omitempty"`
	Data           *MultiActionData   `json:"data,omitempty"`
	Name           string             `json:"name,omitempty"`
	Type           string             `json:"type,omitempty"`
	URI            string             `json:"uri,omitempty"`
	Query          map[string]any     `json:"query,omitempty"`
	Title          string             `json:"title,omitempty"`
	Label          string             `json:"label,omitempty"`
	State          string             `json:"state,omitempty"`
	ListActionData []*ListMultiAction `json:"list,omitempty"`
}

type ListMultiAction struct {
	Title      string      `json:"title"`
	Action     *RootAction `json:"action,omitempty"`
	IsSelected bool        `json:"is_selected"`
}

type RootAction struct {
	Data           RootActionData `json:"data"`
	Type           string         `json:"type"`
	TrackingParams map[string]any `json:"tracking_params,omitempty"`
}

type RootActionData struct {
	Type  string         `json:"type"`
	URI   string         `json:"uri,omitempty"`
	State string         `json:"state,omitempty"`
	Query map[string]any `json:"query,omitempty"`
	Title string         `json:"title,omitempty"`
	Label string         `json:"label,omitempty"`
}

type LabelData struct {
	ImageData Image           `json:"image_data"`
	BGColor   BackgroundColor `json:"background_color"`
	Title     string          `json:"title"`
}

type PaymentOptionCTA struct {
	Label      string              `json:"label"`
	Action     string              `json:"action"`
	StepID     string              `json:"stepId"`
	StepLabel  string              `json:"stepLabel"`
	RenderType string              `json:"renderType"`
	Steps      []PaymentOptionStep `json:"steps"`
}

type PaymentOptionStep struct {
	IconURL          string           `json:"iconUrl"`
	Description      string           `json:"description"`
	StepID           string           `json:"stepId"`
	CompletionAction CompletionAction `json:"completionAction"`
}

type CompletionAction struct {
	Type string             `json:"type"`
	Data CompleteActionData `json:"data"`
}

type CompleteActionData struct {
	URI   string                 `json:"uri"`
	Query map[string]interface{} `json:"query"`
}

type FtuxWidget struct {
	Name    string       `json:"name"`
	Payload *FtuxPayload `json:"payload"`
}

type FtuxPayload struct {
	Sticker    *Style      `json:"sticker,omitempty"`
	Title      string      `json:"title"`
	LabelName  string      `json:"label_name"`
	Subtitle   string      `json:"subtitle"`
	LottieData *LottieInfo `json:"lottie_data,omitempty"`
	PrimaryCta *CTA        `json:"primary_cta,omitempty"`
}

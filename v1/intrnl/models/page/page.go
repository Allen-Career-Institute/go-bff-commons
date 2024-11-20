package page

import (
	"encoding/json"
	calTypes "github.com/Allen-Career-Institute/common-protos/cal/v1/types"
	pbTypes "github.com/Allen-Career-Institute/common-protos/page_service/v1/types"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetPageRequest struct {
	PageURL     string            `json:"page_url"`
	UserContext map[string]string `json:"user_context"`
}

type Info struct {
	ID             string            `json:"id"`
	PageID         string            `json:"page_id"`
	PageMeta       Meta              `json:"page_meta"`
	Name           string            `json:"name"`
	Data           json.RawMessage   `json:"data"`
	LayoutParams   json.RawMessage   `json:"layout_params,omitempty"`
	TrackingParams json.RawMessage   `json:"tracking_params,omitempty"`
	SEOData        json.RawMessage   `json:"seo_data,omitempty"`
	OnloadActions  []*pbTypes.Action `json:"onload_actions,omitempty"`
	Actions        []*pbTypes.Action `json:"actions,omitempty"`
}

type Meta struct {
	PageType        string                   `json:"type"`
	URL             string                   `json:"url"`
	TabMeta         json.RawMessage          `json:"tab_meta,omitempty"`
	VisibilityRules *pbTypes.VisibilityRules `json:"-"`
	FloatingMeta    ArrangementMeta          `json:"floating_meta,omitempty"`
}

type ArrangementMeta struct {
	Stack       string `json:"stack,omitempty"`
	Arrangement string `json:"arrangement,omitempty"`
}

type CommonPageResponse struct {
	PageInfo    Info        `json:"page_info"`
	PageContent ContentData `json:"page_content"`
	TabData     []*TabData  `json:"tab_data,omitempty"`
}

type TabData struct {
	ID              uint32                `json:"id"`
	ConstTabID      string                `json:"const_tab_id"`
	TabID           string                `json:"tab_id"`
	Icon            string                `json:"icon"`
	SelectedIcon    string                `json:"selected_icon"`
	Name            string                `json:"name"`
	Selected        bool                  `json:"selected"`
	SelectedTabIcon *TabIcon              `json:"selected_tab_icon,omitempty"`
	DefaultTabIcon  *TabIcon              `json:"default_tab_icon,omitempty"`
	IntroTooltip    *pbTypes.IntroTooltip `json:"intro_tooltip,omitempty"`
	TabInfo         *TabInfo              `json:"tab_info"`
	TrackingParams  *structpb.Struct      `json:"tracking_params"`
	URL             string                `json:"url"`
}

type TabIcon struct {
	Type  string                       `json:"type"`
	Dark  *pbTypes.TabIcon_TabIconData `json:"dark"`
	Light *pbTypes.TabIcon_TabIconData `json:"light"`
}

type TabInfo struct {
	PageID    string              `json:"page_id"`
	PageType  string              `json:"page_type"`
	URL       string              `json:"url"`
	PageData  *CommonPageResponse `json:"page_data"`
	TabAction *ActionTab          `json:"tab_action"`
}

type ActionTab struct {
	Data *structpb.Struct `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Type string           `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

type ContentData struct {
	HeaderWidgets  []*WidgetData `json:"header_widgets"`
	Widgets        []*WidgetData `json:"widgets"`
	FooterWidgets  []*WidgetData `json:"footer_widgets"`
	OnloadWidgets  []*WidgetData `json:"-"`
	FloatingWidget []*WidgetData `json:"floating_widgets"`
}

// WidgetData TODO@Himanshu: TrackingParams && LayoutParams must be defined and should contain expected fields only
type WidgetData struct {
	ID             uint32           `json:"id"`
	ConstWidgetID  string           `json:"const_widget_id"`
	ViewType       string           `json:"type" validate:"omitempty"`
	TrackingParams *structpb.Struct `json:"tracking_params" validate:"omitempty"`
	LayoutParams   *structpb.Struct `json:"layout_params" validate:"omitempty"`
	Data           *structpb.Struct `json:"data" validate:"omitempty"`
	DataSource     string           `json:"-"`
	WidgetType     string           `json:"-"`
	IsProcessed    bool             `json:"-"`
}

type TabPageResponse struct {
}

type WidgetDataFromLmm struct {
	Content      interface{}                    `json:"content"`
	MaterialInfo *calTypes.LearningMaterialInfo `json:"material_info,omitempty"`
	NACInfo      *calTypes.NACInfo              `json:"nac_info,omitempty"`
}

func GetWidgetDataFromLmm(content interface{}, materialInfo *calTypes.LearningMaterialInfo, nacInfo *calTypes.NACInfo) *WidgetDataFromLmm {
	return &WidgetDataFromLmm{
		Content:      content,
		MaterialInfo: materialInfo,
		NACInfo:      nacInfo,
	}
}

package pageds

import (
	"fmt"
	pbTypes "github.com/Allen-Career-Institute/common-protos/page_service/v1/types"
	frameworkModels "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models/commons"
	internal "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl"
	pageResp "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/models/page"
	log "github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	constants "github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/structpb"
	"strconv"
)

type ResponseMapper struct {
	logger log.Logger
	helper *Helper
}

func NewResponseMapper(logger *log.Logger) *ResponseMapper {
	return &ResponseMapper{logger: *logger, helper: NewHelper()}
}

func (p *ResponseMapper) MapResponse(
	c echo.Context,
	pageInfo *pbTypes.PageInfo,
	resolvedWidgetsMap map[string]bool,
) (pageRes *pageResp.CommonPageResponse, dsl, wPosList []string, err error) {
	if pageInfo == nil || pageInfo.PageMeta == nil {
		return nil, nil, nil, fmt.Errorf("pageInfo data not valid : %s", pageInfo)
	}

	pageRes, dsl, wPosList, err = p.handleListPageResponse(c, pageInfo, resolvedWidgetsMap)
	if err != nil {
		return nil, nil, nil, err
	}

	return pageRes, dsl, wPosList, nil
}

func (p *ResponseMapper) handleListPageResponse(c echo.Context, pageInfo *pbTypes.PageInfo, resolvedWidgetsMap map[string]bool) (listPageResp *pageResp.CommonPageResponse, dsl, wl []string, err error) {
	listPageResp, err = p.mapToListPageResp(c, pageInfo)
	if err != nil {
		return nil, nil, nil, err
	}

	widgetIdxToDataMp := make(map[int]*pageResp.WidgetData)

	processWidgetList := func(widgets []*pbTypes.WidgetInfo, widgetPosType pbTypes.WidgetPosType) {
		for _, widget := range widgets {
			if widget.WidgetType == pbTypes.WidgetType_DYNAMIC {
				dsl = append(dsl, widget.WidgetData.DataSource)
				wl = append(wl, widgetPosType.String())
				idx := len(dsl) - 1
				populateWidgetDataMapHelper(widget, widgetIdxToDataMp, idx)
			} else {
				resolvedWidgetsMap[widget.ConstWidgetId] = true
			}
		}
	}

	processWidgetList(pageInfo.PageContentData.HeaderWidgets, pbTypes.WidgetPosType_HEADER)

	processWidgetList(pageInfo.PageContentData.Widgets, pbTypes.WidgetPosType_NORMAL)

	processWidgetList(pageInfo.PageContentData.FooterWidgets, pbTypes.WidgetPosType_FOOTER)

	processWidgetList(pageInfo.PageContentData.OnloadWidgets, pbTypes.WidgetPosType_ONLOAD)

	processWidgetList(pageInfo.PageContentData.FloatingWidgets, pbTypes.WidgetPosType_FLOATING)

	c.Set(constants.WidgetIndexToWidgetDataMap, widgetIdxToDataMp)
	return listPageResp, dsl, wl, nil
}

func populateWidgetDataMapHelper(w *pbTypes.WidgetInfo, widgetIdxToDataMp map[int]*pageResp.WidgetData, idx int) {
	widgetIdxToDataMp[idx] = &pageResp.WidgetData{
		ConstWidgetID: w.ConstWidgetId,
	}
	if w.WidgetData.DataSource == constants.ResolveLMMWidget || w.WidgetData.Type == constants.Polymorphic {
		widgetIdxToDataMp[idx].Data = w.WidgetData.Data
		widgetIdxToDataMp[idx].DataSource = w.WidgetData.DataSource
	}
}

func (p *ResponseMapper) HandleTabPageResponse(c echo.Context, pInfo *pbTypes.PageInfo) []*pageResp.TabData {
	var td []*pageResp.TabData
	td = []*pageResp.TabData{}
	tc := pInfo.TabData
	for i := range tc {
		tab := p.mapTabContentToTabData(c, tc[i])
		td = append(td, tab)
	}
	return td

}

func (p *ResponseMapper) mapToListPageResp(c echo.Context, pInfo *pbTypes.PageInfo) (*pageResp.CommonPageResponse, error) {
	pm := &pageResp.Meta{}
	var err error
	pm.URL = pInfo.PageMeta.Url
	pm.PageType = pInfo.PageMeta.PageType.String()
	pm.TabMeta, err = p.helper.ConvertStructToRawMessage(pInfo.PageMeta.TabMeta)
	if err != nil {
		return nil, err
	}
	pm.VisibilityRules = pInfo.PageMeta.VisibilityRules
	pm.FloatingMeta = p.mapArrangementMeta(pInfo.PageMeta.FloatingMeta)

	pi, err := p.mapPageMetaAndInfo(c, pInfo, pm)
	if err != nil {
		return nil, err
	}

	li := &pageResp.CommonPageResponse{
		PageInfo: *pi,
		PageContent: pageResp.ContentData{
			HeaderWidgets:  p.mapWidgetListForResp(pInfo.PageContentData.HeaderWidgets),
			Widgets:        p.mapWidgetListForResp(pInfo.PageContentData.Widgets),
			FooterWidgets:  p.mapWidgetListForResp(pInfo.PageContentData.FooterWidgets),
			OnloadWidgets:  p.mapWidgetListForResp(pInfo.PageContentData.OnloadWidgets),
			FloatingWidget: p.mapWidgetListForResp(pInfo.PageContentData.FloatingWidgets),
		},
	}
	return li, nil
}

func (p *ResponseMapper) mapArrangementMeta(meta *pbTypes.ArrangementMeta) pageResp.ArrangementMeta {
	res := pageResp.ArrangementMeta{}
	if meta == nil {
		return res
	}
	res.Stack = meta.Stack.String()
	res.Arrangement = meta.Arrangement.String()
	return res
}

func (p *ResponseMapper) mapPageMetaAndInfo(c echo.Context, pInfo *pbTypes.PageInfo, pm *pageResp.Meta) (*pageResp.Info, error) {
	if pInfo.Name == constants.CareCornerPageName {
		pInfo.PageMeta.TrackingParams = p.addTrackingParamsData(c)
	}

	pData, err := p.helper.ConvertStructToRawMessage(pInfo.PageMeta.Data)
	if err != nil {
		p.logger.WithContext(c).Errorf("error while handling pInfo.PageMeta.Data : %v", err)
		return nil, err
	}
	plp, err := p.helper.ConvertStructToRawMessage(pInfo.PageMeta.LayoutParams)
	if err != nil {
		p.logger.WithContext(c).Errorf("error while handling pInfo.PageMeta.LayoutParams : %v", err)
		return nil, err
	}
	ptp, err := p.helper.ConvertStructToRawMessage(pInfo.PageMeta.TrackingParams)
	if err != nil {
		p.logger.WithContext(c).Errorf("error while handling pInfo.PageMeta.TrackingParams : %v", err)
		return nil, err
	}
	pseo, err := p.helper.ConvertStructToRawMessage(pInfo.PageMeta.SeoData)
	if err != nil {
		p.logger.WithContext(c).Errorf("error while handling pInfo.PageMeta.SeoData : %v", err)
		return nil, err
	}

	pi := &pageResp.Info{
		ID:             strconv.Itoa(int(pInfo.Id)),
		PageID:         pInfo.PageId,
		PageMeta:       *pm,
		Name:           pInfo.Name,
		Data:           pData,
		LayoutParams:   plp,
		TrackingParams: ptp,
		SEOData:        pseo,
		OnloadActions:  pInfo.PageMeta.OnloadActions,
		Actions:        pInfo.PageMeta.Actions,
	}
	return pi, nil
}

func (p *ResponseMapper) mapWidgetListForResp(widgetInfos []*pbTypes.WidgetInfo) []*pageResp.WidgetData {
	wSlice := make([]*pageResp.WidgetData, 0)
	for _, w := range widgetInfos {
		wResp := p.mapWidgetDataForResponse(w)
		wSlice = append(wSlice, wResp)
	}
	return wSlice
}

// mapWidgetDataForResponse just adds 'widgetType' in widgetData
// as it is needed for mapping Dynamic Widgets during processing at bff
func (*ResponseMapper) mapWidgetDataForResponse(widget *pbTypes.WidgetInfo) *pageResp.WidgetData {
	hwCustom := &pageResp.WidgetData{
		ViewType:       widget.WidgetData.Type,
		TrackingParams: widget.WidgetData.TrackingParams,
		LayoutParams:   widget.WidgetData.LayoutParams,
		Data:           widget.WidgetData.Data,
		DataSource:     widget.WidgetData.DataSource,
		WidgetType:     widget.WidgetType.String(),
		ConstWidgetID:  widget.ConstWidgetId,
		ID:             widget.Id,
	}
	return hwCustom
}

func (p *ResponseMapper) MapDataSourceRespToLP(
	c echo.Context,
	ds, wt, widgetID string,
	lpr *pageResp.CommonPageResponse,
	wd *frameworkModels.DSResponse,
	resolvedWidgetsMap map[string]bool,
) (*pageResp.CommonPageResponse, error) {
	switch wt {
	case pbTypes.WidgetPosType_HEADER.String():
		hw := p.populateAndFilterDynamicWidgets(c, lpr.PageContent.HeaderWidgets, ds, widgetID, wd, resolvedWidgetsMap)
		lpr.PageContent.HeaderWidgets = hw
	case pbTypes.WidgetPosType_NORMAL.String():
		w := p.populateAndFilterDynamicWidgets(c, lpr.PageContent.Widgets, ds, widgetID, wd, resolvedWidgetsMap)
		lpr.PageContent.Widgets = w
	case pbTypes.WidgetPosType_FOOTER.String():
		fw := p.populateAndFilterDynamicWidgets(c, lpr.PageContent.FooterWidgets, ds, widgetID, wd, resolvedWidgetsMap)
		lpr.PageContent.FooterWidgets = fw
	case pbTypes.WidgetPosType_ONLOAD.String():
		ow := p.populateAndFilterDynamicWidgets(c, lpr.PageContent.OnloadWidgets, ds, widgetID, wd, resolvedWidgetsMap)
		lpr.PageContent.OnloadWidgets = ow
	case pbTypes.WidgetPosType_FLOATING.String():
		flw := p.populateAndFilterDynamicWidgets(c, lpr.PageContent.FloatingWidget, ds, widgetID, wd, resolvedWidgetsMap)
		lpr.PageContent.FloatingWidget = flw
	default:
		p.logger.WithContext(c).Errorf("error! unsupported widgetPosType received")
		return nil, fmt.Errorf("error! unsupported widgetPosType received")
	}
	return lpr, nil
}

func (p *ResponseMapper) populateAndFilterDynamicWidgets(
	c echo.Context,
	widgets []*pageResp.WidgetData,
	dsName, widgetID string,
	dataSourceResp *frameworkModels.DSResponse,
	resolvedWidgetsMap map[string]bool,
) []*pageResp.WidgetData {
	// slice to store valid widgets that will be returned to frontend to render
	var validWidgets []*pageResp.WidgetData
	validWidgets = []*pageResp.WidgetData{}

	for _, widget := range widgets {
		// populate or filter matching dynamic widget
		if p.isMatchingUnprocessedDynamicWidget(widget, widgetID) {
			// ignoring invalid widget
			if dataSourceResp == nil {
				p.logger.WithContext(c).Errorf("nil response received from ds : %s, disabling widget", dsName)

				resolvedWidgetsMap[widget.ConstWidgetID] = false

				continue
			}
			widgetRespData, err := p.helper.InterfaceToStructPb(dataSourceResp.Data)
			if err != nil {
				p.logger.WithContext(c).Errorf("error while converting incoming widget data from dsName to required format : %v, ds: %s", err, dsName)

				resolvedWidgetsMap[widget.ConstWidgetID] = false

				continue
			}

			// populating widget data in place
			widget.Data = widgetRespData
			if dataSourceResp.ViewType != constants.EmptyString {
				widget.ViewType = dataSourceResp.ViewType
			}
			// marking this dynamic widget as processed, so that even if we have multiple widgets with same name
			// we can identify unprocessed widget simply by checking this flag
			widget.IsProcessed = true
		}

		resolvedWidgetsMap[widget.ConstWidgetID] = true

		validWidgets = append(validWidgets, widget)
	}

	// Setting unsuccessful widgets to null
	for j := len(validWidgets); j < len(widgets); j++ {
		widgets[j] = nil
	}
	return validWidgets
}

func (*ResponseMapper) isMatchingUnprocessedDynamicWidget(widget *pageResp.WidgetData, widgetID string) bool {
	if widget.WidgetType == pbTypes.WidgetType_DYNAMIC.String() &&
		(widget.ConstWidgetID != "" && widget.ConstWidgetID == widgetID) && !widget.IsProcessed {
		return true
	}
	return false
}

func (p *ResponseMapper) mapTabContentToTabData(c echo.Context, tc *pbTypes.TabContent) *pageResp.TabData {
	if tc.Name == constants.CareCornerTabName {
		// calling helper function with context to update tracking params
		tc.TrackingParams = p.addTrackingParamsData(c)
	}

	td := &pageResp.TabData{}
	td.TabID = tc.TabId
	td.URL = tc.Url
	td.Name = tc.Name
	td.TrackingParams = tc.TrackingParams
	td.SelectedIcon = tc.SelectedIcon
	td.Selected = tc.Selected
	td.SelectedTabIcon = p.mapToTabIcon(tc.SelectedTabIcon)
	td.DefaultTabIcon = p.mapToTabIcon(tc.DefaultTabIcon)
	td.IntroTooltip = tc.IntroTooltip
	td.Icon = tc.Icon
	td.TabInfo = p.mapTabPageInfoToTabInfo(tc.TabInfo)
	td.ID = tc.Id
	td.ConstTabID = tc.ConstTabId
	return td
}

// addTrackingParamsData is a helper function that creates a structured protobuf Struct
// with a "current" field containing an "enrolled" field set to the provided isEnrolled string value.
func (*ResponseMapper) addTrackingParamsData(c echo.Context) *structpb.Struct {

	// Fetching userEnrolled status from context
	isEnrolled := internal.GetUserEnrolledStatus(c)

	var isEnrolledFieldValue = &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: isEnrolled,
		},
	}

	var trackingParams = &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"current": {
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"enrolled": isEnrolledFieldValue,
						},
					},
				},
			},
			"global": {
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{},
					},
				},
			},
		},
	}

	return trackingParams

}

// TODO :Handle data inside TAB as List page resp
func (*ResponseMapper) mapTabPageInfoToTabInfo(tpi *pbTypes.TabPageInfo) *pageResp.TabInfo {
	ti := &pageResp.TabInfo{}
	ti.PageID = tpi.PageId
	ti.URL = tpi.Url
	ti.PageType = tpi.PageType.String()

	if tpi.TabAction != nil {
		if tpi.TabAction.Type != pbTypes.ActionType_NO_ACTION_TYPE {
			ti.TabAction = &pageResp.ActionTab{
				Data: tpi.TabAction.Data,
				Type: pbTypes.ActionType_name[int32(tpi.TabAction.Type)],
			}
		}
	}

	return ti
}

func (*ResponseMapper) mapToTabIcon(tpi *pbTypes.TabIcon) *pageResp.TabIcon {
	if tpi == nil {
		return nil
	}
	ti := &pageResp.TabIcon{}
	ti.Type = tpi.Type.String()
	ti.Light = tpi.Light
	ti.Dark = tpi.Dark
	return ti
}

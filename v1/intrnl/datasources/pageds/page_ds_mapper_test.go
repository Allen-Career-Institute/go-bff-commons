package pageds

import (
	"errors"
	pbTypes "github.com/Allen-Career-Institute/common-protos/page_service/v1/types"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	commonModels "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models/commons"
	pageResp "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/models/page"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	constants "github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//nolint:funlen,gocyclo,gocritic
func getTestingParams(t *testing.T) (*gomock.Controller, *config.Config, *echo.Echo, logger.Logger) {
	ctrl := gomock.NewController(t)

	e := echo.New()

	c := &config.Config{Logger: config.Logger{Level: "info"}}
	log := logger.NewAPILogger(c)
	log.InitLogger()

	return ctrl, c, e, log
}

func TestHandleListPageResponse(t *testing.T) {
	_, _, e, log := getTestingParams(t)
	// Define a minimal valid PageInfo with dynamic widgets
	dynamicPageInfo := &pbTypes.PageInfo{
		PageMeta: &pbTypes.PageMeta{
			Url:             "http://test.com",
			PageType:        pbTypes.PageMeta_LIST,
			VisibilityRules: nil, // Fill in according to your logic
			FloatingMeta:    &pbTypes.ArrangementMeta{},
			TabMeta:         nil, // Handle nil values correctly
			OnloadActions:   nil,
		},
		PageContentData: &pbTypes.PageContentData{
			HeaderWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         1,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: constants.ResolveLMMWidget},
				},
			},
			Widgets: []*pbTypes.WidgetInfo{
				{
					Id:         2,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: constants.ResolveLMMWidget},
				},
			},
			FooterWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         3,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: constants.ResolveLMMWidget},
				},
			},
			OnloadWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         4,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: "dynamic-onload"},
				},
			},
			FloatingWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         5,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: "dynamic-floating"},
				},
			},
		},
	}
	staticPageInfo := &pbTypes.PageInfo{
		PageMeta: &pbTypes.PageMeta{
			Url:             "http://test.com",
			PageType:        pbTypes.PageMeta_LIST,
			VisibilityRules: nil, // Fill in according to your logic
			FloatingMeta:    &pbTypes.ArrangementMeta{},
			TabMeta:         nil, // Handle nil values correctly
			OnloadActions:   nil,
		},
		PageContentData: &pbTypes.PageContentData{
			HeaderWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         1,
					WidgetType: pbTypes.WidgetType_STATIC,
					WidgetData: &pbTypes.WidgetData{},
				},
			},
			Widgets: []*pbTypes.WidgetInfo{
				{
					Id:         2,
					WidgetType: pbTypes.WidgetType_STATIC,
					WidgetData: &pbTypes.WidgetData{},
				},
			},
			FooterWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         3,
					WidgetType: pbTypes.WidgetType_STATIC,
					WidgetData: &pbTypes.WidgetData{},
				},
			},
			OnloadWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         4,
					WidgetType: pbTypes.WidgetType_STATIC,
					WidgetData: &pbTypes.WidgetData{},
				},
			},
			FloatingWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         5,
					WidgetType: pbTypes.WidgetType_STATIC,
					WidgetData: &pbTypes.WidgetData{},
				},
			},
		},
	}
	pageInfo := &pbTypes.PageInfo{
		PageMeta: &pbTypes.PageMeta{
			Url:             "http://test.com",
			PageType:        pbTypes.PageMeta_LIST,
			VisibilityRules: nil, // Fill in according to your logic
			FloatingMeta:    &pbTypes.ArrangementMeta{},
			TabMeta:         nil, // Handle nil values correctly
			OnloadActions:   nil,
		},
		PageContentData: &pbTypes.PageContentData{
			HeaderWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         1,
					WidgetType: pbTypes.WidgetType_STATIC,
					WidgetData: &pbTypes.WidgetData{},
				},
			},
			Widgets: []*pbTypes.WidgetInfo{
				{
					Id:         2,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: constants.ResolveLMMWidget},
				},
			},
			FooterWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         3,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: constants.ResolveLMMWidget},
				},
			},
			OnloadWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         4,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: "dynamic-onload"},
				},
			},
			FloatingWidgets: []*pbTypes.WidgetInfo{
				{
					Id:         5,
					WidgetType: pbTypes.WidgetType_DYNAMIC,
					WidgetData: &pbTypes.WidgetData{DataSource: "dynamic-floating"},
				},
			},
		},
	}
	tests := []struct {
		name               string
		pageInfo           *pbTypes.PageInfo
		resolvedWidgetsMap map[string]bool
		expectedDSL        []string
		expectedWL         []string
		expectError        bool
	}{
		{
			name:               "Empty PageContentData",
			pageInfo:           staticPageInfo,
			resolvedWidgetsMap: make(map[string]bool),
			expectedDSL:        []string(nil),
			expectedWL:         []string(nil),
			expectError:        false,
		},
		{
			name:               "Only Static Widgets in Header",
			pageInfo:           staticPageInfo,
			resolvedWidgetsMap: map[string]bool{},
			expectedDSL:        []string(nil),
			expectedWL:         []string(nil),
			expectError:        false,
		},
		{
			name:               "Dynamic Header Widget",
			pageInfo:           dynamicPageInfo,
			resolvedWidgetsMap: map[string]bool{},
			expectedDSL:        []string{"ResolveLMMWigdets", "ResolveLMMWigdets", "ResolveLMMWigdets", "dynamic-onload", "dynamic-floating"},
			expectedWL:         []string{"HEADER", "NORMAL", "FOOTER", "ONLOAD", "FLOATING"},
			expectError:        false,
		},
		{
			name:               "Mixed Static and Dynamic Widgets in Body",
			pageInfo:           pageInfo,
			resolvedWidgetsMap: map[string]bool{},
			expectedDSL:        []string{"ResolveLMMWigdets", "ResolveLMMWigdets", "dynamic-onload", "dynamic-floating"},
			expectedWL:         []string{"NORMAL", "FOOTER", "ONLOAD", "FLOATING"},
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an Echo context for the test
			req := e.NewContext(nil, nil)
			// Create a new ResponseMapper
			rm := NewResponseMapper(&log)

			// Run the function under test
			pageResp, dsl, wl, err := rm.handleListPageResponse(req, tt.pageInfo, tt.resolvedWidgetsMap)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pageResp)
				assert.Equal(t, tt.expectedDSL, dsl)
				assert.Equal(t, tt.expectedWL, wl)
			}
		})
	}
}

func TestMapDataSourceRespToLP(t *testing.T) {
	_, _, e, log := getTestingParams(t)
	tests := []struct {
		name            string
		wt              string
		initialWidgets  []*pageResp.WidgetData
		expectedWidgets []*pageResp.WidgetData
		widgetPosType   string
		expectedError   error
	}{
		{
			name: "Header Widgets",
			wt:   pbTypes.WidgetPosType_HEADER.String(),
			initialWidgets: []*pageResp.WidgetData{
				{ID: 1, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			expectedWidgets: []*pageResp.WidgetData{
				{ID: 1, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			widgetPosType: pbTypes.WidgetPosType_HEADER.String(),
			expectedError: nil,
		},
		{
			name: "ERROR: Header Dynamic Widgets",
			wt:   pbTypes.WidgetPosType_HEADER.String(),
			initialWidgets: []*pageResp.WidgetData{
				{ID: 1, WidgetType: pbTypes.WidgetType_DYNAMIC.String(), ConstWidgetID: "1"},
			},
			expectedWidgets: []*pageResp.WidgetData{},
			widgetPosType:   pbTypes.WidgetPosType_HEADER.String(),
			expectedError:   nil,
		},
		{
			name: "Normal Widgets",
			wt:   pbTypes.WidgetPosType_NORMAL.String(),
			initialWidgets: []*pageResp.WidgetData{
				{ID: 2, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			expectedWidgets: []*pageResp.WidgetData{
				{ID: 2, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			widgetPosType: pbTypes.WidgetPosType_NORMAL.String(),
			expectedError: nil,
		},
		{
			name: "Footer Widgets",
			wt:   pbTypes.WidgetPosType_FOOTER.String(),
			initialWidgets: []*pageResp.WidgetData{
				{ID: 2, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			expectedWidgets: []*pageResp.WidgetData{
				{ID: 2, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			widgetPosType: pbTypes.WidgetPosType_FOOTER.String(),
			expectedError: nil,
		},
		{
			name: "Onload Widgets",
			wt:   pbTypes.WidgetPosType_ONLOAD.String(),
			initialWidgets: []*pageResp.WidgetData{
				{ID: 2, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			expectedWidgets: []*pageResp.WidgetData{
				{ID: 2, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			widgetPosType: pbTypes.WidgetPosType_ONLOAD.String(),
			expectedError: nil,
		},
		{
			name: "Floating Widgets",
			wt:   pbTypes.WidgetPosType_FLOATING.String(),
			initialWidgets: []*pageResp.WidgetData{
				{ID: 2, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			expectedWidgets: []*pageResp.WidgetData{
				{ID: 2, WidgetType: pbTypes.WidgetType_STATIC.String()},
			},
			widgetPosType: pbTypes.WidgetPosType_FLOATING.String(),
			expectedError: nil,
		},
		{
			name: "Unsupported widget position",
			wt:   "unsupported",
			initialWidgets: []*pageResp.WidgetData{
				{ID: 3, WidgetType: pbTypes.WidgetType_DYNAMIC.String()},
			},
			expectedWidgets: []*pageResp.WidgetData{},
			widgetPosType:   "unsupported",
			expectedError:   errors.New("error! unsupported widgetPosType received"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock data for lpr
			lpr := &pageResp.CommonPageResponse{
				PageContent: pageResp.ContentData{
					HeaderWidgets:  make([]*pageResp.WidgetData, 0),
					Widgets:        make([]*pageResp.WidgetData, 0),
					FooterWidgets:  make([]*pageResp.WidgetData, 0),
					OnloadWidgets:  make([]*pageResp.WidgetData, 0),
					FloatingWidget: make([]*pageResp.WidgetData, 0),
				},
			}

			switch tt.widgetPosType {
			case pbTypes.WidgetPosType_HEADER.String():
				lpr.PageContent.HeaderWidgets = tt.initialWidgets
			case pbTypes.WidgetPosType_NORMAL.String():
				lpr.PageContent.Widgets = tt.initialWidgets
			case pbTypes.WidgetPosType_FOOTER.String():
				lpr.PageContent.FooterWidgets = tt.initialWidgets
			case pbTypes.WidgetPosType_FLOATING.String():
				lpr.PageContent.FloatingWidget = tt.initialWidgets
			case pbTypes.WidgetPosType_ONLOAD.String():
				lpr.PageContent.OnloadWidgets = tt.initialWidgets
			}

			// Create ResponseMapper instance
			rm := &ResponseMapper{
				logger: log,
			}

			// Prepare other input parameters
			ds := "data-source"
			widgetId := "1"
			wd := &commonModels.DSResponse{}
			resolvedWidgetsMap := make(map[string]bool)

			req := httptest.NewRequest(http.MethodPatch, "/", nil)
			rec := httptest.NewRecorder()

			ctx := e.NewContext(req, rec)
			// Execute the function
			result, err := rm.MapDataSourceRespToLP(ctx, ds, tt.wt, widgetId, lpr, wd, resolvedWidgetsMap)

			// Check for expected error
			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedError != nil {
				// Check expected results
				switch tt.widgetPosType {
				case pbTypes.WidgetPosType_HEADER.String():
					assert.Equal(t, tt.expectedWidgets, result.PageContent.HeaderWidgets)
				case pbTypes.WidgetPosType_NORMAL.String():
					assert.Equal(t, tt.expectedWidgets, result.PageContent.Widgets)
				case pbTypes.WidgetPosType_FOOTER.String():
					assert.Equal(t, tt.expectedWidgets, result.PageContent.FooterWidgets)
				case pbTypes.WidgetPosType_FLOATING.String():
					assert.Equal(t, tt.expectedWidgets, result.PageContent.FloatingWidget)
				case pbTypes.WidgetPosType_ONLOAD.String():
					assert.Equal(t, tt.expectedWidgets, result.PageContent.OnloadWidgets)
				}
			}

		})
	}
}
func TestMapTabPageInfoToTabInfo(t *testing.T) {
	tests := []struct {
		name           string
		input          *pbTypes.TabPageInfo
		expectedOutput *pageResp.TabInfo
	}{
		{
			name: "Valid Input with SHOW_TOAST Action",
			input: &pbTypes.TabPageInfo{
				PageId:   "page123",
				Url:      "https://example.com",
				PageType: 2,
				TabAction: &pbTypes.ActionTab{
					Data: nil,
					Type: 1,
				},
			},
			expectedOutput: &pageResp.TabInfo{
				PageID:    "page123",
				URL:       "https://example.com",
				PageType:  "2",
				TabAction: &pageResp.ActionTab{Data: nil, Type: "SHOW_TOAST"},
			},
		},
		{
			name: "Nil TabAction",
			input: &pbTypes.TabPageInfo{
				PageId:    "page123",
				Url:       "https://example.com",
				PageType:  2,
				TabAction: nil,
			},
			expectedOutput: &pageResp.TabInfo{
				PageID:    "page123",
				URL:       "https://example.com",
				PageType:  "2",
				TabAction: nil,
			},
		},
		{
			name: "Nil TabAction Data",
			input: &pbTypes.TabPageInfo{
				PageId:   "page123",
				Url:      "https://example.com",
				PageType: 2,
				TabAction: &pbTypes.ActionTab{
					Data: nil,
					Type: 0,
				},
			},
			expectedOutput: &pageResp.TabInfo{
				PageID:    "page123",
				URL:       "https://example.com",
				PageType:  "2",
				TabAction: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := &ResponseMapper{}
			result := rm.mapTabPageInfoToTabInfo(tt.input)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

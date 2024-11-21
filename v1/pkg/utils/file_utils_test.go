// nolint: gocritic,
package utils

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestGetFileFromRequest(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func(echo.Context)
		wantErr      bool
	}{
		{
			name: "no file",
			setupContext: func(c echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				c.SetRequest(req)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			tt.setupContext(c)

			cfg := &config.Config{Logger: config.Logger{Level: "info"}}
			logger := logger.NewAPILogger(cfg)
			logger.InitLogger()

			_, err := GetFileFromRequest(c, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileFromRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetParsedExcelSheet(t *testing.T) {
	tests := []struct {
		name      string
		fileData  []byte
		sheetName string
		wantErr   bool
	}{
		{
			name:      "invalid excel file",
			fileData:  []byte("<invalid excel file content>"),
			sheetName: "validSheet",
			wantErr:   true,
		},
		{
			name:      "valid excel file but invalid sheet",
			fileData:  []byte("<valid excel file content>"),
			sheetName: "invalidSheet",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			cfg := &config.Config{Logger: config.Logger{Level: "info"}}
			log := logger.NewAPILogger(cfg)
			log.InitLogger()

			_, err := GetParsedExcelSheet(c, tt.fileData, log, tt.sheetName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParsedExcelSheet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertExcelSheetToBytes(t *testing.T) {
	tests := []struct {
		name      string
		sheet     [][]string
		sheetName string
		wantErr   bool
	}{
		{
			name:      "valid sheet",
			sheet:     [][]string{{"header1", "header2"}, {"data1", "data2"}},
			sheetName: "validSheet",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("POST", "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			cfg := &config.Config{Logger: config.Logger{Level: "info"}}
			logger := logger.NewAPILogger(cfg)
			logger.InitLogger()

			_, err := ConvertExcelSheetToBytes(c, logger, tt.sheet, tt.sheetName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertExcelSheetToBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

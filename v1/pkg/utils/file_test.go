// nolint : funlen,
package utils

import (
	"errors"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPopulateErrorData(t *testing.T) {
	tests := []struct {
		name       string
		parsedData [][]string
		errors     []string
		expected   [][]string
	}{
		{
			name:       "No errors",
			parsedData: [][]string{{"data1"}, {"data2"}},
			errors:     []string{},
			expected:   [][]string{{"data1"}, {"data2"}},
		},
		{
			name:       "Single error",
			parsedData: [][]string{{"data1"}, {"data2"}},
			errors:     []string{"error1"},
			expected:   [][]string{{"data1"}, {"data2"}, {"error1"}},
		},
		{
			name:       "Multiple errors",
			parsedData: [][]string{{"data1"}, {"data2"}},
			errors:     []string{"error1", "error2"},
			expected:   [][]string{{"data1"}, {"data2"}, {"error1"}, {"error2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PopulateErrorData(tt.parsedData, tt.errors)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PopulateErrorData() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPopulateErrorDataRowWise(t *testing.T) {
	tests := []struct {
		name       string
		parsedData [][]string
		errors     map[int]string
		expected   [][]string
	}{
		{
			name:       "No errors",
			parsedData: [][]string{{"data1"}, {"data2"}},
			errors:     map[int]string{},
			expected:   [][]string{{"data1"}, {"data2"}},
		},
		{
			name:       "Single error",
			parsedData: [][]string{{"data1"}, {"data2"}},
			errors:     map[int]string{1: "error1"},
			expected:   [][]string{{"data1"}, {"data2", "error1"}},
		},
		{
			name:       "Multiple errors",
			parsedData: [][]string{{"data1"}, {"data2"}},
			errors:     map[int]string{0: "error1", 1: "error2"},
			expected:   [][]string{{"data1", "error1"}, {"data2", "error2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PopulateErrorDataRowWise(tt.parsedData, tt.errors)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PopulateErrorDataRowWise() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetParsedData(t *testing.T) {
	tests := []struct {
		name          string
		fileData      []byte
		sheetName     string
		expectedError error
	}{
		{
			name:          "Error case: invalid sheet name",
			fileData:      []byte("valid excel data"),
			sheetName:     "invalid sheet name",
			expectedError: errors.New("zip: not a valid zip file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			log := logger.NewAPILogger(&config.Config{Logger: config.Logger{Level: "info"}})
			log.InitLogger()

			_, err := GetParsedData(c, tt.fileData, tt.sheetName, log)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestAddErrorDataToExcel(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		data     []byte
		err      error
	}{
		{
			name:     "No error",
			fileName: "test.xlsx",
			data:     []byte("test data"),
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			cfg := &config.Config{Logger: config.Logger{Level: "info"}}
			logger := logger.NewAPILogger(cfg)
			logger.InitLogger()

			_, err := AddErrorDataToExcel(c, tt.fileName, tt.data, logger, tt.err)

			assert.Equal(t, "attachment; filename="+tt.fileName, c.Response().Header().Get(echo.HeaderContentDisposition))
			assert.Equal(t, echo.MIMEOctetStream, c.Response().Header().Get(echo.HeaderContentType))

			assert.Equal(t, http.StatusBadRequest, c.Response().Status)

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetDataFromExcel(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		fileOpenErr bool
		fileReadErr bool
		wantErr     bool
	}{
		{
			name:        "Error opening file",
			fileContent: "test data",
			fileOpenErr: true,
			fileReadErr: false,
			wantErr:     true,
		},
		{
			name:        "Error reading file",
			fileContent: "test data",
			fileOpenErr: false,
			fileReadErr: true,
			wantErr:     true,
		},
		{
			name:        "Non-empty file content, error opening file",
			fileContent: "test data",
			fileOpenErr: true,
			fileReadErr: false,
			wantErr:     true,
		},
		{
			name:        "Non-empty file content, error reading file",
			fileContent: "test data",
			fileOpenErr: false,
			fileReadErr: true,
			wantErr:     true,
		},
		{
			name:        "Empty file content, error opening file",
			fileContent: "",
			fileOpenErr: true,
			fileReadErr: false,
			wantErr:     true,
		},
		{
			name:        "Empty file content, error reading file",
			fileContent: "",
			fileOpenErr: false,
			fileReadErr: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.fileContent))
			req.Header.Set(echo.HeaderContentType, echo.MIMEOctetStream)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("formFile", &multipart.FileHeader{})

			if tt.fileOpenErr {
				c.Set("fileOpen", func() (multipart.File, error) {
					return nil, errors.New("mocked file open error")
				})
			} else {
				c.Set("fileOpen", func() (multipart.File, error) {
					return nil, nil
				})
			}

			if tt.fileReadErr {
				c.Set("fileRead", func(r io.Reader) ([]byte, error) {
					return nil, errors.New("mocked file read error")
				})
			} else {
				c.Set("fileRead", io.ReadAll)
			}

			log := logger.NewAPILogger(&config.Config{Logger: config.Logger{Level: "info"}})
			log.InitLogger()

			_, err := GetDataFromExcel(c, log)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// nolint: gocritic,bodyclose
package utils

import (
	"bff-service/config"
	"bff-service/pkg/logger"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestReadRequest(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		wantErr     bool
	}{
		{
			name:        "valid request",
			requestBody: `{"field": "value"}`,
			wantErr:     false,
		},
		{
			name:        "invalid request",
			requestBody: `{"field": ""}`,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var body TestRequest
			err := ReadRequest(c, &body)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type TestRequest struct {
	Field string `json:"field" validate:"required"`
}

func TestGetRequestID(t *testing.T) {
	tests := []struct {
		name      string
		requestID string
	}{
		{
			name:      "Test GetRequestID with valid request ID",
			requestID: "12345",
		},
		{
			name:      "Test GetRequestID with empty request ID",
			requestID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Response().Header().Set(echo.HeaderXRequestID, tt.requestID)

			result := GetRequestID(c)
			assert.Equal(t, tt.requestID, result)
		})
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "Test HandleError with intrnl server error",
			err:             errors.New("intrnl server error"),
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "intrnl server error",
		},
		{
			name:            "Test HandleError with not found error",
			err:             status.Error(codes.NotFound, "not found error"),
			expectedCode:    http.StatusNotFound,
			expectedMessage: "not found error",
		},
		{
			name:            "Test HandleError with unauthorized error",
			err:             status.Error(codes.Unauthenticated, "unauthorized error"),
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "unauthorized error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			c := &config.Config{Logger: config.Logger{Level: "info"}}
			log := logger.NewAPILogger(c)
			log.InitLogger()
			code, message := HandleError(ctx, tt.err, log)

			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedMessage, message)
		})
	}
}

func TestGetRequestCtx(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			name: "Test with request ID",
			req: func() *http.Request {
				req := httptest.NewRequest(echo.GET, "/", nil)
				req.Header.Set(echo.HeaderXRequestID, "test-request-id")
				return req
			}(),
		},
		{
			name: "Test without request ID",
			req:  httptest.NewRequest(echo.GET, "/", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			c := e.NewContext(tt.req, httptest.NewRecorder())

			ctx := GetRequestCtx(c)

			reqID, ok := ctx.Value(ReqIDCtxKey{}).(string)
			assert.True(t, ok)
			assert.Equal(t, GetRequestID(c), reqID)
		})
	}
}

func TestGetRequestCtxWithTimeout(t *testing.T) {
	tests := []struct {
		name      string
		requestID string
		timeout   time.Duration
	}{
		{
			name:      "Test with valid request ID and timeout",
			requestID: "12345",
			timeout:   1 * time.Millisecond,
		},
		{
			name:      "Test with empty request ID and timeout",
			requestID: "",
			timeout:   1 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Response().Header().Set(echo.HeaderXRequestID, tt.requestID)

			ctx, cancel := GetRequestCtxWithTimeout(c, tt.timeout)
			defer cancel()

			reqID, ok := ctx.Value(ReqIDCtxKey{}).(string)
			assert.True(t, ok)
			assert.Equal(t, tt.requestID, reqID)
		})
	}
}

func TestAddAuthHeaderAsMetadata(t *testing.T) {
	tests := []struct {
		name              string
		auth              string
		platform          string
		deviceID          string
		visitorID         string
		requestID         string
		expectedAuth      string
		expectedDevice    string
		expectedReqID     string
		expectedDeviceID  string
		expectedVisitorID string
	}{
		{
			name:           "Test with valid auth, platform and requestID",
			auth:           "Bearer token",
			platform:       "mobile",
			deviceID:       "random-device-id",
			visitorID:      "vis-12345",
			requestID:      "12345",
			expectedAuth:   "Bearer token",
			expectedDevice: "mobile",
			expectedReqID:  "12345",
		},
		{
			name:              "Test with empty auth",
			auth:              "",
			platform:          "web",
			deviceID:          "random-device-id",
			visitorID:         "vis-12345",
			requestID:         "67890",
			expectedAuth:      "",
			expectedDevice:    "web",
			expectedReqID:     "",
			expectedDeviceID:  "random-device-id",
			expectedVisitorID: "vis-12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.Header.Set(echo.HeaderAuthorization, tt.auth)
			req.Header.Set(DeviceType, tt.platform)
			req.Header.Set(DeviceID, tt.deviceID)
			req.Header.Set(VisitorID, tt.visitorID)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Response().Header().Set(echo.HeaderXRequestID, tt.requestID)

			ctx := AddAuthHeaderAsMetadata(context.Background(), c)

			md, _ := metadata.FromOutgoingContext(ctx)
			auth := md.Get(echo.HeaderAuthorization)
			device := md.Get(DeviceType)
			reqID := md.Get(echo.HeaderXRequestID)
			deviceID := md.Get(DeviceID)
			visitorID := md.Get(VisitorID)

			if tt.auth != "" {
				assert.Equal(t, tt.expectedAuth, auth[0])
				assert.Equal(t, tt.expectedDevice, device[0])
				assert.Equal(t, tt.expectedReqID, reqID[0])
			} else {
				assert.Empty(t, auth)
				assert.Empty(t, reqID)
				assert.Equal(t, tt.expectedDevice, device[0])
				assert.Equal(t, tt.expectedDeviceID, deviceID[0])
				assert.Equal(t, tt.expectedVisitorID, visitorID[0])
			}
		})
	}
}

func TestGetIPAddress(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "Test GetIPAddress with valid IP",
			ip:       "192.168.1.1:1234",
			expected: "192.168.1.1:1234",
		},
		{
			name:     "Test GetIPAddress with empty IP",
			ip:       "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.RemoteAddr = tt.ip
			c := e.NewContext(req, httptest.NewRecorder())

			result := GetIPAddress(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAllowedImagesContentTypes(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "Test GetAllowedImagesContentTypes",
			want: map[string]string{
				"image/bmp":                "bmp",
				"image/gif":                "gif",
				"image/png":                "png",
				"image/jpeg":               "jpeg",
				"image/jpg":                "jpg",
				"image/svg+xml":            "svg",
				"image/webp":               "webp",
				"image/tiff":               "tiff",
				"image/vnd.microsoft.icon": "ico",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAllowedImagesContentTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAllowedImagesContentTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHTTPStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		grpcStatusCode codes.Code
		expected       int
	}{
		{
			name:           "Test OK",
			grpcStatusCode: codes.OK,
			expected:       http.StatusOK,
		},
		{
			name:           "Test Canceled",
			grpcStatusCode: codes.Canceled,
			expected:       449,
		},
		{
			name:           "Test Unknown",
			grpcStatusCode: codes.Unknown,
			expected:       http.StatusInternalServerError,
		},
		// Add more test cases for other gRPC status codes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetHTTPStatusCode(tt.grpcStatusCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetErrorCodeAndMessage(t *testing.T) {
	tests := []struct {
		name           string
		grpcStatusCode codes.Code
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "Test OK",
			grpcStatusCode: codes.OK,
			expectedCode:   200,
			expectedMsg:    "",
		},
		{
			name:           "Test Canceled",
			grpcStatusCode: codes.Canceled,
			expectedCode:   449,
			expectedMsg:    "context canceled",
		},
		{
			name:           "Test Unknown",
			grpcStatusCode: codes.Unknown,
			expectedCode:   500,
			expectedMsg:    "unknown error",
		},
		// Add more test cases for other gRPC status codes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := status.Error(tt.grpcStatusCode, tt.expectedMsg)
			resultCode, resultMsg := GetErrorCodeAndMessage(err)
			assert.Equal(t, tt.expectedCode, resultCode)
			assert.Equal(t, tt.expectedMsg, resultMsg)
		})
	}
}

func TestGetHeader(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		value    string
		expected string
	}{
		{
			name:     "Test GetHeader with existing header",
			header:   "Test-Header",
			value:    "Test Value",
			expected: "Test Value",
		},
		{
			name:     "Test GetHeader with non-existing header",
			header:   "Non-Existing-Header",
			value:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.Header.Set(tt.header, tt.value)
			c := e.NewContext(req, httptest.NewRecorder())

			result := GetHeader(c, tt.header)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheckImageFileContentType(t *testing.T) {
	tests := []struct {
		name             string
		fileContent      []byte
		expectedResponse string
		expectedError    error
	}{
		{
			name:             "Valid JPEG",
			fileContent:      []byte("BM"),
			expectedResponse: "bmp",
		},
		{
			name:             "Valid PNG",
			fileContent:      []byte("PNG"),
			expectedResponse: "png",
		},
		{
			name:          "Invalid content type",
			fileContent:   []byte("cdf"),
			expectedError: errors.New("this content type is not allowed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentType, err := CheckImageFileContentType(tt.fileContent)

			if err != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedResponse, contentType)
			}
		})
	}
}

func TestDeleteSessionCookie(t *testing.T) {
	tests := []struct {
		name          string
		sessionName   string
		expectedValue string
	}{
		{
			name:          "Delete session cookie",
			sessionName:   "session",
			expectedValue: "some_value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctx.SetCookie(&http.Cookie{
				Name:  tt.sessionName,
				Value: "some_value",
				Path:  "/",
			})

			DeleteSessionCookie(ctx, tt.sessionName)

			cookies := rec.Result().Cookies()

			var deletedCookie *http.Cookie

			for _, cookie := range cookies {
				if cookie.Name == tt.sessionName {
					deletedCookie = cookie
					break
				}
			}

			assert.NotNil(t, deletedCookie, "Expected cookie not found")
			assert.Equal(t, tt.expectedValue, deletedCookie.Value, "Cookie value after deletion doesn't match expected")
		})
	}
}

func TestGetRequestCtxWithClientIDAndTimeout(t *testing.T) {
	tests := []struct {
		name            string
		timeout         time.Duration
		clientRequestID string
	}{
		{
			name:            "Test with 1 second timeout",
			timeout:         1 * time.Second,
			clientRequestID: "testID1",
		},
		{
			name:            "Test with 5 seconds timeout",
			timeout:         5 * time.Second,
			clientRequestID: "testID2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			c := e.NewContext(req, httptest.NewRecorder())

			ctx, cancel := GetRequestCtxWithClientIDAndTimeout(c, tt.timeout, tt.clientRequestID)
			defer cancel()

			md, _ := metadata.FromOutgoingContext(ctx)
			clientRequestIDs := md.Get(ClientRequestID)

			if len(clientRequestIDs) != 1 || clientRequestIDs[0] != tt.clientRequestID {
				t.Errorf("Expected ClientRequestID to be %s, got %v", tt.clientRequestID, clientRequestIDs)
			}

			deadline, ok := ctx.Deadline()
			if !ok || time.Until(deadline) > tt.timeout {
				t.Errorf("Expected context timeout to be less than or equal to %v, got %v", tt.timeout, time.Until(deadline))
			}
		})
	}
}

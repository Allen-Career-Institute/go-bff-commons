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

func TestNewEchoUtil(t *testing.T) {
	// Define test cases
	d := &config.CommonConfig{Logger: config.Logger{Level: "info"}}
	log := logger.NewAPILogger(d)
	log.InitLogger()
	testCases := []struct {
		name string
		log  logger.Logger
	}{
		{
			name: "Test with nil logger",
			log:  nil,
		},
		{
			name: "Test with non-nil logger",
			log:  log,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			echoUtil := NewEchoUtil(tc.log)

			if echoUtil.logger != tc.log {
				t.Errorf("Expected logger to be %v, but got %v", tc.log, echoUtil.logger)
			}

			if len(echoUtil.claimKeys) != 0 {
				t.Errorf("Expected claimKeys to be an empty map, but got %v", echoUtil.claimKeys)
			}
		})
	}
}

func TestEchoUtil_CloneContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	// Initialize EchoUtil
	d := &config.CommonConfig{Logger: config.Logger{Level: "info"}}
	log := logger.NewAPILogger(d)
	log.InitLogger()
	ec := NewEchoUtil(log)

	testCases := []struct {
		name string
		ctx  echo.Context
	}{
		{
			name: "Test with empty context",
			ctx:  c,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clonedCtx := ec.CloneContext(tc.ctx)
			for claimKey := range ec.getClaimKeys() {
				if clonedCtx.Get(claimKey) != tc.ctx.Get(claimKey) {
					t.Errorf("Expected claim key %s to be %v, but got %v", claimKey, tc.ctx.Get(claimKey), clonedCtx.Get(claimKey))
				}
			}
		})
	}
}

func TestEchoUtil_setClaimKey(t *testing.T) {
	d := &config.CommonConfig{Logger: config.Logger{Level: "info"}}
	log := logger.NewAPILogger(d)
	log.InitLogger()

	ec := NewEchoUtil(log)

	testCases := []struct {
		name      string
		key       string
		expectKey bool
	}{
		{
			name:      "Test with existing key",
			key:       "LoggedIn",
			expectKey: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec.setClaimKey(tc.key)

			_, found := ec.claimKeys[tc.key]
			if found != tc.expectKey {
				t.Errorf("Expected key %s to be %v, but got %v", tc.key, tc.expectKey, found)
			}
		})
	}
}

func TestEchoUtil_initClaimKeys(t *testing.T) {
	d := &config.CommonConfig{Logger: config.Logger{Level: "info"}}
	log := logger.NewAPILogger(d)
	log.InitLogger()

	ec := NewEchoUtil(log)

	ec.initClaimKeys()

	testCases := []struct {
		name      string
		key       string
		expectKey bool
	}{
		{
			name:      "Test LoggedIn key",
			key:       LoggedIn,
			expectKey: true,
		},
		{
			name:      "Test TenantID key",
			key:       TenantID,
			expectKey: true,
		},
		{
			name:      "Test UserID key",
			key:       UserID,
			expectKey: true,
		},
		{
			name:      "Test PersonaType key",
			key:       PersonaType,
			expectKey: true,
		},
		{
			name:      "Test SessionID key",
			key:       SessionID,
			expectKey: true,
		},
		{
			name:      "Test UserExternalID key",
			key:       UserExternalID,
			expectKey: true,
		},
		{
			name:      "Test SharedDataSource key",
			key:       SharedDataSource,
			expectKey: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, found := ec.claimKeys[tc.key]
			if found != tc.expectKey {
				t.Errorf("Expected key %s to be %v, but got %v", tc.key, tc.expectKey, found)
			}
		})
	}
}

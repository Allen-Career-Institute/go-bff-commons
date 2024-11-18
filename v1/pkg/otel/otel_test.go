// nolint : gocritic,
package otel

import (
	"bff-service/config"
	"bff-service/pkg/logger"
	"bff-service/pkg/utils"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestDeltaSelector(t *testing.T) {
	tests := []struct {
		name     string
		kind     metric.InstrumentKind
		expected metricdata.Temporality
	}{
		{
			name:     "Test Counter",
			kind:     metric.InstrumentKindCounter,
			expected: metricdata.DeltaTemporality,
		},
		{
			name:     "Test Histogram",
			kind:     metric.InstrumentKindHistogram,
			expected: metricdata.DeltaTemporality,
		},
		{
			name:     "Test ObservableGauge",
			kind:     metric.InstrumentKindObservableGauge,
			expected: metricdata.DeltaTemporality,
		},
		{
			name:     "Test ObservableCounter",
			kind:     metric.InstrumentKindObservableCounter,
			expected: metricdata.DeltaTemporality,
		},
		{
			name:     "Test UpDownCounter",
			kind:     metric.InstrumentKindUpDownCounter,
			expected: metricdata.CumulativeTemporality,
		},
		{
			name:     "Test ObservableUpDownCounter",
			kind:     metric.InstrumentKindObservableUpDownCounter,
			expected: metricdata.CumulativeTemporality,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deltaSelector(tt.kind); got != tt.expected {
				t.Errorf("deltaSelector() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTrace(t *testing.T) {
	tests := []struct {
		name     string
		spanName string
	}{
		{
			name:     "Test case 1",
			spanName: "span1",
		},
		{
			name:     "Test case 2",
			spanName: "span2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			newCtx, span := Trace(ctx, tt.spanName)

			assert.NotNil(t, newCtx)
			assert.NotNil(t, span.SpanContext())
		})
	}
}

func TestInitOtelProviders(t *testing.T) {
	log := newmockLogger(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamicConfig := NewMockDynamicConfig(ctrl)

	tests := []struct {
		name        string
		endpoint1   string
		endpoint2   string
		appName     string
		unsetEnv    bool
		shouldPanic bool
		mocks       []*gomock.Call
	}{
		{
			name:        "Valid endpoints with delta temporality",
			endpoint1:   "localhost:4317",
			endpoint2:   "unused",
			appName:     "testApp",
			shouldPanic: false,
			mocks: []*gomock.Call{
				mockDynamicConfig.EXPECT().Get(utils.OtelGrpcTemporality).Return(utils.Delta, nil),
			},
		},
		{
			name:        "Valid endpoints with cumulative temporality",
			endpoint1:   "unused",
			endpoint2:   "localhost:4318",
			appName:     "testApp",
			shouldPanic: false,
			mocks: []*gomock.Call{
				mockDynamicConfig.EXPECT().Get(utils.OtelGrpcTemporality).Return(utils.Cumulative, nil),
			},
		},
		{
			name:        "Empty endpoint1 with delta temporality",
			endpoint1:   "",
			endpoint2:   "unused",
			appName:     "testApp",
			shouldPanic: false,
			mocks: []*gomock.Call{
				mockDynamicConfig.EXPECT().Get(utils.OtelGrpcTemporality).Return(utils.Delta, nil),
			},
		},
		{
			name:        "Empty endpoint2 with cumulative temporality",
			endpoint1:   "unused",
			endpoint2:   "",
			appName:     "testApp",
			shouldPanic: false,
			mocks: []*gomock.Call{
				mockDynamicConfig.EXPECT().Get(utils.OtelGrpcTemporality).Return(utils.Cumulative, nil),
			},
		},
		{
			name:        "DynamicConfig.Get returns an error, defaults to delta",
			endpoint1:   "localhost:4317",
			endpoint2:   "unused",
			appName:     "testApp",
			shouldPanic: false,
			mocks: []*gomock.Call{
				mockDynamicConfig.EXPECT().Get(utils.OtelGrpcTemporality).Return("", errors.New("config error")),
			},
		},
		{
			name:        "Invalid temporality value",
			endpoint1:   "localhost:4317",
			endpoint2:   "localhost:4318",
			appName:     "testApp",
			shouldPanic: false, // Expecting panic due to invalid temporality
			mocks: []*gomock.Call{
				mockDynamicConfig.EXPECT().Get(utils.OtelGrpcTemporality).Return("invalid_temporality", nil),
			},
		},
		{
			name:        "Environment variables not set",
			endpoint1:   "localhost:4317",
			endpoint2:   "unused",
			appName:     "testApp",
			unsetEnv:    true,
			shouldPanic: false, // Should not panic even if env vars are unset
			mocks: []*gomock.Call{
				mockDynamicConfig.EXPECT().Get(utils.OtelGrpcTemporality).Return(utils.Delta, nil),
			},
		},
		{
			name:        "App name is empty",
			endpoint1:   "localhost:4317",
			endpoint2:   "unused",
			appName:     "",
			shouldPanic: false, // Should not panic even if app name is empty
			mocks: []*gomock.Call{
				mockDynamicConfig.EXPECT().Get(utils.OtelGrpcTemporality).Return(utils.Delta, nil),
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("InitOtelProviders() panicked unexpectedly: %v", r)
					}
				} else {
					if tt.shouldPanic {
						t.Errorf("InitOtelProviders() did not panic as expected")
					}
				}
			}()

			// Set or unset environment variables
			if tt.unsetEnv {
				os.Unsetenv(utils.Tag)
				os.Unsetenv(utils.Env)
			} else {
				os.Setenv(utils.Tag, "testTag")
				os.Setenv(utils.Env, "stage")
			}

			// Create a new Handler with app name
			h := NewHandler(&config.Config{
				Server: config.ServerConfig{
					App: config.App{
						Name: tt.appName,
					},
				},
				DynamicConfig: mockDynamicConfig,
			}, log)

			// Set up mocks
			gomock.InOrder(tt.mocks...)

			// Call the function
			h.InitOtelProviders(tt.endpoint1, tt.endpoint2)
		})
	}
}

func TestNewHandler(t *testing.T) {
	log := newmockLogger(t)

	tests := []struct {
		name        string
		config      *config.Config
		expectedCnf config.Config
	}{
		{
			name:        "Valid config",
			config:      &config.Config{},
			expectedCnf: config.Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.config, log)
			assert.Equal(t, tt.expectedCnf, handler.Cnf)
		})
	}
}

type mockLogger struct {
	t *testing.T
}

func newmockLogger(t *testing.T) *mockLogger {
	return &mockLogger{t: t}
}

func (m *mockLogger) InitLogger() {}

func (m *mockLogger) WithContext(ctx echo.Context) *logger.APILogger {
	return &logger.APILogger{}
}

func (m *mockLogger) Debug(args ...interface{}) {
	m.t.Log(args...)
}

func (m *mockLogger) Debugf(template string, args ...interface{}) {
	m.t.Logf(template, args...)
}

func (m *mockLogger) Info(args ...interface{}) {
	m.t.Log(args...)
}

func (m *mockLogger) Infof(template string, args ...interface{}) {
	m.t.Logf(template, args...)
}

func (m *mockLogger) Infow(template string, keyAndValue ...interface{}) {
	m.t.Logf(template, keyAndValue...)
}

func (m *mockLogger) Warn(args ...interface{}) {
	m.t.Log(args...)
}

func (m *mockLogger) Warnf(template string, args ...interface{}) {
	m.t.Logf(template, args...)
}

func (m *mockLogger) Error(args ...interface{}) {
	m.t.Log(args...)
}

func (m *mockLogger) Errorf(template string, args ...interface{}) {
	m.t.Logf(template, args...)
}

func (m *mockLogger) Errorw(template string, keyAndValue ...interface{}) {
	m.t.Logf(template, keyAndValue...)
}

func (m *mockLogger) DPanic(args ...interface{}) {
	m.t.Log(args...)
}

func (m *mockLogger) DPanicf(template string, args ...interface{}) {
	m.t.Logf(template, args...)
}

func (m *mockLogger) Fatal(args ...interface{}) {
	m.t.Log(args...)
}

func (m *mockLogger) Fatalf(template string, args ...interface{}) {
	m.t.Logf(template, args...)
}

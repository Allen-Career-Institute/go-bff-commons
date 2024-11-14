package logger

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewAPILogger(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.CommonConfig
		want *APILogger
	}{
		{
			name: "Test with nil config",
			cfg:  nil,
			want: &APILogger{cfg: nil},
		},
		{
			name: "Test with non-nil config",
			cfg:  &config.CommonConfig{Logger: config.Logger{Level: "debug"}},
			want: &APILogger{cfg: &config.CommonConfig{Logger: config.Logger{Level: "debug"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPILogger(tt.cfg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetLoggerLevelMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string]zapcore.Level
	}{
		{
			name: "Test GetLoggerLevelMap",
			want: map[string]zapcore.Level{
				"debug":  zapcore.DebugLevel,
				"info":   zapcore.InfoLevel,
				"warn":   zapcore.WarnLevel,
				"error":  zapcore.ErrorLevel,
				"dpanic": zapcore.DPanicLevel,
				"panic":  zapcore.PanicLevel,
				"fatal":  zapcore.FatalLevel,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLoggerLevelMap()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAPILogger_getLoggerLevel(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.CommonConfig
		want zapcore.Level
	}{
		{
			name: "Test with debug level",
			cfg:  &config.CommonConfig{Logger: config.Logger{Level: "debug"}},
			want: zapcore.DebugLevel,
		},
		{
			name: "Test with info level",
			cfg:  &config.CommonConfig{Logger: config.Logger{Level: "info"}},
			want: zapcore.InfoLevel,
		},
		{
			name: "Test with warn level",
			cfg:  &config.CommonConfig{Logger: config.Logger{Level: "warn"}},
			want: zapcore.WarnLevel,
		},
		{
			name: "Test with error level",
			cfg:  &config.CommonConfig{Logger: config.Logger{Level: "error"}},
			want: zapcore.ErrorLevel,
		},
		{
			name: "Test with non-existent level",
			cfg:  &config.CommonConfig{Logger: config.Logger{Level: "non-existent"}},
			want: zapcore.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewAPILogger(tt.cfg)
			if got := logger.getLoggerLevel(tt.cfg); got != tt.want {
				t.Errorf("APILogger.getLoggerLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestNewLogEncoderConfig(t *testing.T) {
//	tests := []struct {
//		name               string
//		expectedEncoderCfg zapcore.EncoderConfig
//	}{
//		{
//			name: "Valid encoder config",
//			expectedEncoderCfg: zapcore.EncoderConfig{
//				TimeKey:        "TIME",
//				LevelKey:       "LEVEL",
//				NameKey:        "NAME",
//				CallerKey:      "CALLER",
//				MessageKey:     "MESSAGE",
//				FunctionKey:    zapcore.OmitKey,
//				StacktraceKey:  "stacktrace",
//				LineEnding:     zapcore.DefaultLineEnding,
//				EncodeLevel:    zapcore.LowercaseLevelEncoder,
//				EncodeTime:     zapcore.TimeEncoderOfLayout(TimeFormatMs),
//				EncodeDuration: zapcore.SecondsDurationEncoder,
//				EncodeCaller:   zapcore.ShortCallerEncoder,
//			},
//		},
//		// Add more test cases as needed
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			result := NewLogEncoderConfig()
//			assert.Equal(t, tt.expectedEncoderCfg, result, "Encoder config does not match expected")
//		})
//	}
//}

func TestAPILogger_InitLogger(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.CommonConfig
		wantLogs []observer.LoggedEntry
	}{
		{
			name:     "Test with console encoding",
			cfg:      &config.CommonConfig{Logger: config.Logger{Level: "debug", Encoding: "console"}},
			wantLogs: []observer.LoggedEntry{},
		},
		{
			name:     "Test with json encoding",
			cfg:      &config.CommonConfig{Logger: config.Logger{Level: "debug", Encoding: "json"}},
			wantLogs: []observer.LoggedEntry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.DebugLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				cfg:         tt.cfg,
				sugarLogger: logger.Sugar(),
			}

			apiLogger.InitLogger()

			logs := recorded.All()
			if len(logs) != len(tt.wantLogs) {
				t.Errorf("Unexpected number of logs. got=%d, want=%d", len(logs), len(tt.wantLogs))
			}

			for i, log := range logs {
				if log.Message != tt.wantLogs[i].Message {
					t.Errorf("Unexpected log message. got=%s, want=%s", log.Message, tt.wantLogs[i].Message)
				}
			}
		})
	}
}

func TestAPILogger_WithContext(t *testing.T) {
	tests := []struct {
		name string
		ctx  echo.Context
		want *APILogger
	}{
		{
			name: "Test with non-nil context",
			ctx:  echo.New().AcquireContext(),
			want: &APILogger{
				cfg:         nil,
				sugarLogger: nil,
				ctx:         echo.New().AcquireContext(),
			},
		},
		{
			name: "Test with nil context",
			ctx:  nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewAPILogger(nil)
			got := logger.WithContext(tt.ctx)

			if got == nil {
				assert.Equal(t, got, tt.want)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestAPILogger_Debug(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
	}{
		{
			name: "Test with no arguments",
			args: []interface{}{},
		},
		{
			name: "Test with one argument",
			args: []interface{}{"test"},
		},
		{
			name: "Test with multiple arguments",
			args: []interface{}{"test", 123, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, _ := observer.New(zap.DebugLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				cfg:         &config.CommonConfig{},
				sugarLogger: logger.Sugar(),
			}

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Debug() panics with args: %v", tt.args)
				}
			}()

			apiLogger.Debug(tt.args...)
		})
	}
}

func TestAPILogger_Debugf(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []interface{}
	}{
		{
			name:     "Test with no arguments",
			template: "test message",
			args:     []interface{}{},
		},
		{
			name:     "Test with one argument",
			template: "test message %v",
			args:     []interface{}{"arg1"},
		},
		{
			name:     "Test with multiple arguments",
			template: "test message %v %v",
			args:     []interface{}{"arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, _ := observer.New(zap.DebugLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				cfg:         &config.CommonConfig{},
				sugarLogger: logger.Sugar(),
			}

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Debugf() panics with args: %v", tt.args)
				}
			}()

			apiLogger.Debugf(tt.template, tt.args...)
		})
	}
}

func TestAPILogger_Info(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{
			name: "Test with no arguments",
			args: []interface{}{},
			want: "",
		},
		{
			name: "Test with one argument",
			args: []interface{}{"test"},
			want: "test",
		},
		{
			name: "Test with multiple arguments",
			args: []interface{}{"test", 123, true},
			want: "test123 true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				sugarLogger: logger.Sugar(),
			}

			apiLogger.Info(tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Info_WithContext(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{
			name: "Test with no args and with context",
			args: []interface{}{},
			want: "",
		},
		{
			name: "Test with args and with context",
			args: []interface{}{"test"},
			want: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			rec.Header().Set(echo.HeaderXRequestID, "test-request-id")
			req.RemoteAddr = "test-ip"

			spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    trace.TraceID{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef},
				SpanID:     trace.SpanID{0x12, 0x34, 0x56, 0x78},
				TraceFlags: trace.FlagsSampled,
			})

			ctx.SetRequest(req.WithContext(trace.ContextWithSpanContext(req.Context(), spanCtx)))

			ctx.Set("userID", "allen-student")

			apiLogger := &APILogger{
				sugarLogger: logger.Sugar(),
				ctx:         ctx,
			}

			apiLogger.Info(tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Infof(t *testing.T) {
	tests := []struct {
		name     string
		ctx      echo.Context
		template string
		args     []interface{}
		want     string
	}{
		{
			name:     "Test with no arguments and nil context",
			ctx:      nil,
			template: "test message",
			args:     []interface{}{},
			want:     "test message",
		},
		{
			name:     "Test with arguments and nil context",
			ctx:      nil,
			template: "test message %v",
			args:     []interface{}{"arg1"},
			want:     "test message arg1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				cfg:         &config.CommonConfig{},
				sugarLogger: logger.Sugar(),
				ctx:         tt.ctx,
			}

			apiLogger.Infof(tt.template, tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Infow(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		keyAndValue []interface{}
	}{
		{
			name:        "Test with no key-value pairs",
			template:    "test message",
			keyAndValue: []interface{}{},
		},
		{
			name:        "Test with one key-value pair",
			template:    "test message",
			keyAndValue: []interface{}{"key1", "value1"},
		},
		{
			name:        "Test with multiple key-value pairs",
			template:    "test message",
			keyAndValue: []interface{}{"key1", "value1", "key2", "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, _ := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				sugarLogger: logger.Sugar(),
			}

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Infow() panics with keyAndValue: %v", tt.keyAndValue)
				}
			}()

			apiLogger.Infow(tt.template, tt.keyAndValue...)
		})
	}
}

func TestAPILogger_Warn(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{
			name: "Test with no arguments",
			args: []interface{}{},
			want: "",
		},
		{
			name: "Test with one argument",
			args: []interface{}{"test"},
			want: "test",
		},
		{
			name: "Test with multiple arguments",
			args: []interface{}{"test", 123, true},
			want: "test123 true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.WarnLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				sugarLogger: logger.Sugar(),
			}

			apiLogger.Warn(tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Warnf(t *testing.T) {
	tests := []struct {
		name     string
		ctx      echo.Context
		template string
		args     []interface{}
		want     string
	}{
		{
			name:     "Test with no arguments and nil context",
			ctx:      nil,
			template: "test message",
			args:     []interface{}{},
			want:     "test message",
		},
		{
			name:     "Test with arguments and nil context",
			ctx:      nil,
			template: "test message %v",
			args:     []interface{}{"arg1"},
			want:     "test message arg1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.WarnLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				cfg:         &config.CommonConfig{},
				sugarLogger: logger.Sugar(),
				ctx:         tt.ctx,
			}

			apiLogger.Warnf(tt.template, tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Error(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{
			name: "Test with no arguments",
			args: []interface{}{},
			want: "",
		},
		{
			name: "Test with one argument",
			args: []interface{}{"test"},
			want: "test",
		},
		{
			name: "Test with multiple arguments",
			args: []interface{}{"test", 123, true},
			want: "test123 true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.ErrorLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				sugarLogger: logger.Sugar(),
			}

			apiLogger.Error(tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Errorf(t *testing.T) {
	tests := []struct {
		name     string
		ctx      echo.Context
		template string
		args     []interface{}
		want     string
		init     func() echo.Context
	}{
		{
			name:     "Test with no arguments and nil context",
			ctx:      nil,
			template: "test message",
			args:     []interface{}{},
			want:     "test message",
			init:     func() echo.Context { return nil },
		},
		{
			name:     "Test with arguments and nil context",
			ctx:      nil,
			template: "test message %v",
			args:     []interface{}{"arg1"},
			want:     "test message arg1",
			init:     func() echo.Context { return nil },
		},
		{
			name:     "Test with empty context, non-empty template and nil args",
			ctx:      nil,
			template: "test message",
			args:     nil,
			want:     "test message",
			init:     func() echo.Context { return nil },
		},
		{
			name:     "Test with non-empty template and context",
			ctx:      nil,
			template: "test message",
			args:     nil,
			want:     "test message",
			init: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				return e.NewContext(req, rec)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.ErrorLevel)
			logger := zap.New(core)

			// Initialize the context
			ctx := tt.init()

			apiLogger := &APILogger{
				cfg:         &config.CommonConfig{},
				sugarLogger: logger.Sugar(),
				ctx:         ctx,
			}

			apiLogger.Errorf(tt.template, tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Errorw(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		keyAndValue []interface{}
		want        string
	}{
		{
			name:        "Test with no additional context",
			template:    "test message",
			keyAndValue: []interface{}{},
			want:        "test message",
		},
		{
			name:        "Test with additional context",
			template:    "test message",
			keyAndValue: []interface{}{"key1", "value1"},
			want:        "test message",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.ErrorLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				cfg:         &config.CommonConfig{},
				sugarLogger: logger.Sugar(),
			}

			apiLogger.Errorw(tt.template, tt.keyAndValue...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_DPanic(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{
			name: "Test with no arguments",
			args: []interface{}{},
			want: "",
		},
		{
			name: "Test with one argument",
			args: []interface{}{"test"},
			want: "test",
		},
		{
			name: "Test with multiple arguments",
			args: []interface{}{"test", 123, true},
			want: "test123 true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.DPanicLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				sugarLogger: logger.Sugar(),
			}

			apiLogger.DPanic(tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Panic(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{
			name: "Test with no arguments",
			args: []interface{}{},
			want: "",
		},
		{
			name: "Test with one argument",
			args: []interface{}{"test"},
			want: "test",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.PanicLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				cfg:         &config.CommonConfig{},
				sugarLogger: logger.Sugar(),
			}

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Panic() did not panic")
				}
			}()

			apiLogger.Panic(tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

func TestAPILogger_Panicf(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []interface{}
		want     string
	}{
		{
			name:     "Test with no arguments",
			template: "test message",
			args:     []interface{}{},
			want:     "test message",
		},
		{
			name:     "Test with one argument",
			template: "test message %v",
			args:     []interface{}{"arg1"},
			want:     "test message arg1",
		},
		{
			name:     "Test with multiple arguments",
			template: "test message %v %v",
			args:     []interface{}{"arg1", "arg2"},
			want:     "test message arg1 arg2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.PanicLevel)
			logger := zap.New(core)

			apiLogger := &APILogger{
				sugarLogger: logger.Sugar(),
			}

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Panicf() did not panic")
				}
			}()

			apiLogger.Panicf(tt.template, tt.args...)

			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Unexpected number of logs. got=%d, want=1", len(logs))
			}

			if logs[0].Message != tt.want {
				t.Errorf("Unexpected log message. got=%s, want=%s", logs[0].Message, tt.want)
			}
		})
	}
}

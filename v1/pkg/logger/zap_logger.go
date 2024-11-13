package logger

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
)

const (
	TimeFormatMs = "2006-01-02 15:04:05.000"
	userID       = "userID"
)

// Logger methods interface
type Logger interface {
	InitLogger()
	WithContext(ctx echo.Context) *APILogger
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(template string, keyAndValue ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(template string, keyAndValue ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

// ApiLogger Logger
type APILogger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
	ctx         echo.Context
}

// NewAPILogger App Logger constructor
func NewAPILogger(cfg *config.Config) *APILogger {
	return &APILogger{cfg: cfg}
}

// GetLoggerLevelMap is for mapping config logger to app logger levels
func GetLoggerLevelMap() map[string]zapcore.Level {
	return map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
}

func (*APILogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := GetLoggerLevelMap()[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func NewLogEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(TimeFormatMs),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// InitLogger Init logger
func (l *APILogger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)

	logWriter := zapcore.AddSync(os.Stderr)

	encoderCfg := NewLogEncoderConfig()

	var encoder zapcore.Encoder
	if l.cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
}

func (l *APILogger) WithContext(ctx echo.Context) *APILogger {
	if ctx != nil {
		return &APILogger{
			cfg:         l.cfg,
			sugarLogger: l.sugarLogger,
			ctx:         ctx,
		}
	}

	return nil
}

// Logger methods

func (l *APILogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *APILogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *APILogger) Info(args ...interface{}) {
	if l.ctx == nil {
		l.sugarLogger.Info(args...)
		return
	}

	template, args := l.getCustomMessage("", args)
	l.sugarLogger.Infow(template, args...)
}

func (l *APILogger) Infof(template string, args ...interface{}) {
	if l.ctx != nil {
		template, args = l.getCustomMessage(template, args)
		l.sugarLogger.Infow(template, args...)
	} else {
		l.sugarLogger.Infof(template, args...)
	}
}

func (l *APILogger) Infow(template string, keyAndValue ...interface{}) {
	keyAndValue = append(keyAndValue, l.getAdditionalContextArgs()...)
	l.sugarLogger.Infow(template, keyAndValue...)
}

func (l *APILogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *APILogger) Warnf(template string, args ...interface{}) {
	if l.ctx != nil {
		template, args = l.getCustomMessage(template, args)
		l.sugarLogger.Warnw(template, args...)
	} else {
		l.sugarLogger.Warnf(template, args...)
	}
}

func (l *APILogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *APILogger) Errorf(template string, args ...interface{}) {
	if l.ctx != nil {
		template, args = l.getCustomMessage(template, args)
		l.sugarLogger.Errorw(template, args...)
	} else {
		l.sugarLogger.Errorf(template, args...)
	}
}

func (l *APILogger) Errorw(template string, keyAndValue ...interface{}) {
	keyAndValue = append(keyAndValue, l.getAdditionalContextArgs()...)
	l.sugarLogger.Errorw(template, keyAndValue...)
}

func (l *APILogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *APILogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *APILogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *APILogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *APILogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *APILogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

func (l *APILogger) getCustomMessage(template string, args []interface{}) (formattedTemplate string, additionalArgs []interface{}) {
	if args != nil {
		if template == "" {
			formattedTemplate = fmt.Sprint(args...)
		} else {
			formattedTemplate = fmt.Sprintf(template, args...)
		}
	} else {
		formattedTemplate = template
	}

	additionalArgs = l.getAdditionalContextArgs()

	return formattedTemplate, additionalArgs
}

func (l *APILogger) getAdditionalContextArgs() (args []interface{}) {
	if l.ctx == nil {
		return nil
	}

	const (
		traceIDLength = 32
		spanIDLength  = 16
	)

	ipAdd := l.ctx.Request().RemoteAddr
	requestID := l.ctx.Response().Header().Get(echo.HeaderXRequestID)
	traceID := trace.SpanFromContext(l.ctx.Request().Context()).SpanContext().TraceID().String()
	spanID := trace.SpanFromContext(l.ctx.Request().Context()).SpanContext().SpanID().String()
	args = []interface{}{"requestID", requestID, "ip", ipAdd}

	nullTraceID := fmt.Sprintf("%0*s", traceIDLength, "")
	nullSpanID := fmt.Sprintf("%0*s", spanIDLength, "")

	if nullTraceID != traceID {
		args = append(args, "dd.trace_id", traceID)
	}

	if nullSpanID != spanID {
		args = append(args, "dd.span_id", spanID)
	}

	uid := l.ctx.Get(userID)
	if uid != "" {
		args = append(args, "userID", uid)
	}

	return args
}

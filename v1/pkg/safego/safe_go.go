package safego

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/labstack/echo/v4"
	"runtime/debug"
)

// safeGo wraps a function to run as a goroutine with panic recovery and logging.
// It takes an optional Echo context and custom logger to provide detailed logging if a panic occurs.
func safeGo(fn func(), customLogger logger.Logger, ctx echo.Context) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := debug.Stack()
				if customLogger != nil {
					customLogger.WithContext(ctx).Errorf("Recovered_from_panic: %v\n", r)
					customLogger.WithContext(ctx).Errorf("Stack_trace: %v\n", string(stackTrace))
				} else {
					log.NewHelper(log.DefaultLogger).Errorf("Recovered_from_panic: %v\n", r)
					log.NewHelper(log.DefaultLogger).Errorf("Stack_trace: %v\n", string(stackTrace))
				}
			}
			//handlePanic(ctx, customLogger)
		}()
		fn()
	}()
}

// SafeGoWithLogger runs a function in a goroutine with panic recovery and logging using the provided custom logger and context.
func SafeGoWithLogger(fn func(), customLogger logger.Logger, ctx ...echo.Context) {
	var context echo.Context
	if len(ctx) > 0 {
		context = ctx[0]
	} else {
		context = nil
	}
	safeGo(fn, customLogger, context)
}

// SafeGo runs a function in a goroutine with panic recovery and logging using the default logger and provided context.
func SafeGo(fn func(), ctx ...echo.Context) {
	var context echo.Context
	if len(ctx) > 0 {
		context = ctx[0]
	} else {
		context = nil
	}
	safeGo(fn, nil, context)
}

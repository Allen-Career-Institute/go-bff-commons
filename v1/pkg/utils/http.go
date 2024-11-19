package utils

import (
	"context"
	models "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models/commons"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/datasources"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

// GetRequestID Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func HandleError(c echo.Context, err error, log logger.Logger) (code int, message string) {
	log.WithContext(c).Errorf("error received, Err: %v", err)

	e, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, e.Message()
	}

	code = GetHTTPStatusCode(e.Code())
	message = e.Message()

	return code, message
}

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// GetCtxWithReqID Get ctx with timeout and request id from echo context
func GetCtxWithReqID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*15)
	ctx = context.WithValue(ctx, ReqIDCtxKey{}, GetRequestID(c))

	return ctx, cancel
}

func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// GetRequestCtxWithTimeout Get context  with request id
func GetRequestCtxWithTimeout(c echo.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	// timeout should be in milliseconds
	ctx := context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
	ctx, connCancel := context.WithTimeout(ctx, timeout)

	return ctx, connCancel
}

func AddAuthHeaderAsMetadata(ctx context.Context, c echo.Context) context.Context {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	platform := c.Request().Header.Get(DeviceType)
	deviceID := c.Request().Header.Get(DeviceID)
	visitorID := c.Request().Header.Get(VisitorID)

	if auth != EmptyString {
		ctx = metadata.AppendToOutgoingContext(ctx, echo.HeaderAuthorization, auth)
		ctx = metadata.AppendToOutgoingContext(ctx, echo.HeaderXRequestID, GetRequestID(c))
	}
	ctx = metadata.AppendToOutgoingContext(ctx, DeviceID, deviceID)
	ctx = metadata.AppendToOutgoingContext(ctx, DeviceType, platform)
	ctx = metadata.AppendToOutgoingContext(ctx, VisitorID, visitorID)
	ctx = metadata.AppendToOutgoingContext(ctx, ServiceHeader, TracerServiceName)

	return ctx
}

// GetRequestCtxWithClientIDAndTimeout lmm takes clientRequestID field in header, and add to the event
// that vod listens. We need to add this for vod pipeline to trigger and add the recording to the meeting
func GetRequestCtxWithClientIDAndTimeout(c echo.Context, timeout time.Duration, clientRequestID string) (
	context.Context, context.CancelFunc) {
	ctx := context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
	ctx, connCancel := context.WithTimeout(ctx, timeout)
	ctx = metadata.AppendToOutgoingContext(ctx, ClientRequestID, clientRequestID)

	return ctx, connCancel
}

// DeleteSessionCookie Delete session
func DeleteSessionCookie(c echo.Context, sessionName string) {
	c.SetCookie(&http.Cookie{
		Name:   sessionName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

// UserCtxKey is a key used for the OtpRequest object in the context
type UserCtxKey struct{}

// GetIPAddress Get user ip address
func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

// ReadRequest Read request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}

	validate := validator.New()

	return validate.StructCtx(ctx.Request().Context(), request)
}

func getAllowedImagesContentTypes() map[string]string {
	return map[string]string{
		"image/bmp":                "bmp",
		"image/gif":                "gif",
		"image/png":                "png",
		"image/jpeg":               "jpeg",
		"image/jpg":                "jpg",
		"image/svg+xml":            "svg",
		"image/webp":               "webp",
		"image/tiff":               "tiff",
		"image/vnd.microsoft.icon": "ico",
	}
}

func CheckImageFileContentType(fileContent []byte) (string, error) {
	contentType := http.DetectContentType(fileContent)

	extension, ok := getAllowedImagesContentTypes()[contentType]
	if !ok {
		return "", errors.New("this content type is not allowed")
	}

	return extension, nil
}

func GetHTTPStatusCode(grpcStatusCode codes.Code) int {
	codeMapping := map[codes.Code]int{
		codes.OK:                 http.StatusOK,
		codes.Canceled:           449, // ClientClosed
		codes.Unknown:            http.StatusInternalServerError,
		codes.InvalidArgument:    http.StatusBadRequest,
		codes.DeadlineExceeded:   http.StatusGatewayTimeout,
		codes.NotFound:           http.StatusNotFound,
		codes.AlreadyExists:      http.StatusConflict,
		codes.PermissionDenied:   http.StatusForbidden,
		codes.Unauthenticated:    http.StatusUnauthorized,
		codes.ResourceExhausted:  http.StatusTooManyRequests,
		codes.FailedPrecondition: http.StatusBadRequest,
		codes.Aborted:            http.StatusConflict,
		codes.OutOfRange:         http.StatusBadRequest,
		codes.Unimplemented:      http.StatusNotImplemented,
		codes.Internal:           http.StatusInternalServerError,
		codes.Unavailable:        http.StatusServiceUnavailable,
		codes.DataLoss:           http.StatusInternalServerError,
	}

	if statusCode, ok := codeMapping[grpcStatusCode]; ok {
		return statusCode
	}

	return http.StatusInternalServerError
}

func GetErrorCodeAndMessage(err error) (respCode int, errMessage string) {
	e, _ := status.FromError(err)
	errMessage = e.Message()
	respCode = GetHTTPStatusCode(e.Code())

	return respCode, errMessage
}

func GetHeader(ctx echo.Context, name string) string {
	if ctx.Request().Header != nil {
		return strings.TrimSpace(ctx.Request().Header.Get(name))
	}

	return EmptyString
}

func HandleErrorAndConvertResponse(errorResponse error) (models.DSResponse, error) {
	e, _ := status.FromError(errorResponse)
	errMessage := e.Message()
	respCode := GetHTTPStatusCode(e.Code())

	if respCode == http.StatusInternalServerError {
		errMessage = GenericError
	}

	return datasources.PopulateResponse(respCode, errMessage, errorResponse.Error()), nil
}

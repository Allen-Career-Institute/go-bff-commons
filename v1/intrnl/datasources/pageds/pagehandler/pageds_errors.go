package pagehandler

import (
	"github.com/go-kratos/kratos/v2/errors"
)

const (
	ErrorMessageNoTabsToShow          = "This page does not have any tabs"
	ErrorMessageUnsupportedPageType   = "Page Type is not supported"
	ErrorMessageFechingEnrolledStatus = "Error while fetching user enrollment data"
)

const (
	ErrorReasonNoTabsToShow          = "No Tabs to show"
	ErrorReasonUnsupportedPageType   = "Unsupported page type"
	ErrorReasonFechingEnrolledStatus = "No user enrollment data"
)

var (
	ErrorNoTabsToShow          = errors.InternalServer(ErrorReasonNoTabsToShow, ErrorMessageNoTabsToShow)
	ErrorUnsupportedPageType   = errors.InternalServer(ErrorMessageUnsupportedPageType, ErrorReasonUnsupportedPageType)
	ErrorFechingEnrolledStatus = errors.InternalServer(ErrorMessageFechingEnrolledStatus, ErrorReasonFechingEnrolledStatus)
)

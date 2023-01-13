package ecode

import "net/http"

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = http.StatusInternalServerError
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// ClientClosed is non-standard http status code,
	// which defined by nginx.
	// https://httpstatus.in/499/
	ClientClosed = 499
)

var (
	// base error
	Success               = NewV2(1, "SUCCESS", "success")
	RequestErr            = NewV2(40000, "PARAM_ERROR", "request param error")
	UnauthorizedErr       = NewV2(40001, "SIGN_ERROR", "sign error")
	ForbiddenErr          = NewV2(40003, "NO_AUTH", "no auth")
	NotFoundErr           = NewV2(40004, "RESOURCE_NOT_FOUND", "resource not found")
	TooManyRequestErr     = NewV2(40029, "RATELIMIT_EXCEEDED", "ratelimit exceeded")
	ServerErr             = NewV2(50000, "SERVER_ERROR", "server error")
	ServiceUnavailableErr = NewV2(50003, "SERVICE_UNAVAILABLE", "service protected, unavailable")
)

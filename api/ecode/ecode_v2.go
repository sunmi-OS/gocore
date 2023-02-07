package ecode

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode     = http.StatusInternalServerError
	SystemErrorCode = -1
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// ClientClosed is non-standard http status code,
	// which defined by nginx.
	// https://httpstatus.in/499/
	ClientClosed = 499
)

var (
	// base error
	Success     = NewV2(1, "success")
	SystemError = NewV2(SystemErrorCode, "system error")
)

// ErrorV2 struct
type ErrorV2 struct {
	Status
	cause error
}

func (e *ErrorV2) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v cause = %v", e.Status.Code, e.Status.Reason, e.Status.Message, e.Metadata, e.cause)
}

// Code returns the code of the error.
func (e *ErrorV2) Code() int { return int(e.Status.Code) }

// Message returns the message of the error.
func (e *ErrorV2) Message() string { return e.Status.Message }

// Reason returns the reason of the error.
func (e *ErrorV2) Reason() string { return e.Status.Reason }

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *ErrorV2) Unwrap() error { return e.cause }

// Is matches each error in the chain with the target value.
func (e *ErrorV2) Is(err error) bool {
	if se := new(ErrorV2); errors.As(err, &se) {
		return se.Status.Code == e.Status.Code && se.Status.Reason == e.Status.Reason
	}
	return false
}

// Equal matches error from code and reason.
func (e *ErrorV2) Equal(code int, reason ...string) bool {
	se := &ErrorV2{Status: Status{
		Code: int32(code),
	}}
	if len(reason) == 1 {
		se.Status.Reason = reason[0]
		return se.Status.Code == e.Status.Code && se.Status.Reason == e.Status.Reason
	}
	return se.Status.Code == e.Status.Code
}

// GRPCStatus returns the Status represented by error.
func (e *ErrorV2) GRPCStatus() *status.Status {
	gs, _ := status.New(DefaultConverter.ToGRPCCode(int(e.Status.Code)), e.Status.Message).
		WithDetails(&errdetails.ErrorInfo{Reason: e.Status.Reason, Metadata: e.Metadata})
	return gs
}

// WithCause with the underlying cause of the error.
func (e *ErrorV2) WithCause(cause error) *ErrorV2 {
	err := DeepClone(e)
	err.cause = cause
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *ErrorV2) WithMetadata(md map[string]string) *ErrorV2 {
	err := DeepClone(e)
	err.Metadata = md
	return err
}

// ============================================================================================================

// NewV2 returns an error object for the code, message.
func NewV2(code int, message string) *ErrorV2 {
	return &ErrorV2{Status: Status{
		Code:    int32(code),
		Message: message,
	}}
}

// NewV2WithReason returns an error object for the code, message.
func NewV2WithReason(code int, reason string, message string) *ErrorV2 {
	return &ErrorV2{Status: Status{
		Code:    int32(code),
		Reason:  reason,
		Message: message,
	}}
}

// DeepClone deep clone error to a new error.
func DeepClone(err *ErrorV2) *ErrorV2 {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &ErrorV2{
		cause: err.cause,
		Status: Status{
			Code:     err.Status.Code,
			Reason:   err.Status.Reason,
			Message:  err.Status.Message,
			Metadata: metadata,
		},
	}
}

// FromError try to convert an error to *ErrorV2.
// It supports wrapped errors.
func FromError(err error) *ErrorV2 {
	if err == nil {
		return Success
	}
	if se := new(ErrorV2); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return NewV2(SystemErrorCode, err.Error())
	}
	ret := NewV2(DefaultConverter.FromGRPCCode(gs.Code()), gs.Message())
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			ret.Status.Reason = d.Reason
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}

// analyse error info
func AnalyseErrorV2(err error) *ErrorV2 {
	if err == nil {
		return Success
	}
	if e, ok := err.(*ErrorV2); ok {
		return e
	}
	return errStringToErrorV2(err.Error())
}

func errStringToErrorV2(e string) *ErrorV2 {
	if e == "" {
		return Success
	}
	i, err := strconv.Atoi(e)
	if err != nil {
		return NewV2(-1, e)
	}
	return NewV2(i, e)
}

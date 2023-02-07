package ecode

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type statusConverter struct{}

// DefaultConverter default converter.
var DefaultConverter = statusConverter{}

// ToGRPCCode converts an HTTP error code into the corresponding gRPC response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func (c statusConverter) ToGRPCCode(code int) codes.Code {
	switch code {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.Aborted
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case ClientClosed:
		return codes.Canceled
	default:
		if (code > 0 && code <= 99) || code >= 599 {
			// code not http status code
			return codes.Code(code)
		}
		return codes.Unknown
	}
}

// FromGRPCCode converts a gRPC error code into the corresponding HTTP response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func (c statusConverter) FromGRPCCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return ClientClosed
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		if (code > 0 && code <= 99) || code >= 599 {
			// code not http status code
			return int(code)
		}
		return http.StatusInternalServerError
	}
}

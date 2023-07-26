package grpc

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetGRPCStatus(err error) *status.Status {
	grpcErr, ok := status.FromError(err)
	if ok {
		return grpcErr
	}
	return status.New(codes.Unknown, "")
}

func GetHttpStatusCodeFromGrpc(err error) int {
	code := GetGRPCStatus(err).Code()
	return grpcCodeToHTTPStatus(code)
}

func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK

	case codes.Canceled:
		return http.StatusRequestTimeout

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
		return http.StatusInternalServerError
	}
}

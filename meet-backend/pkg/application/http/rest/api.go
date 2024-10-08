package rest

import "github.com/erodriguezg/meet/pkg/core/exception"

const (
	ApiStatusOk            = "OK"
	ApiStatusBusinessError = "BUSINESS_ERROR"
	ApiStatusError         = "ERROR"
	ApiStatusAccessDenied  = "ACCESS_DENIED"
)

type ApiResponse[T any] struct {
	Status  string          `json:"status"`
	Payload *T              `json:"payload,omitempty"`
	Error   *ApiErrorDetail `json:"error,omitempty"`
}

type ApiErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

func ApiOk[T any](payload *T) ApiResponse[T] {
	return ApiResponse[T]{
		Status:  ApiStatusOk,
		Payload: payload,
		Error:   nil,
	}
}

func ApiOkArray[T any](payload []T) ApiResponse[[]T] {
	return ApiResponse[[]T]{
		Status:  ApiStatusOk,
		Payload: &payload,
		Error:   nil,
	}
}

func ApiOkEmpty() ApiResponse[string] {
	return ApiResponse[string]{
		Status:  ApiStatusOk,
		Payload: nil,
		Error:   nil,
	}
}

func ApiAccessDenied() ApiResponse[any] {
	return ApiResponse[any]{
		Status:  ApiStatusAccessDenied,
		Payload: nil,
		Error:   nil,
	}
}

func ApiBusinessException(exception *exception.BusinessException) ApiResponse[any] {
	return ApiResponse[any]{
		Status:  ApiStatusBusinessError,
		Payload: nil,
		Error: &ApiErrorDetail{
			Code:    exception.Code,
			Message: exception.Message,
			Details: exception.Details,
		},
	}
}

func ApiError(err error) ApiResponse[any] {
	return ApiResponse[any]{
		Status:  ApiStatusError,
		Payload: nil,
		Error: &ApiErrorDetail{
			Code:    "error",
			Message: err.Error(),
			Details: nil,
		},
	}
}

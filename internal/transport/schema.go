package transport

import "github.com/danielgtaylor/huma/v2"

type response[T any] struct {
	Message string `json:"message"`
	Data    *T     `json:"data"`
	Success bool   `json:"success"`
}

type BaseResponse[T any] struct {
	Status int
	Body   response[T]
}

type ErrorSchema struct {
	huma.ErrorModel
	Success bool `json:"success"`
}

func SuccessResponse[T any](statusCode int, data *T, message string) *BaseResponse[T] {
	return &BaseResponse[T]{
		Status: statusCode,
		Body: response[T]{
			Message: message,
			Data:    data,
			Success: true,
		},
	}
}

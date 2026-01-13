package httpx

import "rest-fiber/pkg/pagination"

type HttpPaginationResponse struct {
	HttpResponse
	Meta pagination.Meta `json:"meta,omitempty"`
}

func NewHttpPaginationResponse[T any](statusCode int, message string, data T, meta pagination.Meta) HttpPaginationResponse {
	return HttpPaginationResponse{
		HttpResponse{
			StatusCode: statusCode,
			Message:    message,
			Data:       data,
		},
		meta,
	}
}

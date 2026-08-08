package transport

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Service status"`
	}
}

func registerHealthRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Summary:     "Health check",
		Description: "Returns the current health status of the gateway.",
		Method:      http.MethodGet,
		Path:        "api/health",
		Tags:        []string{"Health"},
	}, func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})
}

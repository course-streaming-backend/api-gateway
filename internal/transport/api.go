package transport

import (
	"net/http"
	"os"

	authv1 "github.com/course-streaming-backend/api-gateway/gen/auth/v1"
	"github.com/danielgtaylor/huma/v2"
	grpcpb "google.golang.org/grpc"
)

type Api struct {
	api  huma.API
	auth authv1.AuthServiceClient
}

func InitAPI(api huma.API, conn *grpcpb.ClientConn) *Api {
	a := &Api{
		api:  api,
		auth: authv1.NewAuthServiceClient(conn),
	}

	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		response := &ErrorSchema{
			ErrorModel: huma.ErrorModel{
				Title:  http.StatusText(status),
				Status: status,
				Detail: message,
			},
			Success: false,
		}
		// show internal app errors only in dev
		if len(errs) > 0 && errs[0] != nil && (os.Getenv("LOCAL") == "true" || os.Getenv("IS_DEV") == "true") {
			response.Add(errs[0])
		}
		return response
	}
	a.registerRoutes()
	return a
}

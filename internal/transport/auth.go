package transport

import (
	"context"
	"net/http"

	authv1 "github.com/course-streaming-backend/api-gateway/gen/auth/v1"
	"github.com/danielgtaylor/huma/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// should probably extact the handlers

type AuthHealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Auth service status"`
		DbOk   bool   `json:"db_ok" example:"true" doc:"Whether the auth service can reach its database"`
	}
}

type RegisterInput struct {
	Body struct {
		Email    string  `json:"email" example:"user@example.com" doc:"User email" minLength:"3" maxLength:"320"`
		Username string  `json:"username" example:"johndoe" doc:"Desired username" minLength:"3" maxLength:"32"`
		Password string  `json:"password" example:"secret123" doc:"Password" minLength:"8" maxLength:"128"`
		Phone    *string `json:"phone,omitempty" example:"+15551234567" doc:"Optional phone number"`
	}
}

type RegisterOutput struct {
	Body struct {
		UserId   string `json:"user_id" doc:"Id of the created user"`
		Email    string `json:"email" doc:"User email"`
		Username string `json:"username" doc:"User username"`
	}
}

type LoginInput struct {
	Body struct {
		Identifier string `json:"identifier" example:"user@example.com" doc:"Email or username" minLength:"3" maxLength:"320"`
		Password   string `json:"password" example:"secret123" doc:"Password" minLength:"8" maxLength:"128"`
	}
}

type LoginOutput struct {
	Body struct {
		AccessToken string `json:"access_token" doc:"Access token"`
		TokenType   string `json:"token_type" doc:"Token type" example:"Bearer"`
		ExpiresIn   int64  `json:"expires_in" doc:"Token lifetime in seconds"`
		UserId      string `json:"user_id" doc:"Authenticated user id"`
	}
}

func registerAuthRoutes(api huma.API, auth authv1.AuthServiceClient) {
	huma.Register(api, huma.Operation{
		OperationID: "authHealth",
		Summary:     "Auth service health",
		Description: "Forwards to the auth service over gRPC and returns its health.",
		Method:      http.MethodGet,
		Path:        "api/auth/health",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, input *struct{}) (*BaseResponse[AuthHealthOutput], error) {
		resp, err := auth.Health(ctx, &authv1.HealthRequest{})
		if err != nil {
			return nil, grpcError(err)
		}
		out := &AuthHealthOutput{}
		out.Body.Status = resp.GetStatus()
		out.Body.DbOk = resp.GetDbOk()
		return SuccessResponse(200, out, "Status is available"), nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "register",
		Summary:     "Register a new account",
		Description: "Forwards to the auth service over gRPC.",
		Method:      http.MethodPost,
		Path:        "api/auth/register",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, input *RegisterInput) (*BaseResponse[RegisterOutput], error) {
		req := &authv1.RegisterRequest{
			Email:    input.Body.Email,
			Username: input.Body.Username,
			Password: input.Body.Password,
		}
		if input.Body.Phone != nil {
			req.Phone = *input.Body.Phone
		}
		resp, err := auth.Register(ctx, req)
		if err != nil {
			return nil, grpcError(err)
		}
		out := &RegisterOutput{}
		out.Body.UserId = resp.GetUserId()
		out.Body.Email = resp.GetEmail()
		out.Body.Username = resp.GetUsername()
		return SuccessResponse(200, out, "registered succesfully"), nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Summary:     "Log in",
		Description: "Forwards to the auth service over gRPC.",
		Method:      http.MethodPost,
		Path:        "api/auth/login",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, input *LoginInput) (*BaseResponse[LoginOutput], error) {
		resp, err := auth.Login(ctx, &authv1.LoginRequest{
			Identifier: input.Body.Identifier,
			Password:   input.Body.Password,
		})
		if err != nil {
			return nil, grpcError(err)
		}
		out := &LoginOutput{}
		out.Body.AccessToken = resp.GetAccessToken()
		out.Body.UserId = resp.GetUserId()
		return SuccessResponse(200, out, "logged in succesfully"), nil
	})
}

func grpcError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return huma.NewError(http.StatusBadGateway, err.Error())
	}

	var code int
	switch st.Code() {
	case codes.InvalidArgument:
		code = http.StatusBadRequest
	case codes.Unauthenticated:
		code = http.StatusUnauthorized
	case codes.NotFound:
		code = http.StatusNotFound
	case codes.Unimplemented:
		code = http.StatusNotImplemented
	case codes.DeadlineExceeded, codes.Unavailable:
		code = http.StatusBadGateway
	default:
		code = http.StatusInternalServerError
	}

	return huma.NewError(code, st.Message())
}

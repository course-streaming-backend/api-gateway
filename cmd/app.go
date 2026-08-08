package main

import (
	"log"
	"os"

	"github.com/course-streaming-backend/api-gateway/internal/transport"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	grpcpb "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type app struct {
	Api  *transport.Api
	addr string
}

func (a *app) mount(engine *gin.Engine) {
	config := huma.DefaultConfig(
		"Course streaming backend - API Gateway",
		"1.0.0",
	)
	humaApi := humagin.New(engine, config)

	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authAddr == "" {
		authAddr = "auth-service:50051"
	}

	conn, err := grpcpb.NewClient(authAddr,
		grpcpb.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to auth service: %v", err)
	}

	api := transport.InitAPI(humaApi, conn)
	a.Api = api
}

func (a *app) run() {
	engine := gin.Default()
	a.mount(engine)

	if err := engine.Run(a.addr); err != nil {
		log.Fatal(err)
	}
}

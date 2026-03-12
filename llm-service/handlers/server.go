package handlers

import (
	"context"
	"time"

	"github.com/afrizalsebastian/llm-integration-service/llm-service/application/controllers"
	proto "github.com/afrizalsebastian/llm-integration-service/proto/gen/go/llm/v1"
)

type Server struct {
	proto.UnimplementedLlmServiceServer
	llmGrpcController controllers.ILlmGRPCController
}

func NewServer() (*Server, error) {
	di := initDI()
	server := &Server{
		llmGrpcController: di.llmGrpcController,
	}

	return server, nil
}

func (s *Server) HelloWorld(ctx context.Context, req *proto.HelloWorldRequest) (*proto.HelloWorldResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.llmGrpcController.HelloWorld(ctx, req)
}

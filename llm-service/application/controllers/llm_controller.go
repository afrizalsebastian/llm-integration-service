package controllers

import (
	"context"

	proto "github.com/afrizalsebastian/llm-integration-service/proto/gen/go/llm/v1"
)

type ILlmGRPCController interface {
	HelloWorld(context.Context, *proto.HelloWorldRequest) (*proto.HelloWorldResponse, error)
}

type llmGRPCController struct {
}

func NewLlmGRPCController() ILlmGRPCController {
	return &llmGRPCController{}
}

func (c *llmGRPCController) HelloWorld(ctx context.Context, request *proto.HelloWorldRequest) (*proto.HelloWorldResponse, error) {
	return &proto.HelloWorldResponse{
		Message: "Success",
		Status:  200,
		Data:    "Hello World",
	}, nil
}

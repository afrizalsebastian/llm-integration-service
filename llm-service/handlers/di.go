package handlers

import "github.com/afrizalsebastian/llm-integration-service/llm-service/application/controllers"

type ServerController struct {
	llmGrpcController controllers.ILlmGRPCController
}

func initDI() *ServerController {
	return &ServerController{
		llmGrpcController: llmGrpc(),
	}
}

func llmGrpc() controllers.ILlmGRPCController {
	return controllers.NewLlmGRPCController()
}

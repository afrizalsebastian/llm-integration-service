package controllers

import (
	"context"
	"net/http"

	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/api"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/services"
)

type IHelloController interface {
	GetHello(context.Context, *http.Request) api.WebResponse
}

type helloController struct {
	helloService services.IHelloService
}

func NewHelloController(helloService services.IHelloService) IHelloController {
	return &helloController{
		helloService: helloService,
	}
}

func (h *helloController) GetHello(ctx context.Context, r *http.Request) api.WebResponse {
	return h.helloService.GetHello(ctx)
}

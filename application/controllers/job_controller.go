package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/afrizalsebastian/llm-integration-service/api"
	"github.com/afrizalsebastian/llm-integration-service/application/helper"
	"github.com/afrizalsebastian/llm-integration-service/application/services"
	"github.com/afrizalsebastian/llm-integration-service/domain/models"
)

type IJobController interface {
	EnqueueJob(ctx context.Context, r *http.Request) api.WebResponse
	ResultJob(ctx context.Context, r *http.Request, jobId string) api.WebResponse
}

type jobController struct {
	jobService services.IJobService
}

func NewEvaluateController(
	jobService services.IJobService,
) IJobController {
	return &jobController{
		jobService: jobService,
	}
}

func (e *jobController) EnqueueJob(ctx context.Context, r *http.Request) api.WebResponse {
	request, err := helper.ParseJSONBodyRequest[models.EvaluateRequest](r)
	if err != nil {
		log.Println("error when parse body request")
		return api.CreateWebResponse("invalid request", http.StatusBadRequest, nil, nil)
	}

	// validation
	if err := helper.ValidateParams(ctx, request); err != nil {
		log.Println("validation error")
		return api.CreateWebResponse("validation error", http.StatusBadRequest, nil, err)
	}

	resp := e.jobService.EnqueueJob(ctx, request)
	return resp
}

func (e *jobController) ResultJob(ctx context.Context, r *http.Request, jobId string) api.WebResponse {
	resp := e.jobService.ResultJob(ctx, jobId)
	return resp
}

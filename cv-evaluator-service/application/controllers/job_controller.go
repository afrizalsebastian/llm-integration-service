package controllers

import (
	"context"
	"net/http"

	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/api"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/helper"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/services"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/models/dto"
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
	l := logger.New().WithContext(ctx)

	request, err := helper.ParseJSONBodyRequest[dto.EvaluateRequest](r)
	if err != nil {
		l.Error("error when parse body request").Msg()
		return api.CreateWebResponse("invalid request", http.StatusBadRequest, nil, nil)
	}

	// validation
	if err := helper.ValidateParams(ctx, request); err != nil {
		l.Error("validation error").Msg()
		return api.CreateWebResponse("validation error", http.StatusBadRequest, nil, err)
	}

	resp := e.jobService.EnqueueJob(ctx, request)
	return resp
}

func (e *jobController) ResultJob(ctx context.Context, r *http.Request, jobId string) api.WebResponse {
	resp := e.jobService.ResultJob(ctx, jobId)
	return resp
}

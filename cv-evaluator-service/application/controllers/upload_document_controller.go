package controllers

import (
	"context"
	"net/http"

	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/api"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/helper"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/services"
)

type IUploadDocumentController interface {
	PostDocument(ctx context.Context, r *http.Request) api.WebResponse
}

type uploadDocumentController struct {
	uploadDocumentService services.IUploadDocumentService
}

func NewUploadDocumenController(uploadDocumentService services.IUploadDocumentService) IUploadDocumentController {
	return &uploadDocumentController{
		uploadDocumentService: uploadDocumentService,
	}
}

func (u *uploadDocumentController) PostDocument(ctx context.Context, r *http.Request) api.WebResponse {
	l := logger.New().WithContext(ctx)

	multipartRequest, err := helper.ParseMultipartRequest(r)
	if err != nil {
		l.Error("error when parse body request").Msg()
		return api.CreateWebResponse("Invalid request", http.StatusBadRequest, nil, nil)
	}

	request := helper.MultipartToUploadDocumentRequest(multipartRequest)

	if structErr := helper.ValidateParams(ctx, request); structErr != nil {
		l.Error("error validation on body request").Msg()
		return api.CreateWebResponse("Invalid request", http.StatusBadRequest, nil, structErr)
	}

	resp := u.uploadDocumentService.SaveUploadedDocument(ctx, request)
	return resp

}

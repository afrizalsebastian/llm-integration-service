package services

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/api"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/models/dao"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/models/dto"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IJobService interface {
	EnqueueJob(context.Context, *dto.EvaluateRequest) api.WebResponse
	ResultJob(context.Context, string) api.WebResponse
}

type jobService struct {
	cvEvaluatorJobRepository repository.ICvEvaluatorJobRepository
	kafkaProducer            IKafkaProducer

	//Topic
	cvEvaluatorTopic string
}

func NewEvaluateServce(cvEvaluatorJobRepository repository.ICvEvaluatorJobRepository, kafkaProducer IKafkaProducer, cvEvaluatorTopic string) IJobService {
	return &jobService{
		cvEvaluatorJobRepository: cvEvaluatorJobRepository,
		kafkaProducer:            kafkaProducer,
		cvEvaluatorTopic:         cvEvaluatorTopic,
	}
}

func (e *jobService) EnqueueJob(ctx context.Context, request *dto.EvaluateRequest) api.WebResponse {
	jobId := uuid.New().String()
	jobItem := &dao.CvEvaluatorJob{
		JobId:    jobId,
		JobTitle: request.JobTitle,
		FileId:   request.FileId,
		Status:   dto.StatusQueued,
	}

	if err := e.cvEvaluatorJobRepository.CreateJobItem(ctx, jobItem); err != nil {
		log.Println("failed to create job")
		return api.CreateWebResponse("internal server error", http.StatusInternalServerError, nil, nil)
	}

	go e.kafkaProducer.PublishMessage(ctx, e.cvEvaluatorTopic, nil, jobId)

	resp := &dto.EvaluateResponse{
		JobId:  jobId,
		Status: string(jobItem.Status),
	}

	return api.CreateWebResponse("Success to enqueue the job", http.StatusOK, resp, nil)
}

func (e *jobService) ResultJob(ctx context.Context, jobId string) api.WebResponse {
	jobItem, err := e.cvEvaluatorJobRepository.GetByJobId(ctx, jobId)
	if err != nil {
		log.Println("error when get job")

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.CreateWebResponse("Job Not Found", http.StatusNotFound, nil, nil)
		}

		return api.CreateWebResponse("internal server error", http.StatusInternalServerError, nil, nil)
	}

	resp := &dto.JobItem{
		Id:       jobItem.JobId,
		JobTitle: jobItem.JobTitle,
		FileId:   jobItem.FileId,
		Status:   jobItem.Status,
		Result: dto.JobResult{
			CvMatchRate:     jobItem.CvMatchRate,
			CvFeedback:      jobItem.CvFeedback,
			ProjectScore:    jobItem.ProjectScore,
			ProjectFeedback: jobItem.ProjectFeedback,
			OverallSummary:  jobItem.OverallSummary,
		},
	}

	return api.CreateWebResponse("Success", http.StatusOK, resp, nil)
}

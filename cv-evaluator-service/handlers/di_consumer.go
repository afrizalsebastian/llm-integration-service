package handlers

import (
	controller_consumer "github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/controllers/consumer"
	service_consumer "github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/services/consumer"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/bootstrap"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/repository"
)

type ConsumerController struct {
	CvEvaluatorConsumer controller_consumer.ICvEvaluatorControllerConsumer
}

func initDIConsumer(app *bootstrap.Application) *ConsumerController {
	initDi := &ConsumerController{
		CvEvaluatorConsumer: cvEvaluatorConsumer(app),
	}

	return initDi
}

func cvEvaluatorConsumer(app *bootstrap.Application) controller_consumer.ICvEvaluatorControllerConsumer {
	cvEvaluatorJobItem := repository.NewCvEvaluatorJobRepository(app)
	cvEvaluatorServiceConsumer := service_consumer.NewCvEvaluatorConsumerService(app.GeminiClient, app.ChromaClient, app.Ingest, cvEvaluatorJobItem)
	cvEvaluatorControllerConsumer := controller_consumer.NewCvEvaluatorConsumer(cvEvaluatorServiceConsumer)
	return cvEvaluatorControllerConsumer
}

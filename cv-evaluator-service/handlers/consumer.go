package handlers

import (
	controller_consumer "github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/controllers/consumer"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/bootstrap"
)

type Consumer struct {
	CvEvaluatorConsumer controller_consumer.ICvEvaluatorControllerConsumer
}

func NewConsumer(app *bootstrap.Application) (*Consumer, error) {
	di := initDIConsumer(app)
	consumer := &Consumer{
		CvEvaluatorConsumer: di.CvEvaluatorConsumer,
	}

	return consumer, nil
}

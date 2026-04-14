package controller_consumer

import (
	"context"
	"encoding/json"

	"github.com/afrizalsebastian/go-common-modules/kafka"
	service_consumer "github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/services/consumer"
)

type ICvEvaluatorControllerConsumer interface {
	kafka.ConsumerController
}

type cvEvaluatorControllerConsumer struct {
	cvEvaluatorServiceConsumer service_consumer.ICvEvaluatorConsumerService
}

func NewCvEvaluatorConsumer(
	cvEvaluatorServiceConsumer service_consumer.ICvEvaluatorConsumerService,
) ICvEvaluatorControllerConsumer {
	return &cvEvaluatorControllerConsumer{
		cvEvaluatorServiceConsumer: cvEvaluatorServiceConsumer,
	}
}

func (c *cvEvaluatorControllerConsumer) ProcessMessage(ctx context.Context, msg *kafka.Message) error {
	request := msg.Value
	var jobId string
	_ = json.Unmarshal(request, &jobId)
	return c.cvEvaluatorServiceConsumer.RunningJob(ctx, jobId)
}

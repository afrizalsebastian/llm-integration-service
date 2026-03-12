package handlers

import (
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/controllers"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/application/services"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/bootstrap"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/repository"
)

type ServeController struct {
	Hello          controllers.IHelloController
	UploadDocument controllers.IUploadDocumentController
	Evaluate       controllers.IJobController
}

func initDI(app *bootstrap.Application) *ServeController {
	init := &ServeController{
		Hello:          hello(app),
		UploadDocument: uploadDocument(app),
		Evaluate:       evaluate(app),
	}

	return init
}

func hello(_ *bootstrap.Application) controllers.IHelloController {
	helloService := services.NewHelloService()
	helloController := controllers.NewHelloController(helloService)
	return helloController
}

func uploadDocument(_ *bootstrap.Application) controllers.IUploadDocumentController {
	uploadDocumentService := services.NewUploadDocumentService("./uploaded-file")
	uploadDocumentController := controllers.NewUploadDocumenController(uploadDocumentService)
	return uploadDocumentController
}

func evaluate(app *bootstrap.Application) controllers.IJobController {
	cvEvaluatorJobRepository := repository.NewCvEvaluatorJobRepository(app)
	kafkaProducer := services.NewKafkaProducer(app.KafkaProducer)
	evaluateService := services.NewEvaluateServce(cvEvaluatorJobRepository, kafkaProducer, app.ENV.KafkaCvEvaluatorTopic)
	evaluateController := controllers.NewEvaluateController(evaluateService)
	return evaluateController
}

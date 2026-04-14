package bootstrap

import (
	"context"
	"os"

	"github.com/IBM/sarama"
	appconfig "github.com/afrizalsebastian/go-common-modules/app-config"
	gomysql "github.com/afrizalsebastian/go-common-modules/go-mysql"
	"github.com/afrizalsebastian/go-common-modules/kafka"
	"github.com/afrizalsebastian/go-common-modules/logger"
	chromaclient "github.com/afrizalsebastian/llm-integration-service/modules/chroma-client"
	"github.com/afrizalsebastian/llm-integration-service/modules/config"
	geminiclient "github.com/afrizalsebastian/llm-integration-service/modules/gemini-client"
	ingestdocument "github.com/afrizalsebastian/llm-integration-service/modules/ingest-document"
	"gorm.io/gorm"
)

type Application struct {
	ENV           *config.Config
	GeminiClient  geminiclient.IGeminiClient
	ChromaClient  chromaclient.IChromaClient
	Ingest        ingestdocument.IIngestFile
	DB            *gorm.DB
	KafkaProducer *kafka.Producer
}

func NewApp() *Application {
	ctx := context.Background()
	app := &Application{}
	l := logger.New()

	wd, err := os.Getwd()
	if err != nil {
		l.Error("Failed to get working directory").Msg()
		os.Exit(1)
	}

	app.ENV, err = appconfig.Init[config.Config](wd)
	if err != nil {
		l.Error("failed to initialize configuration").Msg()
		os.Exit(1)
	}

	// Init DB
	dbConfig := &gomysql.MysqlConfig{
		DBUser:     app.ENV.DBUser,
		DBPassword: app.ENV.DBPassword,
		DBHost:     app.ENV.DBHost,
		DBPort:     app.ENV.DBPort,
		DBName:     app.ENV.DBName,
	}
	db, err := gomysql.NewDatabaseConnection(dbConfig)
	if err != nil {
		l.Error("failed to create db connection").Msg()
		os.Exit(1)
	}
	app.DB = db

	// Init Gemini Client
	geminiCient, err := geminiclient.NewGeminiAiCLient(ctx, app.ENV.GeminiApiKey, app.ENV.GeminiModel)
	if err != nil {
		l.Error("failed to init gemini client").Msg()
		os.Exit(1)
	}
	app.GeminiClient = geminiCient

	// Init chroma
	chromaClient, err := chromaclient.NewChromaClient(ctx, app.ENV.ChromaUrl)
	if err != nil {
		l.Errorf("failed to init chroma client, %s", err.Error()).Msg()
		os.Exit(1)
	}
	app.ChromaClient = chromaClient

	// Init ingestDocument
	ingesDocument := ingestdocument.NewIngestFile(chromaClient)
	app.Ingest = ingesDocument

	// Init Kafka Producer client
	kafkaProducer, err := kafka.NewProducer(
		app.ENV.KafkaBroker,
		func(c *sarama.Config) {
			c.Net.SASL.Enable = app.ENV.KafkaSASLEnable
			c.Net.SASL.Handshake = app.ENV.KafkaSASLHandshake
			c.Net.TLS.Enable = app.ENV.KafkaTLS
		},
	)
	if err != nil {
		l.Warn("Kafka producer failed to initialize")
	}
	app.KafkaProducer = kafkaProducer

	return app
}

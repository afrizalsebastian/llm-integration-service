package bootstrap

import (
	"context"
	"log"
	"os"

	"github.com/IBM/sarama"
	appconfig "github.com/afrizalsebastian/llm-integration-service/modules/app-config"
	chromaclient "github.com/afrizalsebastian/llm-integration-service/modules/chroma-client"
	geminiclient "github.com/afrizalsebastian/llm-integration-service/modules/gemini-client"
	gomysql "github.com/afrizalsebastian/llm-integration-service/modules/go-mysql"
	ingestdocument "github.com/afrizalsebastian/llm-integration-service/modules/ingest-document"
	"github.com/afrizalsebastian/llm-integration-service/modules/kafka"
	"gorm.io/gorm"
)

type Application struct {
	ENV           *appconfig.Config
	GeminiClient  geminiclient.IGeminiClient
	ChromaClient  chromaclient.IChromaClient
	Ingest        ingestdocument.IIngestFile
	DB            *gorm.DB
	KafkaProducer *kafka.Producer
}

func NewApp() *Application {
	ctx := context.Background()
	app := &Application{}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory")
	}

	app.ENV, err = appconfig.Init(wd)
	if err != nil {
		log.Fatal("failed to initialize configuration")
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
		log.Fatal("failed to create db connection")
	}
	app.DB = db

	// Init Gemini Client
	geminiCient, err := geminiclient.NewGeminiAiCLient(ctx, app.ENV.GeminiApiKey, app.ENV.GeminiModel)
	if err != nil {
		log.Fatal("failed to init gemini client")
	}
	app.GeminiClient = geminiCient

	// Init chroma
	chromaClient, err := chromaclient.NewChromaClient(ctx, app.ENV.ChromaUrl)
	if err != nil {
		log.Fatalf("failed to init chroma client, %s", err.Error())
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
		log.Println("Kafka producer failed to initialize")
	}
	app.KafkaProducer = kafkaProducer

	return app
}

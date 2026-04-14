package config

import (
	"errors"
)

var (
	ErrOsGetPwd            = errors.New("failed to get working directory")
	ErrReadConfigFile      = errors.New("failed to read config file")
	ErrUnmarshalConfigFile = errors.New("failed to read config file")
)

type Config struct {
	AppPort                    int      `mapstructure:"PORT"`
	GeminiApiKey               string   `mapstructure:"GEMINI_API_KEY"`
	ChromaUrl                  string   `mapstructure:"CHROMA_URL"`
	GeminiModel                string   `mapstructure:"GEMINI_MODEl"`
	DBUser                     string   `mapstructure:"DB_USER"`
	DBPassword                 string   `mapstructure:"DB_PASSWORD"`
	DBHost                     string   `mapstructure:"DB_HOST"`
	DBPort                     string   `mapstructure:"DB_PORT"`
	DBName                     string   `mapstructure:"DB_NAME"`
	KafkaTLS                   bool     `mapstructure:"KAFKA_TLS"`
	KafkaSASLEnable            bool     `mapstructure:"KAFKA_SASL_ENABLE"`
	KafkaSASLHandshake         bool     `mapstructure:"KAFKA_SASL_HANDSHAKE"`
	KafkaBroker                []string `mapstructure:"KAFKA_BROKER"`
	KafkaMaxRetryPolicy        int      `mapstructure:"KAFKA_MAX_RETRY_POLICY"`
	KafkaCvEvaluatorTopic      string   `mapstructure:"KAFKA_CV_EVALUATOR_TOPIC"`
	KafkaCvEvaluatorTopicGroup string   `mapstructure:"KAFKA_CV_EVALUATOR_TOPIC_GROUP"`
	GrpcPort                   string   `mapstructure:"GRPC_PORT"`

	GoogleAuthClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleAuthClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleCallbackUrl      string `mapstructure:"GOOGLE_CALLBACK_URL"`
}

package cli

import (
	"context"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/afrizalsebastian/go-common-modules/kafka"
	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/bootstrap"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/handlers"
	"github.com/spf13/cobra"
)

func init() {
	consumerCommand.Flags().String("topic", "", "kafka consumer topic")
	rootCmd.AddCommand(consumerCommand)
}

var consumerCommand = &cobra.Command{
	Use:   "consumer",
	Short: "Start consumer for Go CV Evaluator",
	PreRun: func(cmd *cobra.Command, args []string) {
		app := bootstrap.NewApp()
		ctx := context.WithValue(cmd.Context(), appKey, app)
		cmd.SetContext(ctx)
	},
	Run: func(cmd *cobra.Command, args []string) {
		app := cmd.Context().Value(appKey).(*bootstrap.Application)
		startConsumer(app, cmd)
	},
}

func startConsumer(app *bootstrap.Application, cmd *cobra.Command) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New()
	consumerTopic, _ := cmd.Flags().GetString("topic")
	if consumerTopic == "" {
		l.Warn("Consumer topic is required. User --topic=<consumer-topic>").Msg()
		os.Exit(1)
	}

	handlerConsumer, err := handlers.NewConsumer(app)
	if err != nil {
		l.Error("Consumer handler failed to initialize.").Msg()
		return
	}

	switch consumerTopic {
	case app.ENV.KafkaCvEvaluatorTopic:
		consumer, err := createConsumer(app, handlerConsumer.CvEvaluatorConsumer, app.ENV.KafkaCvEvaluatorTopic, app.ENV.KafkaCvEvaluatorTopicGroup)
		if err != nil {
			l.Errorf("Error connection consumer, err: %v", err).Msg()
			os.Exit(1)
		}

		ctx = context.WithValue(ctx, "consumer_topic", consumerTopic)
		err = consumer.Run(ctx)
		if err != nil {
			l.Errorf("Error running consumer, err: %v", err).Msg()
			os.Exit(1)
		}

	default:
		l.Errorf("Unknown consumer name: %s", consumerTopic).Msg()
		os.Exit(1)
	}

}

func createConsumer(app *bootstrap.Application, handler kafka.ConsumerController, topic, topicGroup string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(
		kafka.WithBrokers(app.ENV.KafkaBroker...),
		kafka.WithGroupID(topicGroup),
		kafka.WithTopics(topic),
		kafka.WithConsumerController(handler),
		kafka.WithRetryPolicy(app.ENV.KafkaMaxRetryPolicy, 5*time.Second),
		kafka.WithSaramaConfig(func() *sarama.Config {
			cfg := sarama.NewConfig()
			cfg.Net.SASL.Enable = app.ENV.KafkaSASLEnable
			cfg.Net.SASL.Handshake = app.ENV.KafkaSASLHandshake
			cfg.Net.TLS.Enable = app.ENV.KafkaTLS
			cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
			cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
			return cfg
		}()),
		kafka.WithAckMode(kafka.AckModeAuto),
	)

	return consumer, err
}

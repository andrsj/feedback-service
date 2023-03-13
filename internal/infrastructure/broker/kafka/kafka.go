package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"

	"github.com/andrsj/feedback-service/internal/domain/models"
	"github.com/andrsj/feedback-service/pkg/logger"
)

const frequency = 500

type Broker interface {
	SendMessage(models.Feedback) error
	Close() error
}

type apacheKafkaBroker struct {
	logger    logger.Logger
	producer  sarama.SyncProducer
	topicName string
}

var _ Broker = (*apacheKafkaBroker)(nil)

func New(log logger.Logger, addrs []string, topicName string) *apacheKafkaBroker {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = frequency * time.Millisecond

	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		log.Fatal("Failed to create Kafka producer", logger.M{"err": err})
	}

	return &apacheKafkaBroker{
		logger:    log,
		producer:  producer,
		topicName: topicName,
	}
}

func (a *apacheKafkaBroker) SendMessage(feedback models.Feedback) error {
	feedbackJSON, err := json.Marshal(feedback)
	if err != nil {
		a.logger.Error("Failed to marshal Feedback to JSON", logger.M{"err": err})

		return fmt.Errorf("failed to marshal Feedback to JSON: %w", err)
	}

	//nolint
	message := &sarama.ProducerMessage{
		Topic: a.topicName,
		Value: sarama.StringEncoder(feedbackJSON),
	}

	partition, offset, err := a.producer.SendMessage(message)
	if err != nil {
		a.logger.Error("Failed to send Kafka message", logger.M{"err": err})

		return fmt.Errorf("failed to send Kafka message :%w", err)
	}

	a.logger.Info("Sent Kafka message", logger.M{
		"partition": partition,
		"offset":    offset,
	})

	return nil
}

func (a *apacheKafkaBroker) Close() error {
	if err := a.producer.Close(); err != nil {
		return fmt.Errorf("closing error: %w", err)
	}

	return nil
}

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

type Producer interface {
	SendMessage(*models.Feedback) error
	Close() error
}

type apacheKafkaProducer struct {
	logger    logger.Logger
	producer  sarama.SyncProducer
	topicName string
}

var _ Producer = (*apacheKafkaProducer)(nil)

func New(log logger.Logger, addr string, topicName string) (*apacheKafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = frequency * time.Millisecond

	producer, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		log.Error("Failed to create Kafka producer", logger.M{"err": err})

		return nil, fmt.Errorf("can't setting up Kafka Producer: %w", err)
	}

	return &apacheKafkaProducer{
		logger:    log,
		producer:  producer,
		topicName: topicName,
	}, nil
}

func (a *apacheKafkaProducer) SendMessage(feedback *models.Feedback) error {
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

		return fmt.Errorf("failed to send Kafka message: %w", err)
	}

	a.logger.Info("Sent Kafka message", logger.M{
		"partition": partition,
		"offset":    offset,
	})

	return nil
}

func (a *apacheKafkaProducer) Close() error {
	if err := a.producer.Close(); err != nil {
		return fmt.Errorf("closing error: %w", err)
	}

	return nil
}

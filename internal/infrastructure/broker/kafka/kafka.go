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

type Producer struct {
	logger    logger.Logger
	producer  sarama.SyncProducer
	topicName string
}

func New(log logger.Logger, addr string, topicName string) (*Producer, error) {
	log = log.Named("kafka")

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = frequency * time.Millisecond

	// Check list of topics which exist.
	brokers := []string{addr}
	
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	topics, err := client.Topics()
	if err != nil {
		return nil, fmt.Errorf("failed to get topics: %w", err)
	}

	for _, topic := range topics {
		log.Info("Topic", logger.M{"topic": topic})
	}

	producer, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		log.Error("Failed to create Kafka producer", logger.M{
			"err": err,
			"addr": addr,
			"topic": topicName,
		})

		return nil, fmt.Errorf("can't setting up Kafka Producer: %w", err)
	}

	return &Producer{
		logger:    log,
		producer:  producer,
		topicName: topicName,
	}, nil
}

func (a *Producer) SendMessage(feedback *models.Feedback) error {
	feedbackJSON, err := json.Marshal(feedback)
	if err != nil {
		a.logger.Error("Failed to marshal Feedback to JSON", logger.M{"err": err})

		return fmt.Errorf("failed to marshal Feedback to JSON: %w", err)
	}

	//nolint:exhaustivestruct,exhaustruct
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

func (a *Producer) Close() error {
	if err := a.producer.Close(); err != nil {
		return fmt.Errorf("closing error: %w", err)
	}

	return nil
}

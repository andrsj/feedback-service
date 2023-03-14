package feedback

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/andrsj/feedback-service/internal/domain/models"
	"github.com/andrsj/feedback-service/internal/infrastructure/broker/kafka"
	"github.com/andrsj/feedback-service/internal/infrastructure/db/gorm"
	"github.com/andrsj/feedback-service/internal/infrastructure/db/memory"
	"github.com/andrsj/feedback-service/pkg/logger"
)

// TODO VALIDATE INPUT

/*
This is a good practice to separate the interface:
- FeedbackRepoReader
- FeedbackRepoWriter
- FeedbackRepoSearch
If we want to use only some part of logic.
*/
type Repository interface {
	Create(feedback *models.FeedbackInput) (feedbackID uuid.UUID, err error)
	GetByID(feedbackID uuid.UUID) (feedback *models.Feedback, err error)
	GetAll() (feedbacks []*models.Feedback, err error)
}

// Check that actual implementation fits the interface.
var _ Repository = (*gorm.FeedbackRepository)(nil)
var _ Repository = (*memory.FeedbackRepository)(nil)

type Producer interface {
	SendMessage(*models.Feedback) error
	Close() error
}

// Check that actual implementation fits the interface.
var _ Producer = (*kafka.Producer)(nil)

type Service struct {
	logger   logger.Logger
	repo     Repository
	producer Producer
}

func New(feedbackRepository Repository, producer Producer, logger logger.Logger) *Service {
	return &Service{
		logger:   logger.Named("service"),
		repo:     feedbackRepository,
		producer: producer,
	}
}

func (s *Service) Create(feedback *models.FeedbackInput) (string, error) {
	var (
		feedbackID uuid.UUID
		err        error
	)

	s.logger.Info("Creating feedback", logger.M{"feedback": feedback})

	feedbackID, err = s.repo.Create(feedback)
	if err != nil {
		s.logger.Error("Creating feedback error", logger.M{"err": err})
		
		return "", fmt.Errorf("creating feedback error: %w", err)
	}
	
	//nolint:exhaustivestruct,exhaustruct
	err = s.producer.SendMessage(&models.Feedback{
		ID:           feedbackID,
		CustomerName: feedback.CustomerName,
		Email:        feedback.Email,
		FeedbackText: feedback.FeedbackText,
		Source:       feedback.Source,
		// I don't specify the createdAt, updatedAt from DB instance
		// because it's a part of DB logic, not a Kafka.
	})
	if err != nil {
		s.logger.Error("broker sending feedback error", logger.M{"err": err})

		return "", fmt.Errorf("broker sending feedback error: %w", err)
	}

	s.logger.Info("Successfully created feedback", logger.M{"feedbackID": feedbackID.String()})

	return feedbackID.String(), nil
}

func (s *Service) GetByID(feedbackID string) (*models.Feedback, error) {
	var (
		feedback *models.Feedback
		err      error
	)

	s.logger.Info("Getting one feedback", logger.M{"feedbackID": feedbackID})

	feedbackUUID, err := uuid.Parse(feedbackID)
	if err != nil {
		s.logger.Error("Parsing UUID", logger.M{
			"feedbackID": feedbackID,
			"error":      err,
		})

		return nil, fmt.Errorf("can't parse the ID: %w", err)
	}

	feedback, err = s.repo.GetByID(feedbackUUID)
	if err != nil {
		s.logger.Error("Getting by ID", logger.M{
			"feedbackID": feedbackID,
			"error":      err,
		})

		return nil, fmt.Errorf("getting by ID: %w", err)
	}

	s.logger.Info("Returning successful result", logger.M{"feedbackID": feedback.ID})

	return feedback, nil
}

func (s *Service) GetAll() ([]*models.Feedback, error) {
	var (
		feedbacks []*models.Feedback
		err       error
	)

	s.logger.Info("Getting All feedbacks", nil)

	feedbacks, err = s.repo.GetAll()
	if err != nil {
		s.logger.Error("Getting by ID", logger.M{"error": err})

		return nil, fmt.Errorf("error by getting from repository: %w", err)
	}

	s.logger.Info("Returning successful result", logger.M{"result": len(feedbacks)})

	return feedbacks, nil
}

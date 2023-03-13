package feedback

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/andrsj/feedback-service/internal/domain/models"
	"github.com/andrsj/feedback-service/internal/domain/repositories"
	"github.com/andrsj/feedback-service/internal/infrastructure/broker/kafka"
	"github.com/andrsj/feedback-service/pkg/logger"
)

// TODO VALIDATE INPUT

type Service interface {
	Create(feedback *models.FeedbackInput) (string, error)
	GetByID(feedbackID string) (*models.Feedback, error)
	GetAll() ([]*models.Feedback, error)
}

type feedbackService struct {
	logger   logger.Logger
	repo     repositories.FeedbackRepository
	producer kafka.Producer
}

var _ Service = (*feedbackService)(nil)

func New(feedbackRepository repositories.FeedbackRepository, producer kafka.Producer, logger logger.Logger) *feedbackService {
	return &feedbackService{
		logger:   logger.Named("service"),
		repo:     feedbackRepository,
		producer: producer,
	}
}

func (s *feedbackService) Create(feedback *models.FeedbackInput) (string, error) {
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
	
	err = s.producer.SendMessage(&models.Feedback{
		ID:           feedbackID,
		CustomerName: feedback.CustomerName,
		Email:        feedback.Email,
		FeedbackText: feedback.FeedbackText,
		Source:       feedback.Source,
	})
	if err != nil {
		s.logger.Error("broker sending feedback error", logger.M{"err": err})

		return "", fmt.Errorf("broker sending feedback error: %w", err)
	}

	s.logger.Info("Successfully created feedback", logger.M{"feedbackID": feedbackID.String()})

	return feedbackID.String(), nil
}

func (s *feedbackService) GetByID(feedbackID string) (*models.Feedback, error) {
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

func (s *feedbackService) GetAll() ([]*models.Feedback, error) {
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

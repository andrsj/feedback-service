package feedback

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/andrsj/feedback-service/internal/domain/models"
	"github.com/andrsj/feedback-service/internal/domain/repositories"
	"github.com/andrsj/feedback-service/pkg/logger"
)

type Service interface {
	Create(feedback *models.FeedbackInput) (string, error)
	GetByID(feedbackID string) (*models.Feedback, error)
	GetAll() ([]*models.Feedback, error)
}

type feedbackService struct {
	logger logger.Logger
	repo   repositories.FeedbackRepository
}

var _ Service = (*feedbackService)(nil)

func New(feedbackRepository repositories.FeedbackRepository, logger logger.Logger) *feedbackService {
	return &feedbackService{
		logger: logger.Named("service"),
		repo:   feedbackRepository,
	}
}

func (s *feedbackService) Create(feedback *models.FeedbackInput) (string, error) {
	var (
		feedbackID string
		err        error
	)

	s.logger.Info("Creating feedback", logger.M{"feedback": feedback.CustomerName})

	feedbackID, err = s.repo.Create(feedback)
	if err != nil {
		s.logger.Error("Creating feedback error", logger.M{"err": err})

		return "", fmt.Errorf("creating feedback error: %w", err)
	}

	s.logger.Info("Successfully created feedback", logger.M{"feedbackID": feedbackID})

	return feedbackID, nil
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

	// TODO validation and business logic

	feedbacks, err = s.repo.GetAll()
	if err != nil {
		s.logger.Error("Getting by ID", logger.M{"error": err})

		return nil, fmt.Errorf("error by getting from repository: %w", err)
	}

	s.logger.Info("Returning successful result", logger.M{"result": len(feedbacks)})

	return feedbacks, nil
}

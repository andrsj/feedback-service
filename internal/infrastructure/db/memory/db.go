package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/andrsj/feedback-service/internal/domain/models"
	"github.com/andrsj/feedback-service/pkg/logger"
)

type FeedbackRepository struct {
	mu        sync.Mutex
	feedbacks map[string]*models.Feedback
	logger    logger.Logger
}

func New(logger logger.Logger) *FeedbackRepository {
	return &FeedbackRepository{
		mu:        sync.Mutex{},
		feedbacks: make(map[string]*models.Feedback),
		logger:    logger.Named("memoryDB"),
	}
}

func (r *FeedbackRepository) Create(feedback *models.FeedbackInput) (uuid.UUID, error) {
	feedbackID := uuid.New()

	r.logger.Info("Creating feedback", logger.M{"feedbackID": feedbackID})

	feedbackOutput := &models.Feedback{
		ID:           feedbackID,
		CustomerName: feedback.CustomerName,
		Email:        feedback.Email,
		FeedbackText: feedback.FeedbackText,
		Source:       feedback.Source,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	r.mu.Lock()
	r.logger.Info("Saving feedback", logger.M{"feedbackID": feedbackID})
	r.feedbacks[feedbackID.String()] = feedbackOutput
	r.mu.Unlock()
	
	r.logger.Info("Returning feedbackID for successfully saved feedback", logger.M{"feedbackID": feedbackID})
	
	return feedbackID, nil
}

func (r *FeedbackRepository) GetByID(feedbackID uuid.UUID) (*models.Feedback, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("Getting feedback from map", logger.M{"feedbackID": feedbackID})

	feedbackOutput, ok := r.feedbacks[feedbackID.String()]
	if !ok {
		r.logger.Error("Feedback not found for ID", logger.M{"feedbackID": feedbackID})

		return nil, fmt.Errorf("feedback not found for ID '%s'", feedbackID) //nolint:goerr113
	}

	r.logger.Info("Getting feedback from map successfully", logger.M{"feedbackID": feedbackID})

	return feedbackOutput, nil
}

func (r *FeedbackRepository) GetAll() ([]*models.Feedback, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var feedbacks = make([]*models.Feedback, 0, len(r.feedbacks))
	for _, feedback := range r.feedbacks {
		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}

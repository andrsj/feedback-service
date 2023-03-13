package repositories

import (
	"github.com/google/uuid"

	"github.com/andrsj/feedback-service/internal/domain/models"
)

type FeedbackRepository interface {
	Create(feedback *models.FeedbackInput) (feedbackID string, err error)
	GetByID(feedbackID uuid.UUID) (feedback *models.Feedback, err error)
	GetAll() (feedbacks []*models.Feedback, err error)
}

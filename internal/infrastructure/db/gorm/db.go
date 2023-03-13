package gorm

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"

	"github.com/andrsj/feedback-service/internal/domain/models"
	"github.com/andrsj/feedback-service/internal/domain/repositories"
	log "github.com/andrsj/feedback-service/pkg/logger"

)

type feedbackRepository struct {
	db     *gorm.DB
	logger log.Logger
}

var _ repositories.FeedbackRepository = (*feedbackRepository)(nil)

//nolint:varnamelen
func NewFeedbackRepository(db *gorm.DB, logger log.Logger) (*feedbackRepository, error) {
	logger = logger.Named("gormORM")
	
	//nolint
	err := db.AutoMigrate(models.Feedback{})
	if err != nil {
		logger.Error("Can't Auto Migrate the 'Feedback' model", log.M{"err": err})

		return nil, fmt.Errorf("can't Auto Migrate the 'Feedback' model: %w", err)
	}

	logger.Info("Successfully migrated", nil)

	return &feedbackRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *feedbackRepository) Create(feedbackInput *models.FeedbackInput) (string, error) {
	r.logger.Info("Creating 'Feedback'", nil)

	var (
		feedback   *models.Feedback
		feedbackID = uuid.New()
	)

	feedback = &models.Feedback{
		ID:           feedbackID,
		CustomerName: feedbackInput.CustomerName,
		Email:        feedbackInput.Email,
		FeedbackText: feedbackInput.FeedbackText,
		Source:       feedbackInput.Source,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := r.db.Create(feedback).Error
	if err != nil {
		r.logger.Error("Failed to create feedback into DB", log.M{"err": err})

		return "", fmt.Errorf("failed to create feedback into DB: %w", err)
	}

	r.logger.Info("Feedback created successfully", log.M{"id": feedbackID})

	return feedbackID.String(), nil
}

func (r *feedbackRepository) GetByID(feedbackID uuid.UUID) (*models.Feedback, error) {
	var feedback models.Feedback

	r.logger.Info("Getting 'Feedback' by ID", log.M{
		"feedbackID": feedbackID,
	})

	err := r.db.First(&feedback, feedbackID).Error
	if err != nil {
		r.logger.Error("Failed to get feedback from DB", log.M{
			"feedbackID": feedbackID,
			"error":      err.Error(),
		})

		return nil, fmt.Errorf("failed to get feedback from DB: %w", err)
	}

	r.logger.Info("Got 'Feedback' by ID", log.M{"feedbackID": feedbackID})

	return &feedback, nil
}

func (r *feedbackRepository) GetAll() ([]*models.Feedback, error) {
	var feedbacks []*models.Feedback

	r.logger.Info("Get all 'Feedback's", nil)

	err := r.db.Order("created_at").Find(&feedbacks).Error
	if err != nil {
		r.logger.Error("Failed to get feedback from DB", log.M{"error": err.Error()})

		return nil, fmt.Errorf("failed to get feedbacks from DB: %w", err)
	}

	r.logger.Info("Got all 'Feedback's", log.M{"count": len(feedbacks)})

	return feedbacks, nil
}

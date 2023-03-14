package gorm

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/andrsj/feedback-service/internal/domain/models"
	log "github.com/andrsj/feedback-service/pkg/logger"
)

type FeedbackRepository struct {
	db     *gorm.DB
	logger log.Logger
}

//nolint:varnamelen
func NewFeedbackRepository(db *gorm.DB, logger log.Logger) (*FeedbackRepository, error) {
	logger = logger.Named("gormORM")

	//nolint:exhaustivestruct,exhaustruct
	err := db.AutoMigrate(models.Feedback{})
	if err != nil {
		logger.Error("Can't Auto Migrate the 'Feedback' model", log.M{"err": err})

		return nil, fmt.Errorf("can't Auto Migrate the 'Feedback' model: %w", err)
	}

	logger.Info("Successfully migrated", nil)

	return &FeedbackRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *FeedbackRepository) Create(feedbackInput *models.FeedbackInput) (uuid.UUID, error) {
	r.logger.Info("Creating 'Feedback'", nil)

	var (
		feedback *models.Feedback
		// Potential mistake: WE CAN'T BE SURE THAT DB does not have the same ID.
		// But for test task I ignore it for simplify
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

		return uuid.Nil, fmt.Errorf("failed to create feedback into DB: %w", err)
	}

	r.logger.Info("Feedback created successfully", log.M{"id": feedbackID})

	return feedbackID, nil
}

func (r *FeedbackRepository) GetByID(feedbackID uuid.UUID) (*models.Feedback, error) {
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

func (r *FeedbackRepository) GetPage(limit int, next uuid.UUID) ([]*models.Feedback, uuid.UUID, error) {
	var (
		feedbacks []*models.Feedback
		cursor    uuid.UUID
	)

	r.logger.Info("Get page of 'Feedback's", log.M{"limit": limit, "next": next})

	if next != uuid.Nil {
		query := r.db.
			Where("created_at > (SELECT created_at FROM feedbacks WHERE id = ?)", next).
			Order("created_at").
			Limit(limit).
			Find(&feedbacks)

		if err := query.Error; err != nil {
			r.logger.Error("Failed to get feedback page from DB", log.M{"error": err.Error()})

			return nil, uuid.Nil, fmt.Errorf("failed to get feedback page from DB: %w", err)
		}
	} else {
		query := r.db.
			Order("created_at").
			Limit(limit).
			Find(&feedbacks)

		if err := query.Error; err != nil {
			r.logger.Error("Failed to get feedback page from DB", log.M{"error": err.Error()})

			return nil, uuid.Nil, fmt.Errorf("failed to get feedback page from DB: %w", err)
		}
	}

	if len(feedbacks) > 0 {
		cursor = feedbacks[len(feedbacks)-1].ID
	}

	r.logger.Info("Got page of 'Feedback's", log.M{"count": len(feedbacks), "cursor": cursor})

	return feedbacks, cursor, nil
}

func (r *FeedbackRepository) GetAll() ([]*models.Feedback, error) {
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

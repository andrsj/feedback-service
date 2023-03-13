package models

import (
	"time"

	"github.com/google/uuid"
)

type FeedbackInput struct {
	CustomerName string `json:"customer_name"` //nolint:tagliatelle
	Email        string `json:"email"`
	FeedbackText string `json:"feedback_text"` //nolint:tagliatelle
	Source       string `json:"source"`
}

type Feedback struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	CustomerName string    `json:"customer_name"` //nolint:tagliatelle
	Email        string    `json:"email"`
	FeedbackText string    `json:"feedback_text"` //nolint:tagliatelle
	Source       string    `json:"source"`
	CreatedAt    time.Time `json:"-" gorm:"created_at"`
	UpdatedAt    time.Time `json:"-" gorm:"updated_at"`
}

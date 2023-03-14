package feedback

import (
	"errors"
	"net/mail"
	"net/url"

	"github.com/andrsj/feedback-service/internal/domain/models"
)

var (
	errValidEmail = errors.New("invalid email address")
	errValidURL = errors.New("invalid source URL")
)

func Validate(feedback *models.FeedbackInput) error {
	_, err := mail.ParseAddress(feedback.Email)
	if err != nil {
		return errValidEmail
	}

	_, err = url.ParseRequestURI(feedback.Source)
	if err != nil {
		return errValidURL
	}

	return nil
}
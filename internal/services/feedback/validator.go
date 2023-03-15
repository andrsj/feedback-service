package feedback

import (
	"errors"
	"net/mail"
	"regexp"

	"github.com/andrsj/feedback-service/internal/domain/models"
)

var (
	errValidEmail = errors.New("invalid email address")
	errValidURL = errors.New("invalid source URL")

	regexURL = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
)

func Validate(feedback *models.FeedbackInput) error {
	_, err := mail.ParseAddress(feedback.Email)
	if err != nil {
		return errValidEmail
	}

	if !regexURL.MatchString(feedback.Source) {
		return errValidURL
	}

	return nil
}
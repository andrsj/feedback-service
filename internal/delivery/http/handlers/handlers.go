package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrsj/feedback-service/internal/services/feedback"
	"github.com/andrsj/feedback-service/pkg/logger"
)

type Handlers interface {
	Status(w http.ResponseWriter, r *http.Request)
	Token(w http.ResponseWriter, r *http.Request)
	GetFeedback(w http.ResponseWriter, r *http.Request)
	GetAllFeedback(w http.ResponseWriter, r *http.Request)
	CreateFeedback(w http.ResponseWriter, r *http.Request)
}

func New(service feedback.Service, logger logger.Logger) *handlers {
	return &handlers{
		logger:          logger.Named("handlers"),
		feedbackService: service,
	}
}

type handlers struct {
	logger          logger.Logger
	feedbackService feedback.Service
}

var _ Handlers = (*handlers)(nil)

func (h *handlers) Status(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Hit status endpoint", nil)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ok")
}

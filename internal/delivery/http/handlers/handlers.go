package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/andrsj/feedback-service/internal/domain/models"
	"github.com/andrsj/feedback-service/internal/services/feedback"
	"github.com/andrsj/feedback-service/pkg/logger"
)

type Service interface {
	Create(feedback *models.FeedbackInput) (string, error)
	GetByID(feedbackID string) (*models.Feedback, error)
	GetAll() ([]*models.Feedback, error)
	GetPage(limit int, next string) ([]*models.Feedback, string, error)
}

// Check if the actual implementation fits the interface.
var _ Service = (*feedback.Service)(nil)

type Handlers struct {
	logger          logger.Logger
	feedbackService Service
}

func New(service Service, logger logger.Logger) *Handlers {
	return &Handlers{
		logger:          logger.Named("handlers"),
		feedbackService: service,
	}
}

func (h *Handlers) Status(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Hit status endpoint", nil)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ok")
}

func (h *Handlers) handleError(w http.ResponseWriter, statusCode int, err error) {
	h.logger.Error("handler error", logger.M{"err": err})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}) //nolint:errchkjson
}

// ChatGPT's generated code for testing graceful shutdown.
func (h *Handlers) FakeLongWork(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	sleepTime, err := strconv.Atoi(queryValues.Get("time"))
	if err != nil || sleepTime <= 0 {
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid time value"}) //nolint:errchkjson

		return
	}

	time.Sleep(time.Duration(sleepTime) * time.Second)
	w.Write([]byte("Ok")) //nolint:errcheck
}

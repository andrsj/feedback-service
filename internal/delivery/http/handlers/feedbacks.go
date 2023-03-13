package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/andrsj/feedback-service/internal/domain/models"
)

var (
	errIDParamIsMissing = errors.New("id parameter is missing")
)

// GetFeedback GET /feedback/{id}.
func (h *handlers) GetFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	feedbackID := chi.URLParam(r, "id")
	if feedbackID == "" {
		err := errIDParamIsMissing
		h.handleError(w, http.StatusBadRequest, err)

		return
	}

	feedback, err := h.feedbackService.GetByID(feedbackID)
	if err != nil {
		h.handleError(w, http.StatusNotFound, err)

		return
	}

	err = json.NewEncoder(w).Encode(feedback)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, err)

		return
	}
}

// GetAllFeedback GET /feedbacks.
func (h *handlers) GetAllFeedback(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	feedbacks, err := h.feedbackService.GetAll()
	if err != nil {
		h.handleError(w, http.StatusBadRequest, err)

		return
	}

	err = json.NewEncoder(w).Encode(feedbacks)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, err)

		return
	}
}

// CreateFeedback POST /feedback.
func (h *handlers) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback models.FeedbackInput

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		h.handleError(w, http.StatusBadRequest, err)

		return
	}

	feedbackID, err := h.feedbackService.Create(&feedback)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, err)

		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"id": feedbackID})
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, err)

		return
	}
}

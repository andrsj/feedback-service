package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andrsj/feedback-service/internal/domain/models"
	"github.com/andrsj/feedback-service/pkg/logger"
	"github.com/go-chi/chi/v5"
)

var (
	errIDParamIsMissing  = errors.New("id parameter is missing")
)

// GetFeedback GET /feedback/{id}.
func (h *handlers) GetFeedback(w http.ResponseWriter, r *http.Request) {
	feedbackID := chi.URLParam(r, "id")
	if feedbackID == "" {
		err := errIDParamIsMissing
		h.logger.Error(err.Error(), logger.M{"err": err})
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	feedback, err := h.feedbackService.GetByID(feedbackID)
	if err != nil {
		h.logger.Error(err.Error(), logger.M{"err": err})
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	err = json.NewEncoder(w).Encode(feedback)
	if err != nil {
		h.logger.Error(err.Error(), logger.M{"err": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// GetAllFeedback GET /feedbacks.
func (h *handlers) GetAllFeedback(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	feedbacks, err := h.feedbackService.GetAll()
	if err != nil {
		// TODO Format error
		h.logger.Error(err.Error(), logger.M{"err": err})
		http.Error(w, err.Error(), http.StatusBadRequest)

		return 
	}


	err = json.NewEncoder(w).Encode(feedbacks)
	if err != nil {
		h.logger.Error(err.Error(), logger.M{"err": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// CreateFeedback POST /feedback.
func (h *handlers) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback models.FeedbackInput

	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		h.logger.Error(err.Error(), logger.M{"err": err})
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	feedbackID, err := h.feedbackService.Create(&feedback)
	if err != nil {
		// TODO check the error
		h.logger.Error(err.Error(), logger.M{"err": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(map[string]string{"id": feedbackID})
	if err != nil {
		h.logger.Error(err.Error(), logger.M{"err": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

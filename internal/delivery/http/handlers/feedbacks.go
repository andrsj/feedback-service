package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/andrsj/feedback-service/internal/domain/models"
)

const (
	defaultLimit    = 10
	limitQueryParam = "limit"
	nextQueryParam  = "next"
)

var (
	errIDParamIsMissing = errors.New("id parameter is missing")
	errLimitParam       = errors.New("invalid limit parameter")
	errNextParam        = errors.New("invalid next parameter")
	errNoValuesNext     = errors.New("no values after 'next'")
)

// GetFeedback GET /feedback/{id}.
func (h *Handlers) GetFeedback(w http.ResponseWriter, r *http.Request) {
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
func (h *Handlers) GetAllFeedback(w http.ResponseWriter, _ *http.Request) {
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
func (h *Handlers) CreateFeedback(w http.ResponseWriter, r *http.Request) {
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

// GetPageFeedbacks GET /p-feedbacks.
func (h *Handlers) GetPageFeedbacks(w http.ResponseWriter, r *http.Request) {
	var (
		limit      int
		nextInput  string
		nextOutput string
		nextURL    string
		err        error
		feedbacks  []*models.Feedback
	)

	limit, nextInput, err = validatePaginator(r.URL.Query())
	if err != nil {
		h.handleError(w, http.StatusBadRequest, err)
	}

	feedbacks, nextOutput, err = h.feedbackService.GetPage(limit, nextInput)
	if err != nil {
		h.handleError(w, http.StatusBadRequest, err)

		return
	}

	if len(feedbacks) == 0 {
		err = fmt.Errorf("next '%s': %w", nextInput, errNoValuesNext)
		h.handleError(w, http.StatusBadRequest, err)

		return
	}

	nextURL = fmt.Sprintf("/p-feedbacks?%s=%d&%s=%s", limitQueryParam, limit, nextQueryParam, nextOutput)
	w.Header().Set("URL-cursor-next", nextURL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(feedbacks)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, err)

		return
	}
}

func validatePaginator(queryParams url.Values) (int, string, error) {
	var (
		err   error
		limit = defaultLimit //nolint:ineffassign
		next  string
	)

	limit, err = checkLimit(queryParams)
	if err != nil {
		return 0, "", fmt.Errorf("error while check limit: %w", err)
	}

	next, err = checkNext(queryParams)
	if err != nil {
		return 0, "", fmt.Errorf("error while check next: %w", err)
	}

	return limit, next, nil
}

func checkNext(queryParams url.Values) (string, error) {
	var (
		err  error
		next = "" //nolint:ineffassign
	)

	next = queryParams.Get(nextQueryParam)
	if next != "" {
		_, err = uuid.Parse(next)
		if err != nil {
			return next, fmt.Errorf("wrong format of next ID: %w", errNextParam)
		}
	}

	return next, nil
}

func checkLimit(queryParams url.Values) (int, error) {
	var (
		err   error
		limit = defaultLimit
	)

	limitStr := queryParams.Get(limitQueryParam)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, fmt.Errorf("wrong limit param '%s': %w", limitStr, errLimitParam)
		}

		if limit <= 0 {
			return 0, fmt.Errorf("wrong limit param '%d < 0': %w", limit, errLimitParam)
		}
	}

	return limit, nil
}

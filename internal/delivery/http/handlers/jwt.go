package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/andrsj/feedback-service/pkg/logger"
)

const (
	defaultMinutes    = 10
	defaultRole       = "all"
	minutesQueryParam = "minutes"
	roleQueryParam    = "role"
	tokenPrefix       = "Bearer"
)

var (
	errRoleParam    = errors.New("invalid role parameter")
	errMinutesParam = errors.New("invalid minutes parameter")
)

func (h *handlers) Token(w http.ResponseWriter, r *http.Request) {
	var (
		minutes int64
		role    string
		err     error
	)

	queryParams := r.URL.Query()

	h.logger.Info("Hit Token endpoint", logger.M{"queryParams": queryParams})

	minutes, err = checkMinutes(queryParams)
	if err != nil {
		h.handleError(w, http.StatusBadRequest, fmt.Errorf("error while checking minutes: %w", err))

		return
	}

	role, err = checkRole(queryParams)
	if err != nil {
		h.handleError(w, http.StatusBadRequest, fmt.Errorf("error while checking role: %w", err))

		return
	}

	h.logger.Info("Received data", logger.M{
		"role": role,
		"time": minutes,
	})

	token, err := generateJWTToken(minutes, role)
	if err != nil {
		h.handleError(w, http.StatusInternalServerError, fmt.Errorf("can't create a token: %w", err))

		return
	}

	fullToken := fmt.Sprintf("%s %s", tokenPrefix, token)

	h.logger.Info("Successfully returning created token", logger.M{"token": fullToken})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fullToken)
}

func checkRole(queryParams url.Values) (string, error) {
	role := queryParams.Get(roleQueryParam)
	if role != "" {
		switch role {
		case "get", "post", "all":
			return role, nil
		default:
			return "", fmt.Errorf("wrong role '%s': %w", role, errRoleParam)
		}
	}

	return defaultRole, nil
}

func checkMinutes(queryParams url.Values) (int64, error) {
	var (
		err     error
		minutes int64 = defaultMinutes
	)

	minutesStr := queryParams.Get(minutesQueryParam)
	if minutesStr != "" {
		minutes, err = strconv.ParseInt(minutesStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("wrong minutes param: %w", errMinutesParam)
		}

		if minutes <= 0 {
			return 0, fmt.Errorf("wrong value for minutes param '%d': %w", minutes, errMinutesParam)
		}
	}

	return minutes, nil
}

func generateJWTToken(minutes int64, role string) (string, error) {
	const (
		expiredAtKey = "expiredAt"
		roleKey      = "role"
	)

	secret := []byte(os.Getenv("SECRET"))

	expirationTime := time.Now().Add(time.Duration(minutes) * time.Minute)

	claims := jwt.MapClaims{
		expiredAtKey: expirationTime.Unix(),
		roleKey:      role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

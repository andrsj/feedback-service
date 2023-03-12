package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// TODO Refactor this shit
func (h *handlers) Token(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Hit status endpoint", nil)

	token, err := generateJWTToken()
	if err != nil {
		http.Error(w, "Can't create a token :(", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Baerer %s", token)
}

func generateJWTToken() (string, error) {
	secret := []byte(os.Getenv("SECRET"))

	expirationTime := time.Now().Add(1 * time.Minute)

	claims := jwt.MapClaims{
		"expiredAt": expirationTime.Unix(),
		"role":      "all",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/andrsj/feedback-service/internal/infrastructure/cache"
	"github.com/andrsj/feedback-service/pkg/logger"

)

const (
	authorizationHeader = "Authorization"
	tokenPrefix         = "Bearer"
)

var (
	errMissingAuthHeader   = errors.New("missing authorization header")
	errInvalidAuthHeader   = errors.New("invalid authorization header, use 'Bearer <Token>'")
	errSigningMethod       = errors.New("unexpected signing method")
	errTokenExpired        = errors.New("token is expired")
	errTokenMissingExpired = errors.New("token missing expired value")
	errTokenRole           = errors.New("token has wrong role")
	errTokenMissingRole    = errors.New("token missing role value")
	errInvalidToken        = errors.New("invalid token")
)

func CacheMiddleware(cache cache.Cache) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)

				return
			}

			// TODO Header

			cacheKey := r.URL.String()
			if val, ok := cache.Get(cacheKey); ok {
				w.Write(val) //nolint

				return
			}

			rw := NewResponseWriter(w, http.StatusProcessing)
			next.ServeHTTP(rw, r)

			if rw.Status() == http.StatusOK {
				cache.Set(cacheKey, rw.Body.Bytes())
			}
		})
	}
}

func JWTMiddleware(log logger.Logger) func(next http.Handler) http.Handler {
	var errMSG string

	log = log.Named("jwt")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := validateHeaders(r.Header)
			if err != nil {
				errMSG = "Wrong authorization header"
				log.Error(errMSG, logger.M{"err": err})
				http.Error(w, errMSG, http.StatusBadRequest)

				return
			}

			token, err := parseToken(tokenStr)
			if err != nil {
				errMSG = "Wrong token authorization"
				log.Error(errMSG, nil)
				http.Error(w, errMSG, http.StatusUnauthorized)

				return
			}

			err = validateToken(token, r)
			if err != nil {
				errMSG = "Validating token error"
				log.Error(errMSG, logger.M{"err": err})
				http.Error(w, errMSG, http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong alg '%v': %w", token.Header["alg"], errSigningMethod)
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return token, fmt.Errorf("parsing error: %w", err)
	}

	return token, nil
}

func validateHeaders(header http.Header) (string, error) {
	const lenBearerToken = 2

	authHeader := header.Get(authorizationHeader)
	if authHeader == "" {
		return "", errMissingAuthHeader
	}

	// token format: "Bearer_<Token>"
	splittedBearerToken := strings.Split(authHeader, " ")
	if len(splittedBearerToken) != lenBearerToken {
		return "", errInvalidAuthHeader
	}

	if splittedBearerToken[0] != tokenPrefix {
		return "", errInvalidAuthHeader
	}

	return splittedBearerToken[1], nil
}

func validateToken(token *jwt.Token, r *http.Request) error {
	var err error

	claims, claimsIsValid := token.Claims.(jwt.MapClaims)
	if !claimsIsValid && !token.Valid {
		return errInvalidToken
	}

	err = validateExpiredAt(claims)
	if err != nil {
		return fmt.Errorf("validating 'expiredAt' error: %w", err)
	}

	err = validateRole(claims, r.Method)
	if err != nil {
		return fmt.Errorf("validating 'role' error: %w", err)
	}

	return nil
}

func validateRole(claims jwt.MapClaims, httpMethod string) error {
	role, roleIsValid := claims["role"].(string)
	if !roleIsValid {
		return fmt.Errorf("%w", errTokenMissingRole)
	}

	switch role {
	case "get":
		if httpMethod != http.MethodGet {
			return fmt.Errorf("wrong role for 'GET': %w", errTokenRole)
		}
	case "post":
		if httpMethod != http.MethodPost {
			return fmt.Errorf("wrong role for 'POST': %w", errTokenRole)
		}
	case "all":
		return nil
	default:
		return fmt.Errorf("not existing role role: %w", errTokenRole)
	}

	return nil
}

func validateExpiredAt(claims jwt.MapClaims) error {
	expiredAt, expiredAtIsValid := claims["expiredAt"].(float64)
	if !expiredAtIsValid {
		return fmt.Errorf("%w", errTokenMissingExpired)
	}

	if int64(expiredAt) < time.Now().Unix() {
		return fmt.Errorf("expired (%d): %w", time.Now().Unix()-int64(expiredAt), errTokenExpired)
	}

	return nil
}

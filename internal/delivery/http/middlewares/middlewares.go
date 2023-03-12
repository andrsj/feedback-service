package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/andrsj/feedback-service/internal/infrastructure/cache"
	"github.com/andrsj/feedback-service/pkg/logger"
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

// TODO Refactor this shit
func JWTMiddleware(log logger.Logger) func(next http.Handler) http.Handler {
	log = log.Named("jwt")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			log.Info("Authorization", logger.M{"token": authHeader})

			if authHeader == "" {
				log.Error("Missing authorization header", nil)
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)

				return
			}

			// if authHeader != os.Getenv("SECRET") {
			// 	log.Error("Wrong authorization header", nil)
			// 	http.Error(w, "Wrong authorization header", http.StatusUnauthorized)

			// 	return
			// }

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 {
				log.Error("Invalid authorization header", logger.M{"token": authHeader})
				http.Error(w, "Invalid authorization header, use 'Bearer <Token>'", http.StatusUnauthorized)

				return
			}

			tokenStr := bearerToken[1]

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(os.Getenv("SECRET")), nil
			})

			if err != nil {
				log.Error("Wrong authorization header", nil)
				http.Error(w, "Wrong authorization header", http.StatusUnauthorized)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if expiredAt, ok := claims["expiredAt"].(float64); ok {

					log.Info("expiredAt", logger.M{"expiredAt": expiredAt})
					log.Info("expiredAt", logger.M{"expiredAt": int64(expiredAt)})

					if int64(expiredAt) < time.Now().Unix() {
						log.Error("Token is expired", logger.M{"expiredAt": expiredAt})
						http.Error(w, "token is expired", http.StatusUnauthorized)

						return
					}
				} else {
					log.Error("ExpiredAt is missing", nil)
					http.Error(w, "ExpiredAt is missing", http.StatusUnauthorized)

					return
				}

				if role, ok := claims["role"].(string); ok {

					log.Info("role", logger.M{"role": role})

					switch role {
					case "reading":
						if r.Method != "GET" {
							log.Error("Token has invalid role", logger.M{"role": role})
							http.Error(w, "token has invalid role", http.StatusUnauthorized)

							return
						}
					case "creating":
						if r.Method != "POST" {
							log.Error("Token has invalid role", logger.M{"role": role})
							http.Error(w, "token has invalid role", http.StatusUnauthorized)

							return
						}
					case "all":
						// allow everything
					default:
						log.Error("Token has invalid role", logger.M{"role": role})
						http.Error(w, "token has invalid role", http.StatusUnauthorized)

						return
					}

				} else {
					log.Error("Role is missing", nil)
					http.Error(w, "role is missing", http.StatusUnauthorized)

					return
				}

			} else {
				log.Error("Invalid token claims", nil)
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
		})
	}
}

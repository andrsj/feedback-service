package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/andrsj/feedback-service/internal/delivery/http/middlewares"
	"github.com/andrsj/feedback-service/internal/infrastructure/cache"
	"github.com/andrsj/feedback-service/pkg/logger"
)

type Router struct {
	router *chi.Mux
	logger logger.Logger

	cacheMiddleware func(next http.Handler) http.Handler
	jwtMiddleware   func(next http.Handler) http.Handler
}

func New(cache cache.Cache, logger logger.Logger) *Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	jwtMiddleware := middlewares.JWTMiddleware(logger)
	cacheMiddleware := middlewares.CacheMiddleware(cache)

	return &Router{
		router:          router,
		logger:          logger.Named("router"),
		cacheMiddleware: cacheMiddleware,
		jwtMiddleware:   jwtMiddleware,
	}
}

type Handlers interface {
	Status(w http.ResponseWriter, r *http.Request)
	Token(w http.ResponseWriter, r *http.Request)
	GetFeedback(w http.ResponseWriter, r *http.Request)
	GetAllFeedback(w http.ResponseWriter, r *http.Request)
	CreateFeedback(w http.ResponseWriter, r *http.Request)
	GetPageFeedbacks(w http.ResponseWriter, r *http.Request)
	FakeLongWork(w http.ResponseWriter, r *http.Request)
}

func (r *Router) Register(handler Handlers) {
	r.logger.Info("Registering handlers", nil)

	// Status checker.
	r.router.Get("/", handler.Status)
	
	// Token generation.
	r.router.Get("/token", handler.Token)
	
	// No cache all feedbacks.
	r.router.With(r.jwtMiddleware).Get("/feedbacks", handler.GetAllFeedback)
	r.router.Group(
		func(router chi.Router) {
			router.Use(r.cacheMiddleware)
			router.Use(r.jwtMiddleware)
			
			// Specific ID.
			router.Get("/feedback/{id}", handler.GetFeedback)
			// Paginated cursor list of feedbacks.
			router.Get("/p-feedbacks", handler.GetPageFeedbacks)
			// Create feedback.
			router.Post("/feedback", handler.CreateFeedback)
		},
	)

	// Testing router for checking Graceful Shutdown.
	r.router.Get("/l", handler.FakeLongWork)

	// ChatGPT's generated code for logging registered <Method: URLs>
	err := chi.Walk(
		r.router,
		func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
			r.logger.Info(fmt.Sprintf("%-5s -> %s", method, route), nil)

			return nil
		},
	)
	if err != nil {
		panic(err)
	}
}

func (r *Router) GetChiMux() *chi.Mux {
	return r.router
}

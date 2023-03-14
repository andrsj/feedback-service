package router

import (
	"fmt"
	"net/http"
	"time"

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
}

func (r *Router) Register(handler Handlers) {
	r.logger.Info("Registering handlers", nil)

	r.router.Get("/token", handler.Token)
	
	r.router.Group(
		func(router chi.Router) {
			router.Use(r.cacheMiddleware)
			router.Use(r.jwtMiddleware)
			
			router.Get("/", handler.Status)
			router.Get("/feedbacks", handler.GetAllFeedback)
			router.Get("/feedback/{id}", handler.GetFeedback)
			router.Get("/p-feedbacks", handler.GetPageFeedbacks)
			
			router.Post("/feedback", handler.CreateFeedback)
		},
	)

	r.router.Get("/l", func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(10 * time.Second)
		w.Write([]byte("Ok"))
	})

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

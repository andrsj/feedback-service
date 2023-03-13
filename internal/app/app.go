package app

import (
	"fmt"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/andrsj/feedback-service/internal/delivery/http/handlers"
	"github.com/andrsj/feedback-service/internal/delivery/http/router"
	"github.com/andrsj/feedback-service/internal/delivery/http/server"
	"github.com/andrsj/feedback-service/internal/infrastructure/cache/memcached"
	repo "github.com/andrsj/feedback-service/internal/infrastructure/db/gorm"
	"github.com/andrsj/feedback-service/internal/services/feedback"
	log "github.com/andrsj/feedback-service/pkg/logger"
)

type App struct {
	server *http.Server
	logger log.Logger
}

type Params struct {
	DsnDB            string
	CacheSecondsLive int32
	CacheHost        string
	Logger           log.Logger
}

func New(params *Params) (*App, error) {
	logger := params.Logger.Named("app")

	//nolint
	db, err := gorm.Open(
		postgres.Open(params.DsnDB),
		&gorm.Config{},
	)
	if err != nil {
		logger.Error("Can't connect to DB", log.M{"err": err, "dsn": params.DsnDB})

		return nil, fmt.Errorf("can't connect to DB: %w", err)
	}

	feedbackRepo, err := repo.NewFeedbackRepository(db, logger)
	if err != nil {
		logger.Error("Can't up repository", log.M{"err": err})

		return nil, fmt.Errorf("can't up repository: %w", err)
	}

	service := feedback.New(feedbackRepo, logger)
	handlers := handlers.New(service, logger)

	cache := memcached.New(params.CacheHost, params.CacheSecondsLive, logger)
	router := router.New(cache, logger)
	router.Register(handlers)

	server := server.New(router)

	return &App{
		server: server,
		logger: logger,
	}, nil
}

func (a *App) Start() error {
	a.logger.Info("Starting the application", log.M{"address": a.server.Addr})

	err := a.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("can't listen and serve: %w", err)
	}

	return nil
}

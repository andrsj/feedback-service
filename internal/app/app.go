package app

import (
	"fmt"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/andrsj/feedback-service/internal/delivery/http/handlers"
	"github.com/andrsj/feedback-service/internal/delivery/http/router"
	"github.com/andrsj/feedback-service/internal/delivery/http/server"
	"github.com/andrsj/feedback-service/internal/infrastructure/cache/memory"
	repo "github.com/andrsj/feedback-service/internal/infrastructure/db/gorm"
	// TODO remove this shit
	// repo "github.com/andrsj/feedback-service/internal/infrastructure/db/memory"
	"github.com/andrsj/feedback-service/internal/services/feedback"
	log "github.com/andrsj/feedback-service/pkg/logger"
)

type App struct {
	server *http.Server
	logger log.Logger
}

func New(dsn string, logger log.Logger) (*App, error) {
	logger = logger.Named("app")

	// business logic
	// feedbackRepo := repo.New(logger)
	
	//nolint
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	) 
	if err != nil {
		logger.Error("Can't connect to DB", log.M{"err": err, "dsn": dsn})

		return nil, fmt.Errorf("can't connect to DB: %w", err)
	}

	feedbackRepo, err := repo.NewFeedbackRepository(db, logger)
	if err != nil {
		logger.Error("Can't up repository", log.M{"err": err})

		return nil, fmt.Errorf("can't up repository: %w", err)
	}

	service := feedback.New(feedbackRepo, logger)
	handlers := handlers.New(service, logger)

	cache := memory.New(logger)
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

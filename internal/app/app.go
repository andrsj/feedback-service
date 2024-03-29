package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/andrsj/feedback-service/internal/delivery/http/handlers"
	"github.com/andrsj/feedback-service/internal/delivery/http/router"
	"github.com/andrsj/feedback-service/internal/delivery/http/server"
	"github.com/andrsj/feedback-service/internal/infrastructure/broker/kafka"
	"github.com/andrsj/feedback-service/internal/infrastructure/cache/memcached"
	repo "github.com/andrsj/feedback-service/internal/infrastructure/db/gorm"
	"github.com/andrsj/feedback-service/internal/services/feedback"
	log "github.com/andrsj/feedback-service/pkg/logger"
)

const timeoutShutdown = 5

type App struct {
	server *http.Server
	logger log.Logger
}

type Params struct {
	DsnDB            string
	CacheSecondsLive int32
	CacheHost        string
	KafkaHost        string
	KafkaTopic       string
	Logger           log.Logger
}

func New(params *Params) (*App, error) {
	logger := params.Logger.Named("app")

	//nolint:varnamelen
	db, err := gorm.Open(
		postgres.Open(params.DsnDB),
		//nolint:exhaustivestruct,exhaustruct
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		},
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

	broker, err := kafka.New(logger, params.KafkaHost, params.KafkaTopic)
	if err != nil {
		logger.Error("Can't up broker", log.M{"err": err})

		return nil, fmt.Errorf("can't up broker: %w", err)
	}

	service := feedback.New(feedbackRepo, broker, logger)
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

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Server error", log.M{"error": err.Error()})
		}
	}()

	sig := <-osSignals

	a.logger.Info("Received signal", log.M{"signal": sig})

	ctx, cancel := context.WithTimeout(context.Background(), timeoutShutdown*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server shutdown error", log.M{"error": err.Error()})
	}

	a.logger.Info("Application stopped", nil)

	return nil
}

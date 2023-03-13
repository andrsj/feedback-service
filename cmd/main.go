package main

import (
	"fmt"
	"os"

	"github.com/andrsj/feedback-service/internal/app"
	log "github.com/andrsj/feedback-service/pkg/logger"
	zap "github.com/andrsj/feedback-service/pkg/logger/zap"
	"github.com/joho/godotenv"
)

func main() {
	zap := zap.New()

	err := godotenv.Load("config.env")
	if err != nil {
		msg := fmt.Sprintf("Error loading .env file: %s", err)
		zap.Fatal(msg, nil)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	app, err := app.New(dsn, zap)
	if err != nil {
		zap.Fatal("Can't configure the app", log.M{"err": err})
	}

	err = app.Start()
	if err != nil {
		zap.Fatal("Can't start the application", log.M{"err": err})
	}
}

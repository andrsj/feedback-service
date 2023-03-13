package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/andrsj/feedback-service/internal/app"
	log "github.com/andrsj/feedback-service/pkg/logger"
	zap "github.com/andrsj/feedback-service/pkg/logger/zap"

)

func main() {
	zap := zap.New()

	err := godotenv.Load("config.env")
	if err != nil {
		zap.Fatal("error loading .env file", log.M{"err": err})
	}

	// Postgresql config
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	// Memcached config
	memcachedHost := fmt.Sprintf(
		"%s:%s",
		os.Getenv("MEMCACHED_HOST"),
		os.Getenv("MEMCACHED_PORT"),
	)
	memcachedSecondsLiveStr := os.Getenv("MEMCACHED_LIVE_TIME")
	memcachedSecondsLive, err := strconv.Atoi(memcachedSecondsLiveStr)

	if err != nil {
		zap.Fatal("can't convert the Memcached live time seconds into integer", log.M{"err": err})
	}

	// Kafka config
	kafkaHost := fmt.Sprintf(
		"%s:%s",
		os.Getenv("KAFKA_HOST"),
		os.Getenv("KAFKA_PORT"),
	)
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	// App creating
	app, err := app.New(&app.Params{
		DsnDB:            dsn,
		CacheSecondsLive: int32(memcachedSecondsLive),
		CacheHost:        memcachedHost,
		KafkaHost:        kafkaHost,
		KafkaTopic:       kafkaTopic,
		Logger:           zap,
	})
	if err != nil {
		zap.Fatal("can't configure the app", log.M{"err": err})
	}

	// App starting
	err = app.Start()
	if err != nil {
		zap.Fatal("can't start the application", log.M{"err": err})
	}
}

package main

import (
	"flag"
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

	var configFile string
	
	flag.StringVar(&configFile, "c", "", "path to config file (required)")
	flag.Parse()

	if configFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := godotenv.Load(configFile)
	if err != nil {
		zap.Fatal("error loading .env file", log.M{"err": err})
	}

	// Postgresql config
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	zap.Info("DB Configuration", log.M{
		"host": dbHost,
		"port": dbPort,
		"user": dbUser,
		"pass": dbPass,
		"name": dbName,
	})

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	zap.Info("DB URL", log.M{"dsn": dsn})

	// Memcached config
	memcachedHost := fmt.Sprintf(
		"%s:%s",
		os.Getenv("MEMCACHED_HOST"),
		os.Getenv("MEMCACHED_PORT"),
	)
	memcachedSecondsLiveStr := os.Getenv("MEMCACHED_LIVE_TIME")
	memcachedSecondsLive, err := strconv.Atoi(memcachedSecondsLiveStr)

	zap.Info("Memcached data", log.M{
		"URL": memcachedHost,
		"live time": memcachedSecondsLiveStr,
	})

	if err != nil {
		zap.Fatal("can't convert the Memcached live time seconds into integer", log.M{"err": err})
	}

	// Kafka config
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaURL := fmt.Sprintf(
		"%s:%s",
		kafkaHost,
		kafkaPort,
	)
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	
	zap.Info("Apache Kafka Configuration", log.M{
		"host": kafkaHost,
		"port": kafkaPort,
		"topic": kafkaTopic,
	})

	// App creating
	app, err := app.New(&app.Params{
		DsnDB:            dsn,
		CacheSecondsLive: int32(memcachedSecondsLive),
		CacheHost:        memcachedHost,
		KafkaHost:        kafkaURL,
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

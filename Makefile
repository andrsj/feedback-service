BINARY_NAME=app

include config.env

.PHONY: build brun clean
build:
	go build -o build/${BINARY_NAME} cmd/main.go

brun: build
	./build/${BINARY_NAME}

clean:
	go clean
	rm build/${BINARY_NAME}

.PHONY: up down go bgo stop
up:
	docker-compose --env-file=config.env up -d

down:
	docker-compose --env-file=config.env down

bgo:
	docker build -t my-golang-app .
	docker-compose --env-file=config.env up -d my-golang-app
	docker image prune -f

go:
	docker-compose --env-file=config.env up -d my-golang-app

stop:
	docker-compose --env-file=config.env stop my-golang-app
	docker-compose --env-file=config.env rm -f my-golang-app
	docker image rm my-golang-app
	docker image prune -f

.PHONY: db-up db-down
db-up:
	docker-compose --env-file=config.env up -d postgresql

db-down:
	docker-compose --env-file=config.env stop postgresql
	docker-compose --env-file=config.env rm -f postgresql

.PHONY: cache-up cache-down
cache-up:
	docker-compose --env-file=config.env up -d memcached

cache-down:
	docker-compose --env-file=config.env stop memcached
	docker-compose --env-file=config.env rm -f memcached

.PHONY: kafka-up kafka-down consumer
kafka-up:
	docker-compose --env-file=config.env up -d kafka

kafka-down:
	docker-compose --env-file=config.env stop kafka
	docker-compose --env-file=config.env stop zookeeper
	docker-compose --env-file=config.env rm -f kafka
	docker-compose --env-file=config.env rm -f zookeeper

consumer:
	docker exec -it \
		kafka kafka-console-consumer.sh \
		--bootstrap-server $(KAFKA_HOST):$(KAFKA_PORT) \
		--topic $(KAFKA_TOPIC) \
		--from-beginning

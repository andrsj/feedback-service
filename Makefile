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

.PHONY: up down
up:
	docker-compose --file=docker-compose.yml --env-file=config.env up -d

down:
	docker-compose --file=docker-compose.yml --env-file=config.env down

.PHONY: db-up db-down
db-up:
	docker-compose --file=postgresql.yml --env-file=config.env up -d

db-down:
	docker-compose --file=postgresql.yml --env-file=config.env down

.PHONY: cache-up cache-down
cache-up:
	docker-compose --file=memcached.yml --env-file=config.env up -d

cache-down:
	docker-compose --file=memcached.yml --env-file=config.env down

.PHONY: kafka-up kafka-down consumer
kafka-up:
	docker-compose --file=kafka.yml --env-file=config.env up -d

kafka-down:
	docker-compose --file=kafka.yml --env-file=config.env down

consumer:
	docker exec -it \
		kafka kafka-console-consumer.sh \
		--bootstrap-server $(KAFKA_HOST):$(KAFKA_PORT) \
		--topic $(KAFKA_TOPIC) \
		--from-beginning

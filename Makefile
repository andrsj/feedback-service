BINARY_NAME=app

.PHONY: build brun clean db-up db-down

build:
	go build -o build/${BINARY_NAME} cmd/main.go

brun: build
	./build/${BINARY_NAME}

clean:
	go clean
	rm build/${BINARY_NAME}

db-up:
	docker-compose --file=postgresql.yml --env-file=config.env up -d

db-down:
	docker-compose --file=postgresql.yml --env-file=config.env down
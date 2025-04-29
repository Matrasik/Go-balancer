.PHONY: up
up:
	go mod tidy
	docker-compose up -d
down:
	docker-compose down

build:
	docker-compose build

restart:
	docker-compose down
	docker-compose up -d
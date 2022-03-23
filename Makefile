#!make
up_all:
	docker-compose up

up_db:
	docker-compose up -d db

start_srv:
	go run cmd/main.go

run_test:
	go test -v ./...

up_mac: up_db start_srv
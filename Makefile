.PHONY:
.SILENT:

build:
	docker compose build

up:
	docker compose up

bash:
	docker compose run --service-ports web bash

run:
	#air cmd/main.go -b 0.0.0.0
	air cmd/main.go

test:
	go test ./...
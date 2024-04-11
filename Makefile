.PHONY: build
build:
	go build -v ./cmd/apiserver

.DEFAULT-GOAL := build

.PHONY: run
run:
	docker-compose --env-file ./config/environ/db.env up -d

.PHONY: run-ex
run-ex:
	docker-compose --env-file ./config/environ/db.env up

.PHONY: run-rebuild
run-rebuild:
	docker-compose --env-file ./config/environ/db.env up --build --force-recreate

.PHONY: down
down:
	docker-compose down

.PHONY: down-v
down-v:
	docker-compose down --volumes

.PHONY: stop
stop:
	docker-compose stop
.PHONY: build
build:
	go build -v ./cmd/apiserver

.DEFAULT-GOAL := build

.PHONY: run
run:
	docker-compose --env-file ./config/environ/db.env -f docker-compose.yaml up -d

.PHONY: run-ex
run-ex:
	docker-compose --env-file ./config/environ/db.env -f docker-compose.yaml up

.PHONY: run-rebuild
run-rebuild:
	docker-compose --env-file ./config/environ/db.env -f docker-compose.yaml up --build --force-recreate

.PHONY: test
test:
	docker-compose --env-file ./test/config/environ/db.env -f docker-compose-test.yaml up --build -d
	go test -v "test/main_test.go" "test/get_banner_test.go" && docker-compose --env-file ./test/config/environ/db.env -f docker-compose-test.yaml stop

.PHONY: down
down:
	docker-compose -f docker-compose.yaml down

.PHONY: down-tests
down-tests:
	docker-compose -f docker-compose-test.yaml down

.PHONY: down-v
down-v:
	docker-compose -f docker-compose.yaml down --volumes

.PHONY: stop
stop:
	docker-compose -f docker-compose.yaml stop

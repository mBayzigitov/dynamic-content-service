run:
	docker-compose up -d

run-ex:
	docker-compose up

run-rebuild:
	docker-compose up --build --force-recreate

down:
	docker-compose down

down-v:
	docker-compose down --volumes

stop:
	docker-compose stop
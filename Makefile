include .env

stop_containers:
	@echo "Stopping other docker containers"
	if [ $$(docker ps -q) ]; then \
		echo "Found and stopped containers"; \
		docker stop $$(docker ps -q); \
	else \
		echo "No containers running..."; \
	fi

create_container:
	docker run --name ${DB_DOCKER_CONTAINER} -p ${DB_PORT}:5432 -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:15-alpine

create_db:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

start_container:
	docker start ${DB_DOCKER_CONTAINER}

migrate_up:
	migrate -path db/migrations -database ${DSN} -verbose up

migrate_down:
	migrate -path db/migrations -database ${DSN} -verbose down

sqlc_generate:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: stop_containers create_container create_db start_container migrate_up migrate_down sqlc_generate test
include .env
export

.PHONY: mongodb-up mongodb-down postgres-up postgres-down run docker-up docker-down mongod nginx-start nginx-stop postres-start postgres-stop

mongodb-up:
	@migrate -database 'mongodb://${MONGODB_HOST}:${MONGODB_PORT}/${MONGODB_DB}' -path migrations/mongodb up

mongodb-down:
	@migrate -database 'mongodb://${MONGODB_HOST}:${MONGODB_PORT}/${MONGODB_DB}' -path migrations/mongodb down

postgres-up:
	@migrate -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}' -path migrations/postgres up

postgres-down:
	@migrate -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}' -path migrations/postgres down

run:
	@go run cmd/main.go

docker-up:
	@docker-compose up --build

docker-down:
	@docker-compose down -v

mongod:
	@cd /d ${MONGOD_PATH} && mongod.exe

nginx-start:
	@cd /d ${NGINX_PATH} && start nginx.exe -c ${LOCAL_CONF_PATH}/local.conf

nginx-stop:
	@cd /d ${NGINX_PATH} && nginx.exe -s stop

postgres-start:
	@net start ${POSTGRES_SRV}

postgres-stop:
	@net stop ${POSTGRES_SRV} /y
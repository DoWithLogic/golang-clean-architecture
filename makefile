.PHONY: help database-up database-down migration-up migration-down local run

help:
	@echo "Available targets:"
	@echo "  make database-up   	- Start the database container"
	@echo "  make database-down 	- Stop and remove the database container"
	@echo "  make migration-up  	- Run database migrations"
	@echo "  make migration-down 	- Rollback database migrations"
	@echo "  make local         	- Run the application locally"
	@echo "  make run           	- Start the database, run migrations, and start the application locally"
	@echo "  make down           	- Shutdown the database and down migrations"



# Directory where migration files are located
MIGRATION_DIR := database/mysql/migration

# This target waits for the MySQL container to become available
wait-for-mysql:
	@echo "Waiting for MySQL container to start..."
	@until docker compose exec mysql-db mysql -umysql -ppwd -hlocalhost -e "SELECT 1"; do \
		sleep 6; \
	done
	@echo "MySQL is up and running!"

database-up: 
	docker compose up mysql-db -d

service-up:
	docker compose up golang-clean-architecture -d

docker-down:
	docker compose down 

migration-up: wait-for-mysql
	GOOSE_DRIVER=mysql GOOSE_DBSTRING="mysql:pwd@tcp(localhost:3306)/users?parseTime=true" goose -dir=$(MIGRATION_DIR) up

migration-down: 
	GOOSE_DRIVER=mysql GOOSE_DBSTRING="mysql:pwd@tcp(localhost:3306)/users?parseTime=true" goose -dir=$(MIGRATION_DIR) down


run: database-up migration-up service-up

down : migration-down docker-down

mock-repository:
	mockgen -source internal/users/repository/repository.go -destination internal/users/mock/repository_mock.go -package=mocks

mock-usecase:
	mockgen -source internal/users/usecase/usecase.go -destination internal/users/mock/usecase_mock.go -package=mocks


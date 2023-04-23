help:
	@echo ''
	@echo 'Usage: make [TARGET] [EXTRA_ARGUMENTS]'
	@echo 'Targets:'
	@echo 'make dev: make dev for development work'
	@echo 'make build: make build container'
	@echo 'clean: clean for all clear docker images'

dev:
	docker-compose down
	docker-compose up

build:
	docker-compose down
	docker-compose build

clean:
	docker-compose down -v

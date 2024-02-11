include .env

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

psql:
	psql paycrest

gen-ent:
	go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema/

run: gen-ent
	air

seed:
	go run main.go -seed-db

test:
	go test -v ./...

start-rpc:
	ganache -m "$(HD_WALLET_MNEMONIC)" --chain.chainId 1337 -l 21000000

test-coverage:
	go test $(go list ./... | grep -v /ent | grep -v /config | grep -v /database | grep -v /routers)  -coverprofile=coverage.out ./...

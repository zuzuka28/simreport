generate_api:
	oapi-codegen --generate gorilla,types,strict-server api/rest/doc/openapi.yaml > api/rest/gen/server.gen.go

wire:
	wire ./...

startup_env:
	docker-compose -f build/docker-compose.yaml up -d

shutdown_env:
	docker-compose -f build/docker-compose.yaml down --remove-orphans

run_api:
	go run cmd/simrep/*.go --config config/config.yaml run-api

run_async_api:
	go run cmd/simrep/*.go --config config/config.yaml run-async-api

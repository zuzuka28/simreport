generate_api:
	go tool oapi-codegen --generate gorilla,types,strict-server,embedded-spec api/rest/openapi.yaml > internal/handler/rest/gen/server.gen.go

generate_proto:
	protoc proto/similarity.proto --go_out=. --natsmicro_out=.

wire:
	go tool wire ./...

startup_env:
	docker-compose -f build/docker-compose.yaml up -d

shutdown_env:
	docker-compose -f build/docker-compose.yaml down --remove-orphans

run_api:
	go run cmd/similarity/*.go --config config/config.yaml run-api

run_async_processing:
	go run cmd/similarity/*.go --config config/config.yaml run-async-processing

run_app:
	go run cmd/similarity/*.go --config config/config.yaml

format_code:
	go run mvdan.cc/gofumpt@latest -w .

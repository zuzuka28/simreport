generate_api:
	oapi-codegen --generate gorilla,types,strict-server api/rest/doc/openapi.yaml > api/rest/gen/server.gen.go

wire:
	wire ./...


run_api:
	go run cmd/anysave/*.go --config config/config.yaml run-api

run_app:
	go run cmd/anysave/*.go --config config/config.yaml

wire:
	wire ./...

generate_proto:
	protoc proto/fulltextindex.proto --go_out=. --natsmicro_out=.

run_api:
	go run cmd/fulltextindex/*.go --config config/config.yaml run-api

run_async_processing:
	go run cmd/fulltextindex/*.go --config config/config.yaml run-async-processing

run_app:
	go run cmd/fulltextindex/*.go --config config/config.yaml

format_code:
	go run mvdan.cc/gofumpt@latest -w .

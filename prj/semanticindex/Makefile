wire:
	wire ./...

run_api:
	go run cmd/semanticindex/*.go --config config/config.yaml run-api

run_async_processing:
	go run cmd/semanticindex/*.go --config config/config.yaml run-async-processing

run_app:
	go run cmd/semanticindex/*.go --config config/config.yaml

format_code:
	go run mvdan.cc/gofumpt@latest -w .

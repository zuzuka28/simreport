ARG GO_VERSION=1.24

FROM golang:${GO_VERSION} AS build

RUN apt update && apt install -y git  ca-certificates && update-ca-certificates

WORKDIR /app

COPY . .

RUN export CGO_ENABLED=0 && \
    go test ./... && \
    go build -installsuffix 'static' -o app ./cmd/allinone

FROM gcr.io/distroless/static AS final

COPY --from=build /app .

ENTRYPOINT [ "/app" ]

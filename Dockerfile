# Step 1: Modules caching
FROM golang:1.26-alpine3.23 AS modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Step 2: Swagger docs generation
FROM golang:1.26-alpine3.23 AS swagger

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN go tool swag init --parseDependency -g internal/controller/restapi/router.go

# Step 3: Builder
FROM golang:1.26-alpine3.23 AS builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app

# Step 4: Final
FROM scratch

COPY --from=builder /app/config /config
COPY --from=swagger /app/docs /docs
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app"]

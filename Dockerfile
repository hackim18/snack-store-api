FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/web

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/api/openapi.yaml /app/api/openapi.yaml
COPY --from=builder /app/internal/migrations/json /app/internal/migrations/json
COPY .env.example /app/.env.example

EXPOSE 8080

USER app

ENTRYPOINT ["/app/server"]
CMD ["--drop-table", "--migrate", "--seed", "--run"]

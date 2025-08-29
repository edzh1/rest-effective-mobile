FROM golang:1.24.6-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN swag init -g cmd/main.go --parseDependency --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

COPY --from=builder /app/migrations ./migrations

RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

RUN chown -R appuser:appuser /root

USER appuser
EXPOSE 8080

CMD ["./main"]
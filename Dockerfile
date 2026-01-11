FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gochat_app ./cmd/gochat

FROM alpine:3
WORKDIR /app
COPY --from=builder /app/gochat_app .
CMD [ "/app/gochat_app" ]

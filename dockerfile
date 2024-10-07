# Use Go official image
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init --dir ./cmd,./internal/server --generalInfo main.go --output ./docs/openapi


ARG ENV
ARG SERVER_PORT
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_DATABASE
ARG POSTGRES_HOST
ARG POSTGRES_PORT

EXPOSE ${SERVER_PORT}

CMD ["go", "run", "./cmd/main.go"]

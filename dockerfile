FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

#RUN swag init  # RUN swag init  # Generate Swagger docs during the build process

COPY . .

ARG ENV
ARG SERVER_PORT
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_DATABASE
ARG POSTGRES_HOST
ARG POSTGRES_PORT


EXPOSE ${SERVER_PORT}

CMD ["go", "run", "./cmd/main.go"]

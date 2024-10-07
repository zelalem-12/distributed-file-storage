# Distributed File Storage Server

Distributed File Storage Server is a Go-based web application that allows users to upload, retrieve, and download files. The server uses PostgreSQL for storage and supports parallel file downloads for better performance. Swagger is used to generate the API documentation.

## Features

. File Upload: Upload multiple files via multipart form data.
. File Metadata Storage: Store file metadata (name, size, type) in PostgreSQL.
. Parallel Download: Efficient parallel file download mechanism.
. API Documentation: Automatically generated Swagger documentation.
. Docker Support: Easily deployable with Docker.

## Technologies Used

. Go (Golang)
. PostgreSQL
. Swag (for Swagger documentation)
. Docker
. Httprouter (for routing)

## Table of Contents

1. [Installation](#installation)
2. [Running Locally](#running-locally)
3. [Docker Setup](#docker-setup)
4. [API Endpoints](#api-endpoints)
5. [Swagger Documentation](#swagger-documentation)
6. [Environment Variables](#environment-variables)

## Installation

### Prerequisites

. Go 1.23+
. PostgreSQL 14+
. Docker (optional, for containerized deployment)

### Setup

1. Clone the repository:

```bash
git clone https://github.com/your-username/distributed-file-storage.git
cd distributed-file-storage
```

2. Install dependencies:

```bash
go mod download
```

3. Set up your .env file:

```bash
cp .env.example .env
# Fill out the required environment variables (POSTGRES credentials, SERVER_PORT, etc.)
```

## Running Locally

1. Ensure PostgreSQL is running and the database is set up with the proper credentials.

2. Start the server:

```bash
go run ./cmd/main.go
```

The server should now be running on the port specified in your .env file.

## Docker Setup

### Build Docker Image and Run Docker container

1. Build the Docker image:

```bash
docker compose up --build
```

## API Endpoints

## Upload Files

1. Endpoint: /api/v1/upload
2. Method: POST
3. Description: Upload one or multiple files
4. Request: Multipart form data, field name: files

## List All Files

1. Endpoint: /api/v1/files
2. Method: GET
3. Description: Get metadata of all uploaded files.

## Download File by ID

1. Endpoint: /api/v1/downloads/:id
2. Method: GET
3. Description: Download the file using its unique ID.

# Swagger Documentation

The API documentation is generated using Swagger and is available at:

```bash
http://localhost:8080/swagger/index.html
```

To regenerate the Swagger documentation while running without docker, run:

```bash
swag init --dir ./cmd,./internal/server --generalInfo main.go --output ./docs/openapi
```

## Environment Variables

1. SERVER_PORT: The port the server listens on.
2. POSTGRES_USER: PostgreSQL username.
3. POSTGRES_PASSWORD: PostgreSQL password.
4. POSTGRES_DATABASE: PostgreSQL database name.
5. POSTGRES_HOST: PostgreSQL host (e.g., localhost).
6. POSTGRES_PORT: PostgreSQL port (e.g., 5432).
7. ENV: Application environment (development, production).

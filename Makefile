.PHONY: run test build docker-build docker-run lint

# Run locally (requires WEATHER_API_KEY env var)
run:
	go run ./cmd/server

# Run all tests
test:
	go test ./... -v

# Run tests with coverage
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Build the binary
build:
	go build -o server ./cmd/server

# Build Docker image
docker-build:
	docker build -t cep-weather .

# Run Docker container (requires WEATHER_API_KEY env var)
docker-run:
	docker run -p 8080:8080 -e WEATHER_API_KEY=$(WEATHER_API_KEY) cep-weather

# Run with Docker Compose
compose-up:
	docker-compose up --build

# Stop Docker Compose
compose-down:
	docker-compose down

# CEP Weather API

A Go microservice that receives a Brazilian ZIP code (CEP), identifies the city, and returns the current temperature in Celsius, Fahrenheit, and Kelvin.

## 🌐 Live URL (Cloud Run)

> **Replace this with your actual Cloud Run URL after deploying:**
> `https://cep-weather-XXXXXXXXXX-uc.a.run.app`

Example request:
```
GET https://cep-weather-XXXXXXXXXX-uc.a.run.app/01310100
```

---

## 📋 API Contract

### `GET /{cep}`

| Scenario | HTTP Code | Response |
|---|---|---|
| Success | 200 | `{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5}` |
| Invalid CEP format | 422 | `invalid zipcode` |
| CEP not found | 404 | `can not find zipcode` |

**Rules:**
- CEP must be exactly **8 numeric digits** (no dashes)
- Temperatures: `C`, Fahrenheit `F = C * 1.8 + 32`, Kelvin `K = C + 273`

---

## 🚀 Running Locally

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) (optional)
- A free [WeatherAPI](https://www.weatherapi.com/) key

### ⚡ Quick Start — 3 Steps (2 min)

**1. Clone and navigate:**
```bash
git clone https://github.com/joaogarbeline/cep-weather.git
cd cep-weather
```

**2. Set your API key (choose your shell):**

PowerShell:
```powershell
$env:WEATHER_API_KEY = "your_key_api"
# For permanent setup (new shell required):
setx WEATHER_API_KEY "your_key_api"
```

Bash/Mac/Linux:
```bash
export WEATHER_API_KEY=your_key_api
```

**3. Run and test:**
```bash
# Terminal 1: Start server (keeps running on port 8080)
go run ./cmd/server

# Terminal 2: Test the API
curl http://localhost:8080/01310100
```

**Expected response:**
```json
{"temp_C":19,"temp_F":66.2,"temp_K":292}
```

---

### Option 1 — Go directly

```bash
# Clone the repository
git clone https://github.com/joaogarbeline/cep-weather.git
cd cep-weather

# Set your API key
export WEATHER_API_KEY=your_key_api

# Run the server
go run ./cmd/server

# Test it
curl http://localhost:8080/01310100
```

### Option 2 — Docker

```bash
# Build the image
docker build -t cep-weather .

# Run the container
docker run -p 8080:8080 -e WEATHER_API_KEY=your_key_api cep-weather

# Test it
curl http://localhost:8080/01310100
```

### Option 3 — Docker Compose

```bash
# Copy and fill in your API key
cp .env.example .env
# Edit .env and set WEATHER_API_KEY=your_key_api

# Start the service
docker-compose up --build

# Test it
curl http://localhost:8080/01310100
```

---

## 🧪 Running Tests

**All tests (no API key required):**
```bash
go test ./... -v
```

**With coverage report:**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**By package:**
```bash
go test ./internal/handler -v       # HTTP handler tests
go test ./internal/service -v       # Business logic tests
go test ./internal/integration -v   # Integration tests
```

**Using Make:**
```bash
make test              # Run all tests
make test-cover        # Run with coverage report
```

### Test coverage includes:
- **Unit tests** — temperature conversion formulas (°C → °F → K)
- **Unit tests** — CEP validation (length, numeric check)
- **Unit tests** — HTTP handler responses (200, 404, 422, 405)
- **Integration tests** — full pipeline using mock HTTP servers

---

## ☁️ Deploying to Google Cloud Run

### Prerequisites

- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) installed and authenticated
- A GCP project with billing enabled (free tier works)

### Deploy

```bash
# Make the script executable
chmod +x deploy.sh

# Deploy (replace with your GCP project ID and WeatherAPI key)
./deploy.sh my-gcp-project-id your_weatherapi_key

# The script outputs your service URL at the end
```

### Manual deploy steps

```bash
export PROJECT_ID=my-gcp-project-id
export WEATHER_API_KEY=your_key_api
export SERVICE_NAME=cep-weather
export REGION=us-central1
export IMAGE=gcr.io/${PROJECT_ID}/${SERVICE_NAME}

# Enable APIs
gcloud services enable run.googleapis.com containerregistry.googleapis.com

# Build and push
gcloud builds submit --tag ${IMAGE}

# Deploy
gcloud run deploy ${SERVICE_NAME} \
  --image ${IMAGE} \
  --platform managed \
  --region ${REGION} \
  --allow-unauthenticated \
  --set-env-vars "WEATHER_API_KEY=${WEATHER_API_KEY}" \
  --port 8080
```

---

## 🏗️ Project Structure

```
cep-weather/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── client/
│   │   ├── viacep.go            # ViaCEP API client
│   │   └── weatherapi.go        # WeatherAPI client
│   ├── handler/
│   │   ├── handler.go           # HTTP handler
│   │   └── handler_test.go      # Handler unit tests
│   ├── service/
│   │   ├── weather.go           # Business logic + conversions
│   │   └── weather_test.go      # Service unit tests
│   └── integration/
│       └── integration_test.go  # Integration tests
├── Dockerfile
├── docker-compose.yml
├── deploy.sh
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## 🔑 External APIs

| API | Purpose | Docs |
|---|---|---|
| [ViaCEP](https://viacep.com.br/) | CEP → City lookup | Free, no key required |
| [WeatherAPI](https://www.weatherapi.com/) | City → Temperature | Free tier: 1M calls/month |

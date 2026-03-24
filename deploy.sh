#!/bin/bash
# deploy.sh — Deploy to Google Cloud Run
# Usage: ./deploy.sh <PROJECT_ID> <WEATHER_API_KEY>
# Example: ./deploy.sh my-gcp-project abc123key

set -euo pipefail

PROJECT_ID="${1:?Usage: ./deploy.sh <PROJECT_ID> <WEATHER_API_KEY>}"
WEATHER_API_KEY="${2:?Usage: ./deploy.sh <PROJECT_ID> <WEATHER_API_KEY>}"

SERVICE_NAME="cep-weather"
REGION="us-central1"
IMAGE="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

echo "==> Configuring project: ${PROJECT_ID}"
gcloud config set project "${PROJECT_ID}"

echo "==> Enabling required APIs..."
gcloud services enable \
  run.googleapis.com \
  containerregistry.googleapis.com

echo "==> Building and pushing Docker image..."
gcloud builds submit --tag "${IMAGE}"

echo "==> Deploying to Cloud Run..."
gcloud run deploy "${SERVICE_NAME}" \
  --image "${IMAGE}" \
  --platform managed \
  --region "${REGION}" \
  --allow-unauthenticated \
  --set-env-vars "WEATHER_API_KEY=${WEATHER_API_KEY}" \
  --port 8080

echo ""
echo "==> Deployment complete!"
echo "Service URL:"
gcloud run services describe "${SERVICE_NAME}" \
  --platform managed \
  --region "${REGION}" \
  --format "value(status.url)"

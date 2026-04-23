#!/bin/bash

set -euo pipefail

ENV_FILE="${ENV_FILE:-.env}"
if [[ -f "${ENV_FILE}" ]]; then
  set -o allexport
  source "${ENV_FILE}"
  set +o allexport
else
  echo "Warning: .env file not found at '${ENV_FILE}', using defaults" >&2
fi

START_TIME=$(date +%s)
START_DATETIME=$(date "+%Y-%m-%d %H:%M:%S")
echo "Deploy script initiated at: ${START_DATETIME}"

# Configuration
APP_NAME="${APP_NAME:-svt}"
BUILD_PATH="${BUILD_PATH:-./cmd/ssh/main.go}"
OUTPUT_BIN="${OUTPUT_BIN:-${APP_NAME}}"
REMOTE_USER="${REMOTE_USER:-root}"
REMOTE_HOST="${REMOTE_HOST:?REMOTE_HOST is not set in environment or .env}"
REMOTE_PORT="${REMOTE_PORT:?REMOTE_PORT is not set in environment or .env}"
REMOTE_DIR="${REMOTE_DIR:-/usr/local/bin}"
REMOTE_TARGET="${REMOTE_DIR}/${APP_NAME}"
REMOTE_BACKUP="${REMOTE_TARGET}.bak.$(date +%Y%m%d%H%M%S)"

cleanup() {
  rm -f "${OUTPUT_BIN}"
}
trap cleanup EXIT

echo "Building Linux amd64 binary: ${OUTPUT_BIN}"
GOOS=linux GOARCH=amd64 go build -ldflags="-X main.buildMode=prod" -o "${OUTPUT_BIN}" "${BUILD_PATH}"

if [[ ! -f "${OUTPUT_BIN}" ]]; then
  echo "Build failed: ${OUTPUT_BIN} was not created" >&2
  exit 1
fi

echo "Uploading binary to ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_TARGET}"
scp -P "${REMOTE_PORT}" "${OUTPUT_BIN}" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_TARGET}.new"

echo "Activating new binary on remote host"
ssh -p "${REMOTE_PORT}" "${REMOTE_USER}@${REMOTE_HOST}" "
  set -euo pipefail
  if [[ -f '${REMOTE_TARGET}' ]] || [[ -L '${REMOTE_TARGET}' ]]; then
    mv '${REMOTE_TARGET}' '${REMOTE_BACKUP}'
  fi
  mv '${REMOTE_TARGET}.new' '${REMOTE_BACKUP}'
  ln -sf '${REMOTE_BACKUP}' '${REMOTE_TARGET}'
  chmod 755 '${REMOTE_BACKUP}'
  systemctl restart ${APP_NAME}
"

END_TIME=$(date +%s)
DEPLOYMENT_TIME=$((END_TIME - START_TIME))
MINUTES=$((DEPLOYMENT_TIME / 60))
SECONDS=$((DEPLOYMENT_TIME % 60))
echo "Deployment completed in ${MINUTES}m ${SECONDS}s"
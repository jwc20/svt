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

REMOTE_USER="${REMOTE_USER:-root}"
REMOTE_HOST="${REMOTE_HOST:?REMOTE_HOST is not set in environment or .env}"
REMOTE_PORT="${REMOTE_PORT:?REMOTE_PORT is not set in environment or .env}"
APP_NAME="${APP_NAME:-svt-app}"
REMOTE_DIR="${REMOTE_DIR:-/usr/local/bin}"

echo "Rolling back ${APP_NAME} on ${REMOTE_HOST}..."

ssh -p "${REMOTE_PORT}" "${REMOTE_USER}@${REMOTE_HOST}" "
  set -euo pipefail
  LATEST=\$(ls -t ${REMOTE_DIR}/${APP_NAME}.bak.* | head -1)
  if [[ -z \"\${LATEST}\" ]]; then
    echo 'No backup found' >&2
    exit 1
  fi
  echo \"Rolling back to: \${LATEST}\"
  ln -sf \"\${LATEST}\" '${REMOTE_DIR}/${APP_NAME}'
  systemctl restart ${APP_NAME}
"

echo "Rollback complete."
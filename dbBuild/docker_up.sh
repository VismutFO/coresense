#!/bin/bash
set -e
SCRIPT_DIRECTORY=$(dirname "$(realpath "$0")")
source "${SCRIPT_DIRECTORY}/config.env"

INFRA_COUNT=$(docker network ls | grep ${INFRA_NETWORK_NAME} | wc -l)
if [[ "${INFRA_COUNT}" -eq "0" ]]; then
  docker network create -d overlay --attachable ${INFRA_NETWORK_NAME}
fi

echo "Docker-compose up..."
docker compose -p ${INFRA_PROJECT_NAME} \
  -f docker-compose.db.yml \
  up -d
echo "Docker-compose is up"

#!/bin/bash
set -e
SCRIPT_DIRECTORY=$(dirname "$(realpath "$0")")
source "${SCRIPT_DIRECTORY}/config.env"

echo Stopping all running threads...
# pushd . > /dev/null
# cd ../../admin
# ./admin.sh stopall
# popd > /dev/null

echo "Docker-compose down..."
docker compose -p ${INFRA_PROJECT_NAME} \
  -f docker-compose.db.yml \
  down

docker volume prune --force
echo "Docker-compose is cleaned up"

echo "Cleaning up volumes..."
volumes=( \
  infra_postgres-data \
)

for volume in ${volumes[@]}; do
  if docker volume inspect $volume > /dev/null 2>&1; then
    echo Removing volume $volume...
    docker volume rm $volume > /dev/null
  fi
done
echo "Volumes are cleaned up"

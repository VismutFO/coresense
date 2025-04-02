#!/usr/bin/env bash
set -e

echo "Executing init.sh"

start=`date +%s`

SCRIPT_DIRECTORY=$(dirname "$(realpath "$0")")

# Check if Docker Swarm is already initialized
if ! docker info | grep -q 'Swarm: active'; then
    echo "Initializing Docker Swarm..."
    docker swarm init --advertise-addr 127.0.0.1
else
    echo "Docker Swarm is already initialized."
fi

# Check if network 'infra' already exists
if ! docker network ls | grep -q 'infra'; then
    echo "Creating network 'infra'..."
    docker network create --driver overlay --attachable infra
else
    echo "Network 'infra' already exists."
fi

source "$SCRIPT_DIRECTORY/docker_down.sh"
source "$SCRIPT_DIRECTORY/docker_up.sh"

containers=( \
  liquibase-claimix-postgres \
)

for container_name in ${containers[@]}; do
  if [ $(docker wait $container_name) -ne 0 ]; then
    echo "Container '$container_name' has failed"
    docker logs $container_name
    exit $exit_code
  else
    echo "Container '$container_name' has finished"
    docker container rm $container_name > /dev/null # remove only if migration was successful
  fi
done

echo "Apply fixtures..."
( cd $SCRIPT_DIRECTORY/../db/ ; source ./02_run_fixtures.sh )

end=$(date +%s)
echo "Done: $((end-start))s"
#!/bin/bash
set -e

EXTENSIONS=$(cat <<-EOSQL
    CREATE EXTENSION pg_stat_statements;
EOSQL
)

echo "Creating extensions in claimix database..."
psql -U postgres -d claimix -v ON_ERROR_STOP=1 <<< ${EXTENSIONS}

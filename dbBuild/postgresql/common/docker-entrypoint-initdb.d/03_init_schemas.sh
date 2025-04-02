#!/bin/bash
set -e

echo "creating database..."
psql -U postgres -d claimix -v ON_ERROR_STOP=1 <<-EOSQL
        CREATE SCHEMA claimix;
EOSQL

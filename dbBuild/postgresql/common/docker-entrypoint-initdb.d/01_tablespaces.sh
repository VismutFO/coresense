#!/bin/bash
set -e

echo "01_tablespace claimix"

for tablespace in claimix_entities claimix_param claimix_oper; do
  echo "creating $tablespace tablespace..."
        mkdir -pv "/var/lib/postgresql/data/$tablespace"
done

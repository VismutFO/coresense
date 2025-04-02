#!/usr/bin/env bash

declare -A db_schema
db_schema["claimix"]="claimix"

for db in "${!db_schema[@]}"; do
  # "schema1 schema2 schema3" -> ('schema1','schema2','schema3')
  schemas="($(echo ${db_schema[$db]} | sed "s/ /','/g" | sed "s/^/'/" | sed "s/$/'/"))"
  query="select count(*) from information_schema.schemata where schema_name in $schemas;"
  count=$(echo $query | psql -t -U postgres -d $db | xargs)

  expected=$(echo ${db_schema[$db]} | wc -w)
  if [[ $count -ne $expected ]]; then
    exit 1
  fi
done

exit 0
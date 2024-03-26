#!/bin/sh
sh ./build.sh

docker-compose up --no-start
docker-compose run --rm db_migrate_job sh migrate.sh
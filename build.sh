#!/bin/sh
docker-compose --env-file=.env build --no-cache $@
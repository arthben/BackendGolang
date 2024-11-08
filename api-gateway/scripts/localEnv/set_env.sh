#!/bin/bash

# DB_HOST="postgresdb"
DB_HOST="127.0.0.1"
DB_DRIVER="postgres"
DB_USER="epiphyte"
DB_PASSWORD="@3p1phyte-corp.com"
DB_NAME="backendGo"
DB_PORT="5432"
DB_MIN_POOL="1"
DB_MAX_POOL="10"

DB_DSN="${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&timezone=UTC"

export DB_HOST
export DB_DRIVER
export DB_USER
export DB_PASSWORD
export DB_NAME
export DB_PORT
export DB_MIN_POOL
export DB_MAX_POOL
export DB_DSN


# set ini di env docker
SERVER_PORT=3000
SERVER_TIMEOUT=30

export SERVER_PORT
export SERVER_TIMEOUT


# OPENWEATHER_APIKEY=598eb70eacf1d53a1eec6ef3e6da25c2
OPENWEATHER_APIKEY=e4d25c020947523c7c18b8e4af1ce00e

export OPENWEATHER_APIKEY

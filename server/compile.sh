#!/usr/bin/env bash

export SANCHO_ENV="dev"
export HTTP_PORT=8081
export DB_PATH="database/database.sancho"
export SANCHO_PATH="/home/alejandro/StreamripDownloads/sancho"
export CGO_ENABLED=1
export FRONTEND_PATH="../client/build"

if go build -o ./bin ./cmd/sancho; then
  echo "Build exitoso."
else
  echo "Error al compilar."
  exit 1
fi

./bin/sancho-api

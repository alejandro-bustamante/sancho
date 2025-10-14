#!/usr/bin/env bash

export SANCHO_ENV="dev"
export HTTP_PORT=5400
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null)"
export DB_PATH="$REPO_ROOT/test/test_env/sancho_test.db"
export SANCHO_PATH="$REPO_ROOT/test/test_env/"
export CGO_ENABLED=1
export FRONTEND_PATH="../client/build"

if go build -o ./bin ./cmd/sancho; then
  echo "Build exitoso."
else
  echo "Error al compilar."
  exit 1
fi

./bin/sancho

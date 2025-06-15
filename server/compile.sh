#!/usr/bin/env bash
#go build -o ./bin ./cmd/sancho-api

# Establecer la variable de entorno CGO_ENABLED solo para este comando
export CGO_ENABLED=1

# Intentar construir con la etiqueta icu. Si falla, mostrar un mensaje.
if go build -tags "icu" -o ./bin ./cmd/sancho-api; then
#if go build -o ./bin ./cmd/sancho-api; then
  echo "Build exitoso con soporte para go-sqlite3 e ICU."
else
  echo "Error al construir con soporte para ICU. Asegúrate de tener libicu-dev (o equivalente) instalado en tu sistema."
  exit 1
fi

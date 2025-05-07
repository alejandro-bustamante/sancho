#!/usr/bin/env bash

# Intentar construir con la etiqueta icu. Si falla, mostrar un mensaje.
go build -o ./bin ./cmd/sancho-indexer;

echo "Script de build completado."

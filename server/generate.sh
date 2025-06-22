#!/usr/bin/env bash
set -e
echo "🛠 Generando código con sqlc..."
sqlc generate
echo "✅ Código generado en internal/db/"

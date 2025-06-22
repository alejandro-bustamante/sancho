#!/usr/bin/env bash
set -e
DB_PATH="database/dev.sancho"

echo "🧱 Migrando base de datos: $DB_PATH"
migrate -database sqlite3://$DB_PATH -path migrations up
echo "✅ Migración completada."

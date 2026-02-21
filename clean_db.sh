#!/usr/bin/env bash
set -e
DB_PATH="database/database.sancho"
BACKUP_DB_PATH="database/database.sancho.backup"

echo "Eliminando la base de datos usada"
rm $DB_PATH
echo "Cambiando la base de datos backup a la nueva"
mv $BACKUP_DB_PATH $DB_PATH
echo "Respaldando la base de datos mediante una copia"
cp $DB_PATH $BACKUP_DB_PATH
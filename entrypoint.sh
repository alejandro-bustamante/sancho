#!/bin/sh
set -e

CONFIG_DIR="/root/.config/streamrip"
CONFIG_FILE="$CONFIG_DIR/config.toml"
TEMPLATE_FILE="$CONFIG_DIR/config.template.toml"

echo "Configurando streamrip con variables de entorno..."

# Verificar que el template existe
if [ ! -f "$TEMPLATE_FILE" ]; then
    echo "Error: Template de configuración no encontrado en $TEMPLATE_FILE"
    exit 1
fi

# Copiar el template al archivo final
cp "$TEMPLATE_FILE" "$CONFIG_FILE"

# QOBUZ PASSWORD
if [ -n "$QOBUZ_PASSWORD_OR_TOKEN" ]; then
    echo "Configurando token de Qobuz..."
    sed -i "s|password_or_token = \".*\"|password_or_token = \"$QOBUZ_PASSWORD_OR_TOKEN\"|g" "$CONFIG_FILE"
else
    echo "Advertencia: QOBUZ_PASSWORD_OR_TOKEN no está configurado"
fi

# QOBUZ USER_ID
if [ -n "$QOBUZ_USER_ID" ]; then
    echo "Configurando User ID de Qobuz..."
    sed -i "s|email_or_userid = \".*\"|email_or_userid = \"$QOBUZ_USER_ID\"|g" "$CONFIG_FILE"
else
    echo "Advertencia: QOBUZ_USER_ID no está configurado"
fi

# QOBUZ APP_ID
if [ -n "$QOBUZ_APP_ID" ]; then
    echo "Configurando App ID de Qobuz..."
    sed -i "s|app_id = \".*\"|app_id = \"$QOBUZ_APP_ID\"|g" "$CONFIG_FILE"
else
    echo "Advertencia: QOBUZ_APP_ID no está configurado"
fi

# QOBUZ SECRETS
if [ -n "$QOBUZ_SECRETS" ]; then
    echo "Configurando secrets de Qobuz..."
    # Reemplaza comas por "," para formatear como array de strings en TOML
    FORMATTED_SECRETS=$(echo "$QOBUZ_SECRETS" | sed 's/,/","/g')
    sed -i "s|secrets = \[.*\]|secrets = \[\"$FORMATTED_SECRETS\"\]|g" "$CONFIG_FILE"
else
    echo "Advertencia: QOBUZ_SECRETS no está configurado"
fi


echo "Configuración completada. Iniciando sancho..."

# Ejecutar el comando principal (sancho) con todos los argumentos pasados
exec sancho "$@"

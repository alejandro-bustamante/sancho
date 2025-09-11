# =========================
# Etapa 1: Build Go binary
# =========================
FROM golang:1.24.4-alpine3.22 AS go-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server ./
RUN go build -o /bin/sancho ./cmd/sancho

# =========================
# Etapa 2: Preparar streamrip
# =========================
FROM python:3.12-alpine3.22 AS streamrip-builder
RUN apk add --no-cache git gcc musl-dev libffi-dev openssl-dev curl
# Instalar Poetry
RUN curl -sSL https://install.python-poetry.org | python3 -
ENV PATH="/root/.local/bin:$PATH"
# Instalar plugin de exportación
RUN poetry self add poetry-plugin-export
# Clonar tu fork
WORKDIR /src
RUN git clone --branch sancho --depth 1 https://github.com/alejandro-bustamante/streamrip.git
WORKDIR /src/streamrip
# Exportar las dependencias a requirements.txt
RUN poetry export -f requirements.txt --without-hashes --output requirements.txt
# Crear un wheel del proyecto
RUN poetry build

# ============================
# Etapa final: Imagen mínima
# ============================
FROM alpine:3.22

# Instalar Python y herramientas necesarias
RUN apk add --no-cache \
    ffmpeg \
    libstdc++ \
    sqlite-libs \
    python3 \
    py3-pip \
    py3-cryptography \
    py3-urllib3 \
    py3-certifi \
    py3-charset-normalizer \
    py3-idna \
    py3-requests \
    ncurses \
    ncurses-terminfo \
    gcc \
    musl-dev \
    libffi-dev \
    openssl-dev

# Crear carpeta config esperada por streamrip
RUN mkdir -p /root/.config/streamrip

# Copiar template de configuración (sin tokens)
COPY config.template.toml /root/.config/streamrip/config.template.toml

# Directorio de trabajo
WORKDIR /app

# Copiar binario backend
COPY --from=go-builder /bin/sancho /usr/local/bin/sancho

# Frontend build
COPY client/build ./build

# Migraciones
COPY server/migrations ./migrations

# Copiar requirements.txt y wheel
COPY --from=streamrip-builder /src/streamrip/requirements.txt /tmp/requirements.txt
COPY --from=streamrip-builder /src/streamrip/dist/*.whl /tmp/

# Instalar dependencias de streamrip
RUN pip install --no-cache-dir --break-system-packages -r /tmp/requirements.txt
RUN pip install --no-cache-dir --break-system-packages /tmp/*.whl

# Limpiar archivos temporales y dependencias de compilación
RUN rm -rf /tmp/* && \
    apk del gcc musl-dev libffi-dev openssl-dev

# Copiar script de entrada
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Variables de entorno con valores por defecto vacíos
ENV QOBUZ_PASSWORD_OR_TOKEN=""
ENV QOBUZ_APP_ID=""

# Documentativo, no obligatorio
EXPOSE 5400

# Usar el script de entrada
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
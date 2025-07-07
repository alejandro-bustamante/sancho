# =========================
# Etapa 1: Build Go binary
# =========================
FROM golang:1.24.4-alpine3.22 AS go-builder

# Instala dependencias necesarias para SQLite + CGO
RUN apk add --no-cache gcc musl-dev

# Crear directorio de trabajo
WORKDIR /app

# Descargar dependencias
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copiar el resto del código
COPY server ./

# Compilar el binario
RUN go build -o /bin/sancho-api ./cmd/sancho-api

# ============================
# Etapa final: Imagen mínima
# ============================
FROM alpine:3.22

# Instalar solo las dependencias de ejecución
RUN apk add --no-cache ffmpeg libstdc++ sqlite-libs

# Directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el binario compilado
COPY --from=go-builder /bin/sancho-api /usr/local/bin/sancho-api

# Copiar frontend (build sveltekit ya generado previamente)
COPY client/build ./build

# Copiar las migraciones para que el binario las ejecute
COPY server/migrations ./migrations

# Entrypoint por defecto
ENTRYPOINT ["sancho-api"]

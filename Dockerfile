# =========================
# Etapa 1: Build ffmpeg estático
# =========================
FROM alpine:3.18 AS ffmpeg-builder
RUN apk add --no-cache \
    build-base \
    yasm \
    nasm \
    zlib-dev \
    zlib-static \
    curl \
    xz

WORKDIR /tmp
RUN curl -L https://ffmpeg.org/releases/ffmpeg-8.0.tar.xz -o ffmpeg.tar.xz && \
    tar xf ffmpeg.tar.xz && \
    mv ffmpeg-8.0 ffmpeg

WORKDIR /tmp/ffmpeg
RUN ./configure \
    --prefix=/opt/ffmpeg \
    --disable-shared \
    --enable-static \
    --pkg-config-flags="--static" \
    --extra-ldflags="-static" \
    --disable-debug \
    --disable-doc \
    --disable-ffplay \
    --disable-ffprobe \
    --disable-network \
    --disable-autodetect \
    --disable-everything \
    --enable-protocol=file \
    --enable-demuxer=flac \
    --enable-demuxer=mp3 \
    --enable-demuxer=image2 \
    --enable-decoder=flac \
    --enable-decoder=mp3 \
    --enable-decoder=mp3float \
    --enable-decoder=mjpeg \
    --enable-decoder=png \
    --enable-decoder=bmp \
    --enable-encoder=mjpeg \
    --enable-encoder=png \
    --enable-muxer=image2 \
    --enable-filter=scale \
    --enable-parser=mpegaudio \
    --enable-swscale \
    --enable-zlib && \
    make -j$(nproc) && \
    make install

RUN strip /opt/ffmpeg/bin/ffmpeg

# =========================
# Etapa 2: Build Frontend (Svelte)
# =========================
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY client/package.json client/package-lock.json* client/.npmrc ./
# Use `ci` for reproducible builds in CI environments
RUN npm cache clean --force && npm ci
COPY client ./
RUN npm run build

# =========================
# Etapa 3: Build Go binary
# =========================
FROM golang:1.24.4-alpine3.22 AS go-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server ./
RUN go build -o /bin/sancho ./cmd/sancho

# =========================
# Etapa 4: Preparar streamrip
# =========================
FROM python:3.12-alpine3.22 AS streamrip-builder
RUN apk add --no-cache git gcc musl-dev libffi-dev openssl-dev curl
RUN curl -sSL https://install.python-poetry.org | python3 -
ENV PATH="/root/.local/bin:$PATH"
RUN poetry self add poetry-plugin-export

WORKDIR /src
RUN git clone --branch sancho --depth 1 https://github.com/alejandro-bustamante/streamrip.git
WORKDIR /src/streamrip
RUN poetry export -f requirements.txt --without-hashes --output requirements.txt
RUN poetry build

# ============================
# Etapa final: Imagen mínima
# ============================
FROM alpine:3.22

RUN apk add --no-cache \
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

# Copiar tu ffmpeg custom en lugar del de Alpine
COPY --from=ffmpeg-builder /opt/ffmpeg/bin/ffmpeg /usr/local/bin/ffmpeg

RUN mkdir -p /root/.config/streamrip
COPY config.template.toml /root/.config/streamrip/config.template.toml

WORKDIR /app

COPY --from=go-builder /bin/sancho /usr/local/bin/sancho
COPY --from=frontend-builder /app/build ./build
COPY server/migrations ./migrations
COPY --from=streamrip-builder /src/streamrip/requirements.txt /tmp/requirements.txt
COPY --from=streamrip-builder /src/streamrip/dist/*.whl /tmp/

RUN pip install --no-cache-dir --break-system-packages -r /tmp/requirements.txt
RUN pip install --no-cache-dir --break-system-packages /tmp/*.whl

RUN rm -rf /tmp/* && \
    apk del gcc musl-dev libffi-dev openssl-dev

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENV QOBUZ_PASSWORD_OR_TOKEN=""
ENV QOBUZ_USER_ID=""

EXPOSE 5400

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
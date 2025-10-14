# syntax=docker/dockerfile:1.6

ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-bookworm AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# ─── DEV ──────────────────────────────────────────────────────────────────────
FROM base AS dev
COPY . .
# Air para live reload
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
 && chmod +x install.sh && sh install.sh && mv ./bin/air /usr/local/bin/air
EXPOSE 8080
CMD ["air"]

# ─── BUILD ────────────────────────────────────────────────────────────────────
FROM base AS builder
COPY . .
ARG ENVIRONMENT=prod
ENV ENV=${ENVIRONMENT}
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/main ./cmd/main.go

# ─── PROD ─────────────────────────────────────────────────────────────────────
FROM debian:bookworm-slim AS prod
WORKDIR /app
COPY --from=builder /out/main ./main

# (Opcional) herramientas para el actualizador de DNS
RUN apt-get update \
 && apt-get install -y curl perl libwww-perl ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Instala dinaIP (si lo necesitas dentro del contenedor)
RUN curl -fsSLo /tmp/dinaIP-consola.tar.gz \
      https://dinahosting.com/utilidades/estandar/aplicaciones/dinaIP-consola.tar.gz \
 && tar xzf /tmp/dinaIP-consola.tar.gz -C /tmp \
 && sh /tmp/dinaIP-consola/install.sh

# Vars solo en runtime (no en build args)
ENV DINAHOSTING_DOMAIN="" DINAHOSTING_USER="" DINAHOSTING_PASSWORD=""
EXPOSE 80 443
CMD sh -c 'dinaip -u "$DINAHOSTING_USER" -p "$DINAHOSTING_PASSWORD" -a "$DINAHOSTING_DOMAIN" && exec ./main'

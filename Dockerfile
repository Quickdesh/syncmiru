# build web
FROM node:18.7.0-alpine3.16 AS web-builder
WORKDIR /web
COPY web/package.json web/pnpm-lock.yaml ./
# install pnpm
RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile --prod
COPY web/ .
RUN pnpm run build

# build app
FROM golang:1.20-alpine3.16 AS app-builder

ARG VERSION=dev
ARG REVISION=dev
ARG BUILDTIME

RUN apk add --no-cache git make build-base tzdata

ENV SERVICE=syncmiru

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

COPY --from=web-builder /web/dist ./web/dist
COPY --from=web-builder /web/build.go ./web

RUN go build -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${REVISION} -X main.date=${BUILDTIME}" -o bin/syncmiru main.go

# build final image
FROM alpine:latest

LABEL org.opencontainers.image.source = "https://github/Quickdesh/SyncMiru"

ENV HOME="/config" \
XDG_CONFIG_HOME="/config" \
XDG_DATA_HOME="/config"

RUN apk add --no-cache ca-certificates curl tzdata jq

WORKDIR /app

VOLUME /config

COPY --from=app-builder /src/bin/syncmiru /usr/local/bin/

EXPOSE 8282

ENTRYPOINT ["/usr/local/bin/syncmiru", "--config", "/config"]
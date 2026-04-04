# Tunnel server (control TCP :9000 + public HTTP edge :3000)
# Build:  docker build -t devtunnel-server .
# Run:    docker run --rm -p 3000:3000 -p 9000:9000 devtunnel-server

# syntax=docker/dockerfile:1

# Match go.mod; bump if Docker Hub adds a newer tag (e.g. 1.25-alpine).
ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static binary for minimal runtime image (GOARCH follows build platform)
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/mytunneld ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 65532 -g nobody tunnel

COPY --from=build /out/mytunneld /usr/local/bin/mytunneld

USER tunnel:tunnel

EXPOSE 3000/tcp 9000/tcp

ENTRYPOINT ["/usr/local/bin/mytunneld"]

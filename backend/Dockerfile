# --- build stage ---
FROM golang:1.22-alpine AS build
WORKDIR /src

# The module has zero external dependencies (stdlib only, no go.sum),
# so this download step is a no-op today but keeps the Dockerfile
# correct if a dependency is ever added.
COPY go.mod ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 produces a fully static binary with no dynamic libc
# dependency at all (Go's own pure-Go net/DNS resolver is used instead
# of glibc's), so it runs unmodified on Alpine's musl libc in the final
# stage below.
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

# --- final stage ---
FROM alpine:3.20
RUN adduser -D -H appuser
COPY --from=build /out/api /api
USER appuser

# Matches the SERVER_PORT value docker-compose.yml sets for this
# service (see the compose file) — EXPOSE is documentation only, not
# an actual port binding, so it is hardcoded rather than templated.
EXPOSE 8080

ENTRYPOINT ["/api"]

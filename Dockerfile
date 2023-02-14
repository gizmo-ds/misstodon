FROM docker.io/library/node:18.13.0-alpine AS mfm-builder
WORKDIR /app
COPY ./internal/mfm /app/internal/mfm
COPY ./package.json /app/package.json
COPY ./pnpm-lock.yaml /app/pnpm-lock.yaml
RUN wget -qO /bin/pnpm "https://github.com/pnpm/pnpm/releases/latest/download/pnpm-linuxstatic-x64" && chmod +x /bin/pnpm
RUN pnpm install && \
    pnpm build

FROM docker.io/library/golang:1.20-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . /app
ENV CGO_ENABLED=0
COPY --from=mfm-builder /app/internal/mfm/out.js /app/internal/mfm/out.js
RUN go build -trimpath -tags timetzdata \
    -ldflags "-s -w -X github.com/gizmo-ds/misstodon/internal/global.AppVersion=$(git describe --tags --always)" \
    -o misstodon \
    cmd/misstodon/main.go

FROM docker.io/library/alpine:latest
WORKDIR /app
COPY --from=builder /app/misstodon /app/misstodon
COPY --from=builder /app/config_example.toml /app/config.toml
ENTRYPOINT ["/app/misstodon", "start"]
EXPOSE 3000

FROM docker.io/oven/bun:latest AS mfm-builder
WORKDIR /app
COPY ./internal/mfm /app/internal/mfm
COPY ./package.json /app/package.json
COPY ./bun.lockb /app/bun.lockb
RUN bun install && bun run build

FROM docker.io/library/golang:1.20-alpine AS builder
WORKDIR /app
COPY . /app
ENV CGO_ENABLED=0
ARG version=development
COPY --from=mfm-builder /app/internal/mfm/out.js /app/internal/mfm/out.js
RUN go mod download
RUN go build -trimpath -tags timetzdata \
    -ldflags "-s -w -X github.com/gizmo-ds/misstodon/internal/global.AppVersion=$version" \
    -o misstodon \
    ./cmd/misstodon

FROM gcr.io/distroless/static-debian11:latest
WORKDIR /app
COPY --from=builder /app/misstodon /app/misstodon
COPY --from=builder /app/config_example.toml /app/config.toml
ENTRYPOINT ["/app/misstodon", "start"]
EXPOSE 3000

FROM docker.io/library/node:18.13.0-alpine AS mfm-builder
WORKDIR /app
COPY ./internal/mfm /app/internal/mfm
COPY ./package.json /app/package.json
COPY ./pnpm-lock.yaml /app/pnpm-lock.yaml
RUN echo "update-notifier=false" >> /app/.npmrc
RUN corepack enable && \
    corepack prepare pnpm@7.27.0 --activate
RUN pnpm install && \
    pnpm build

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

FROM docker.io/library/alpine:latest
WORKDIR /app
COPY --from=builder /app/misstodon /app/misstodon
COPY --from=builder /app/config_example.toml /app/config.toml
ENTRYPOINT ["/app/misstodon", "start"]
EXPOSE 3000

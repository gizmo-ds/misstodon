FROM docker.io/library/node:18.13.0-alpine AS mfm-builder
WORKDIR /misstodon
COPY ./internal/mfm /misstodon/internal/mfm
COPY ./package.json /misstodon/package.json
COPY ./pnpm-lock.yaml /misstodon/pnpm-lock.yaml
RUN wget -qO /bin/pnpm "https://github.com/pnpm/pnpm/releases/latest/download/pnpm-linuxstatic-x64" && chmod +x /bin/pnpm
RUN pnpm install && \
    pnpm build

FROM docker.io/library/golang:1.19.5-alpine AS builder
RUN apk add --no-cache git
WORKDIR /misstodon
COPY . /misstodon
ENV CGO_ENABLED=0
COPY --from=mfm-builder /misstodon/internal/mfm/out.js /misstodon/internal/mfm/out.js
RUN go build -trimpath -tags timetzdata \
    -ldflags "-s -w -X github.com/gizmo-ds/misstodon/internal/global.AppVersion=$(git describe --tags --always)" \
    -o misstodon \
    cmd/misstodon/main.go

FROM docker.io/library/alpine:latest
WORKDIR /misstodon
COPY --from=builder /misstodon/misstodon /misstodon/misstodon
COPY --from=builder /misstodon/config_example.toml /misstodon/config.toml
ENTRYPOINT ["/misstodon/misstodon"]
EXPOSE 3000

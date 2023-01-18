FROM docker.io/library/golang:1.19.5-alpine AS builder
RUN mkdir -p /go/src/misstodon && \
    apk add --no-cache git
WORKDIR /go/src/misstodon
COPY . /go/src/misstodon
ENV CGO_ENABLED=0
RUN go build -trimpath -tags timetzdata \
    -ldflags "-s -w -X github.com/gizmo-ds/misstodon/internal/global.AppVersion=$(git describe --tags --always --dirty)" \
    -o misstodon \
    cmd/misstodon/main.go

FROM docker.io/library/alpine:latest
RUN mkdir -p /home/misstodon
WORKDIR /home/misstodon
COPY --from=builder /go/src/misstodon/misstodon /home/misstodon/misstodon
ENTRYPOINT ["/home/misstodon/misstodon"]
EXPOSE 3000

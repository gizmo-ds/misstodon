NAME=misstodon
OUTDIR=build
PKGNAME=github.com/gizmo-ds/misstodon
MAIN=cmd/misstodon/main.go
VERSION=$(shell git describe --tags --always)
FLAGS+=-trimpath
FLAGS+=-tags timetzdata
FLAGS+=-ldflags "-s -w -X $(PKGNAME)/internal/global.AppVersion=$(VERSION)"
export CGO_ENABLED=0

PLATFORMS := linux windows darwin

all: build-all

generate:
	go generate ./...

build-all: $(PLATFORMS)

$(PLATFORMS): generate
	GOOS=$@ GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME)-$@-amd64$(if $(filter windows,$@),.exe) $(MAIN)

sha256sum:
	cd $(OUTDIR); for file in *; do sha256sum $$file > $$file.sha256; done

zip:
	cp config_example.toml $(OUTDIR)/config.toml
	for platform in $(PLATFORMS); do \
		zip -jq9 $(OUTDIR)/$(NAME)-$$platform-amd64.zip $(OUTDIR)/$(NAME)-$$platform-amd64* $(OUTDIR)/config.toml README.md LICENSE; \
	done

clean:
	rm -rf $(OUTDIR)/*

build-image: generate
	docker build --no-cache --build-arg version=$(shell git describe --tags --always) -t ghcr.io/gizmo-ds/misstodon:latest -f Dockerfile .

build-develop-image: generate
	docker build --no-cache --build-arg version=$(shell git describe --tags --always) -t ghcr.io/gizmo-ds/misstodon:develop -f Dockerfile .

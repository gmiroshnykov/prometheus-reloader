PWD := $(shell pwd)
VERSION := $(shell git describe --tags --always)
GIT_HASH := $(shell git log -1 --pretty=format:"%H")
BUILD_LDFLAGS :=
BUILD_LDFLAGS += -X main.Version=${VERSION}
BUILD_LDFLAGS += -X main.GitHash=${GIT_HASH}

DOCKER_IMAGE ?= "laggyluke/prometheus-reloader"


.PHONY: run
run:
	go run \
		-ldflags='$(BUILD_LDFLAGS)' \
		cmd/main.go \
		-config-file testdata/prometheus.yml \
		-v 1


.PHONY: build
build:
	go build \
  -ldflags='$(BUILD_LDFLAGS)' \
  -o out/prometheus-reloader \
  cmd/main.go


.PHONY: docker-build
docker-build:
	docker build \
		--build-arg "VERSION=$(VERSION)" \
		--build-arg "GIT_HASH=$(GIT_HASH)" \
		--tag $(DOCKER_IMAGE):$(VERSION) \
		--tag $(DOCKER_IMAGE):latest \
		.


.PHONY: docker-run
docker-run:
	docker run \
		-it \
		-v $(PWD)/testdata:/etc/prometheus \
		$(DOCKER_IMAGE):latest

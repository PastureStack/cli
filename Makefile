TARGETS := $(shell ls scripts)
DAPPER_IMAGE ?= pasturestack-cli-dapper
DAPPER_HOST_ARCH ?= amd64
DAPPER_UID ?= $(shell id -u)
DAPPER_GID ?= $(shell id -g)
SOURCE_DIR := $(shell pwd)
DOCKER_BUILD_NETWORK_ARG := $(if $(DOCKER_BUILD_NETWORK),--network=$(DOCKER_BUILD_NETWORK),)

.dapper-image: Dockerfile.dapper
	docker build $(DOCKER_BUILD_NETWORK_ARG) \
		--build-arg DAPPER_HOST_ARCH=$(DAPPER_HOST_ARCH) \
		--build-arg UBUNTU_MIRROR=$(or $(UBUNTU_MIRROR),http://archive.ubuntu.com/ubuntu) \
		-t $(DAPPER_IMAGE) \
		-f Dockerfile.dapper .

$(TARGETS): .dapper-image
	docker run --rm \
		-v "$(SOURCE_DIR)":/go/src/github.com/PastureStack/cli \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-e DAPPER_UID=$(DAPPER_UID) \
		-e DAPPER_GID=$(DAPPER_GID) \
		-e ARCH=$(DAPPER_HOST_ARCH) \
		-e HOST_ARCH=$(DAPPER_HOST_ARCH) \
		-e IMAGE_NAME \
		-e REPO \
		-e TAG \
		-e VERSION_OVERRIDE \
		-e GOOS \
		-e CROSS \
		-e DOCKER_BUILD_NETWORK \
		-e UBUNTU_MIRROR \
		$(DAPPER_IMAGE) $@

.DEFAULT_GOAL := ci

.PHONY: $(TARGETS)

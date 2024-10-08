.PHONY: all build build-in-docker build-docker clean clean-docker test

# Default target, which test and builds the web service binaries
all: test build

# This target builds the web service binaries on the host machine
build:
ifeq ($(GOOS),)
	$(error GOOS is not set)
endif
ifeq ($(GOARCH),)
	$(error GOARCH is not set)
endif
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o bin/risk-api-$(GOOS)-$(GOARCH) -ldflags="-w -s" ./cmd/risk-api
	@if [ "$(GOOS)" = "windows" ]; then mv bin/risk-api-$(GOOS)-$(GOARCH) bin/risk-api-$(GOOS)-$(GOARCH).exe; fi 

# This target builds the web service binaries in a Docker container
build-in-docker:
ifeq ($(GOOS),)
	$(error GOOS is not set)
endif
ifeq ($(GOARCH),)
	$(error GOARCH is not set)
endif
	@docker build \
		-f build/package/Dockerfile \
		--build-arg GOOS=$(GOOS) \
		--build-arg GOARCH=$(GOARCH) \
		--target build \
		-t risk-api-builder \
		.

	@docker create --name risk-api-builder-tmp risk-api-builder
	@if [ "$(GOOS)" = "windows" ]; then \
		docker cp risk-api-builder-tmp:/app/bin/risk-api-$(GOOS)-$(GOARCH).exe bin/risk-api-$(GOOS)-$(GOARCH).exe; \
	else \
		docker cp risk-api-builder-tmp:/app/bin/risk-api-$(GOOS)-$(GOARCH) bin/risk-api-$(GOOS)-$(GOARCH); \
	fi
	@docker rm risk-api-builder-tmp

# This target builds the Docker image for the web service
build-docker:
	@docker build \
		--progress=plain \
		-f build/package/Dockerfile \
		--target run \
		-t risk-api \
		.

# This target removes the web service binaries from the host machine
clean:
	@rm bin/risk-api-*

# This target removes the images from the local Docker registry
clean-docker:
	@docker rmi -f risk-api-builder risk-api

# This target runs the unit tests for the web service
test:
	go test ./...

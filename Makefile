BINARIES=bin/stackctl
IMAGE=containers.trusch.io/stackctl:latest
BASE_IMAGE=gcr.io/distroless/base:latest
BUILD_IMAGE=golang:1.14

GOOS=linux
GOARCH=amd64
GOARM=7 # for crosscompiling

COMMIT=$(shell git log --format="%H" -n 1)
VERSION=$(shell git describe)

default: image

# rebuild and run the server
run: image
	podman run --rm -p 3001:3001 -p 8080:8080 $(IMAGE) /bin/stackctl serve

# put binaries into image
image: .image
.image: $(BINARIES) Makefile
	$(eval ID=$(shell buildah from $(BASE_IMAGE)))
	buildah copy $(ID) $(shell pwd)/bin/* /bin/
	buildah commit $(ID) $(IMAGE)
	buildah rm $(ID)
	touch .image

# build binaries
bin/%: $(shell find ./ -name "*.go")
	podman run \
		--rm \
		-v ./:/app \
		-w /app \
		-e GOOS=${GOOS} \
		-e GOARCH=${GOARCH} \
		-e GOARM=${GOARM} \
		-v go-build-cache:/root/.cache/go-build \
		-v go-mod-cache:/go/pkg/mod $(BUILD_IMAGE) \
			go build -v -o $@ -ldflags "-X github.com/trusch/stackctl/cmd/stackctl/cmd.Version=$(VERSION) -X github.com/trusch/stackctl/cmd/stackctl/cmd.Commit=$(COMMIT)" cmd/stackctl/main.go

# install locally
install: ${BINARIES}
	cp -v ${BINARIES} $(shell go env GOPATH)/bin/

crosscompile:
	$(MAKE) bin/stackctl-darwin-amd64 GOOS=darwin GOARCH=amd64
	$(MAKE) bin/stackctl-linux-amd64 GOOS=linux GOARCH=amd64
	$(MAKE) bin/stackctl-linux-arm64 GOOS=linux GOARCH=arm64
	$(MAKE) bin/stackctl-linux-arm GOOS=linux GOARM=7 GOARCH=arm

# cleanup
clean:
	-rm -r bin .image .buildimage /tmp/protoc-download
	-podman volume rm  go-build-cache go-mod-cache

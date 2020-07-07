BINARIES=bin/stackctl
IMAGE=containers.trusch.io/examples/stackctl:latest
BASE_IMAGE=gcr.io/distroless/base:latest
BUILD_IMAGE=golang:1.14

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
	buildah copy $(ID) ./bin/* /bin/
	buildah commit $(ID) $(IMAGE)
	buildah rm $(ID)
	touch .image

# build binaries
bin/%: $(shell find ./ -name "*.go")
	podman run \
		--rm \
		-v ./:/app \
		-w /app \
		-v go-build-cache:/root/.cache/go-build \
		-v go-mod-cache:/go/pkg/mod $(BUILD_IMAGE) \
			go build -v -o $@ -ldflags "-X github.com/trusch/stackctl/cmd/stackctl/cmd.Version=$(VERSION) -X github.com/trusch/stackctl/cmd/stackctl/cmd.Commit=$(COMMIT)" cmd/$(shell basename $@)/main.go

# install locally
install: ${BINARIES}
	cp -v ${BINARIES} $(shell go env GOPATH)/bin/

# cleanup
clean:
	-rm -r bin .image .buildimage /tmp/protoc-download
	-podman volume rm  go-build-cache go-mod-cache

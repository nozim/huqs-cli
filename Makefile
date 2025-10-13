.PHONY: build

OUTPUT_NAME = huqs-cli-${GOOS}-${GOARCH}
GIT_SHORT_REV=${GIT_SHORT_REV}

ifeq ($(GOOS),windows)
    OUTPUT_NAME := $(OUTPUT_NAME).exe
endif

clean:
	rm -rf build

build: clean
	go build -o build/huqs-cli main.go
install:
	mv ./build/huqs ${GOBIN}

release:
	@if [ -z "${GITHUB_TOKEN}" ]; then \
			echo "Error: GITHUB_TOKEN environment variable is not set"; \
			exit 1; \
	fi
	./scripts/release.sh

binary: clean
	CGO_ENABLED=0 GOGC=off GOOS=${GOOS} GOARCH=${GOARCH} \
	go build -installsuffix nocgo -o "./build/${OUTPUT_NAME}" main.go
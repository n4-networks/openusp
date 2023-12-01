TOP_DIR := $(shell pwd)
SUBDIRS := $(TOP_DIR)/cmd/apiserver
SUBDIRS += $(TOP_DIR)/cmd/cli
SUBDIRS += $(TOP_DIR)/cmd/controller


.PHONY: build install clean
build installall clean:
	for dir in $(SUBDIRS); do \
	  $(MAKE) -C $$dir -f Makefile $@; \
	done

.PHONY: controller
controller:
	  $(MAKE) -C cmd/controller -f Makefile build

.PHONY: apiserver
apiserver:
	  $(MAKE) -C cmd/apiserver -f Makefile build

.PHONY: cli
cli:
	  $(MAKE) -C cmd/cli -f Makefile build

.PHONY: images
images:
	docker buildx build -t n4networks/openusp-controller:latest -f build/controller/Dockerfile --push --platform=linux/amd64,linux/arm64 .
	docker buildx build -t n4networks/openusp-apiserver:latest -f build/apiserver/Dockerfile --push --platform=linux/amd64,linux/arm64 .
	docker buildx build -t n4networks/openusp-cli:latest -f build/cli/Dockerfile --push --platform=linux/amd64,linux/arm64 .


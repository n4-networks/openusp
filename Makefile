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

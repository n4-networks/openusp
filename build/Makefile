TOP_DIR := $(shell pwd)/..
SUBDIRS := $(TOP_DIR)/cmd/apiserver
SUBDIRS += $(TOP_DIR)/cmd/cli
SUBDIRS += $(TOP_DIR)/cmd/controller
SUBDIRSCLEAN=$(SUBDIRS)


install clean:
	for dir in $(SUBDIRS); do \
	  $(MAKE) -C $$dir -f Makefile $@; \
  done

#target: $(SUBDIRS)

#$(SUBDIRS):
#	$(MAKE) -C $@

#.PHONY: target $(SUBDIRS) clean $(SUBDIRSCLEAN)

#clean: $(SUBDIRSCLEAN)
#$(SUBDIRSCLEAN):
#	$(MAKE) -C $@ clean

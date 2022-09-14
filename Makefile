# kernel-style V=1 build verbosity
ifeq ("$(origin V)", "command line")
	BUILD_VERBOSE = $(V)
endif

ifeq ($(BUILD_VERBOSE),1)
	Q =
else
	Q = @
endif

VERSION=$(shell git describe --dirty)
REPO=github.com/flatcar/locksmith
LD_FLAGS="-w -s"

export GOPATH=$(shell pwd)/gopath

.PHONY: all
all: bin/locksmithctl

bin/%:
	$(Q)go build -o $@ -ldflags $(LD_FLAGS) $(REPO)/$*

.PHONY: test
test:
	$(Q)./scripts/test

.PHONY: vendor
vendor:
	$(Q)go mod vendor

.PHONY: clean
clean:
	$(Q)rm -rf bin

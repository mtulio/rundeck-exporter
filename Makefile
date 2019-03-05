include Makefile-glob.mk
include Makefile-fc.mk

# Defaut APP name to build
APP_NAME ?= rundeck-exporter

# Build a version
.PHONY: build
build: clean
	@test -d $(PWD)/bin || mkdir $(PWD)/bin
	@$(foreach dircmd,$(shell ls cmd/), \
	cd cmd/$(dircmd); \
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build \
		-ldflags "$(LDFLAGS)" \
		$(BUILD_TAGS) \
		-o $(BIN_NAME)/$(APP_NAME) && strip $(BIN_NAME)/$(APP_NAME) \
	; cd -)

.PHONY: run
run:
	$(BIN_NAME)

.PHONY: version
version: build
	$(BIN_NAME) -v

.PHONY: clean
clean:
	@rm -f bin/$(BIN_NAME)

# ##################
# Release Management
#

# Release
tag:
	$(call deps_tag,$@)
	git tag -a $(shell cat VERSION) -m "$(message)"
	git push origin $(shell cat VERSION)


# Goreleaser
# https://goreleaser.com/introduction/
GORELEASE_ROOT ?= ../../
gorelease-init:
	goreleaser init

release:
	@$(foreach dircmd,$(shell ls cmd/), \
	cd cmd/$(dircmd); \
	. $(GORELEASE_ROOT)/hack/env-build.sh && \
		goreleaser --rm-dist -f $(GORELEASE_ROOT)/.goreleaser.yml \
	; cd -)

release-snap:
	@$(foreach dircmd,$(shell ls cmd/), \
	cd cmd/$(dircmd); \
	. $(GORELEASE_ROOT)/hack/env-build.sh && \
		goreleaser --rm-dist --snapshot -f $(GORELEASE_ROOT)/.goreleaser.yml \
	; cd -)

# DEV Release builder
# Using ghr to avoid dependencies in goreleaser
dev-release: build
	ghr $(RELEASE_VERSION) bin/

dev-release-master: build
	ghr --recreate $(RELEASE_VERSION) bin/

# ################
# Dev Dependencies

# Goreleaser
# Go project release management
# https://goreleaser.com/
dep-install-goreleaser:
	ggo get github.com/goreleaser/goreleaser

# GHR
# GitHub releaser
# https://github.com/tcnksm/ghr
dep-install-ghr:
	go get -u github.com/tcnksm/ghr


dep-all: dep-install-goreleaser dep-install-ghr
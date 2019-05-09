
#  Makefile for Go
#
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_BUILD_RACE=$(GO_CMD) build -race
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_INSTALL=$(GO_CMD) install -v
GO_CLEAN=$(GO_CMD) clean
GO_DEPS=$(GO_CMD) get -d -v
GO_DEPS_UPDATE=$(GO_CMD) get -d -v -u
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
GO_LINT=golint
DEP_COMMAND=dep
DEP_INSTALL=$(DEP_COMMAND) ensure -v

# Packages
PACKAGE_LIST := main.go struct.go utils.go endpoints.go 

.PHONY: all build install clean fmt vet lint list

all: build

build: vet
		echo "==> Build package ...";
		$(GO_BUILD)  -o bin/webex-teams-grafana-alerts-webhook.osx main.go struct.go utils.go endpoints.go  || exit 1;
		GOOS=linux GOARCH=amd64 $(GO_BUILD)  -o bin/webex-teams-grafana-alerts-webhook.linux main.go struct.go utils.go endpoints.go  || exit 1;

build-race: vet

install:
	$(DEP_INSTALL);

clean:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Clean $$p ..."; \
		$(GO_CLEAN) $$p; \
	done

fmt:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Formatting $$p ..."; \
		$(GO_FMT) $$p || exit 1; \
	done

vet:
	$(GO_VET) $(PACKAGE_LIST);

lint:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Lint $$p ..."; \
		$(GO_LINT) src/$$p; \
	done

list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

# vim: set noexpandtab shiftwidth=8 softtabstop=0:

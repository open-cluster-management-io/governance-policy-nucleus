PWD := $(shell pwd)

## Location to install dependencies to
LOCAL_BIN ?= $(PWD)/bin
$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

# Keep an existing GOPATH, make a private one if it is undefined
GOPATH_DEFAULT := $(PWD)/.go
export GOPATH ?= $(GOPATH_DEFAULT)
GOBIN_DEFAULT := $(GOPATH)/bin
export GOBIN ?= $(GOBIN_DEFAULT)

# Set PATH so that locally installed things will be used first
export PATH=$(LOCAL_BIN):$(GOBIN):$(shell echo $$PATH)

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# go-install will 'go install' any package $1 to LOCAL_BIN
# Note: this replaces `go-get-tool`.
go-install = @set -e ; mkdir -p $(LOCAL_BIN) ; GOBIN=$(LOCAL_BIN) go install $(1)

# Define local utilities near the top so they work correctly as targets
# Note: this pattern of variables, paths, and target names allows users to
#  override the version used, and helps Make by not using PHONY targets.

CONTROLLER_GEN ?= $(LOCAL_BIN)/controller-gen
$(CONTROLLER_GEN): $(LOCAL_BIN)
	$(call go-install,sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0)

ENVTEST ?= $(LOCAL_BIN)/setup-envtest
$(ENVTEST): $(LOCAL_BIN)
	$(call go-install,sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

KUSTOMIZE ?= $(LOCAL_BIN)/kustomize
$(KUSTOMIZE): $(LOCAL_BIN)
	$(call go-install,sigs.k8s.io/kustomize/kustomize/v4@v4.5.5)

GOLANGCI ?= $(LOCAL_BIN)/golangci-lint
$(GOLANGCI): $(LOCAL_BIN)
	$(call go-install,github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2)

.PHONY: manifests
manifests: $(CONTROLLER_GEN) ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: $(CONTROLLER_GEN) ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

# Note: this target is not used by Github Actions. Instead, each linter is run 
# separately to automatically decorate the code with the linting errors.
# Note: this target will fail if yamllint is not installed.
.PHONY: lint
lint: $(GOLANGCI)
	$(GOLANGCI) run
	yamllint .

# ENVTEST_K8S_VERSION = 1.23
.PHONY: test
test: manifests generate $(ENVTEST) ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use 1.23 -p path)" go test ./... -coverprofile cover.out

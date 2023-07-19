LOCAL_BIN ?= $(PWD)/../../bin
GOPATH_DEFAULT ?= $(PWD)/../../.go

include ../../Makefile

.PHONY: manifests
manifests: $(CONTROLLER_GEN) ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

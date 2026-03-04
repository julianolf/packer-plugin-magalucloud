ifneq (,$(wildcard .env))
include .env
export
endif

NAME=magalucloud
BINARY=packer-plugin-$(NAME)

COUNT?=1
TEST?=$(shell go list ./...)

PACKER_SDK=github.com/hashicorp/packer-plugin-sdk
PACKER_SDK_VERSION?=$(shell go list -m $(PACKER_SDK) | cut -d " " -f2)
PACKER_SDC=$(PACKER_SDK)/cmd/packer-sdc@$(PACKER_SDK_VERSION)
PLUGIN_FQN=$(shell grep -E '^module' <go.mod | sed -E 's/module \s*//')

IMAGE_URL?=
IMAGE_FILE=$(shell pwd)/post-processor/magalucloud/test-fixtures/image.qcow2

.PHONY: dev docs env

build:
	@go build -o $(BINARY)

gen:
	@go generate ./...

dev: gen
	go build -ldflags="-X '$(PLUGIN_FQN)/version.VersionPrerelease=dev'" -o $(BINARY)
	packer plugins install --path $(BINARY) "$(shell echo "$(PLUGIN_FQN)" | sed 's/packer-plugin-//')"

test:
	@go test -race -count $(COUNT) $(TEST) -timeout=5m

$(IMAGE_FILE):
	@curl -sL $(IMAGE_URL) -o $@

tools:
	@go install $(PACKER_SDC)

testacc: $(IMAGE_FILE) tools dev
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m

plugin-check: tools build
	@packer-sdc plugin-check $(BINARY)

.env:
	@cp .env.example $@

env: .env tools

docs: tools
	@rm -rf .docs/
	@packer-sdc renderdocs -src docs -partials docs-partials/ -dst .docs/
	@./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs" "julianolf"

clean:
	@rm -rf .docs/ packer-plugin-magalucloud .coverage* coverage.*

NAME:= amazon-bedrock-metadata
BINARY := bootstrap
BUCKET_NAME ?= greenscale-ai-public
LOCAL_ARTIFACTS_PATH := artifacts/s3
SHA := $(shell git rev-parse HEAD | cut -c 1-7)
VERSION := 1.0
BUILD := $(VERSION)-$(SHA)
GOARCH ?= arm64
LDFLAGS="-X main.Version=$(BUILD) -s -w"

all: clean fmt-check vet-check build test                    ## Build the code
.PHONY: all

clean:                                                       ## Remove installed packages/temporary files
	@go clean ./...
	@rm -rf $(LOCAL_ARTIFACT_PATH) bootstrap
	@rm -rf $(LOCAL_ARTIFACTS_PATH)
.PHONY: clean

build:                                                       ## Build the amazon-bedrock-metadata binary
	@GOOS=linux GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o bootstrap github.com/greenscale-ai/amazon-bedrock-metadata/cmd
.PHONY: build

artifacts:                                                   ## Create artifacts to be uploaded to S3
	@rm -rf $(LOCAL_ARTIFACTS_PATH)
	@mkdir -p $(LOCAL_ARTIFACTS_PATH)
	@zip $(LOCAL_ARTIFACTS_PATH)/$(NAME)-$(BUILD).zip bootstrap models.json
	@cp -f cloudformation/template.json $(LOCAL_ARTIFACTS_PATH)/template_build_$(BUILD).json
	@sed -e "s#1.0#$(BUILD)#" $(LOCAL_ARTIFACTS_PATH)/template_build_$(BUILD).json > $(LOCAL_ARTIFACTS_PATH)/template_build_$(BUILD).json.new
	@mv -- $(LOCAL_ARTIFACTS_PATH)/template_build_$(BUILD).json.new $(LOCAL_ARTIFACTS_PATH)/template_build_$(BUILD).json
	@mv -f $(LOCAL_ARTIFACTS_PATH)/template_build_$(BUILD).json $(LOCAL_ARTIFACTS_PATH)/template.json
.PHONY: artifacts

upload: artifacts                                            ## Upload artifacts to S3
	aws s3 sync artifacts/s3/ s3://$(BUCKET_NAME)/amazon-bedrock-metadata
.PHONY: upload

vet-check:                                                   ## Verify vet compliance
	@go vet -all ./...
.PHONY: vet-check

fmt-check:                                                   ## Verify fmt compliance
	@sh -c 'test -z "$$(gofmt -l -s -d . | tee /dev/stderr)"'
.PHONY: fmt-check

test:                                                        ## Test go code and coverage
	@go test -covermode=count ./...
.PHONY: test

help:                                                        ## Show this help
	@printf "Rules:\n"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
.PHONY: help

DOCKER := docker

BASE_IMAGE = centos:centos8

REGISTRY_PREFIX ?= igame-sg.tencentcloudcr.com/igame/view_proto

## iamge.build : Build iamge
.PHONY : image.build.%
image.build.%: go.build.%
	$(eval IMAGE := $(COMMAND))
	@echo "===========> Building docker image $(IMAGE)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)
	@cat $(ROOT_DIR)/build/docker/$(IMAGE)/Dockerfile\
		| sed "s#BASE_IMAGE#$(BASE_IMAGE)#g" > $(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp $(OUTPUT_DIR)/$(IMAGE) $(TMP_DIR)/$(IMAGE)
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION) $(TMP_DIR)/$(IMAGE))
	@echo "$(DOCKER) build $(BUILD_SUFFIX)"
	@$(DOCKER) build $(BUILD_SUFFIX) ;
	
.PHONY: image.push.%
image.push.%: image.build.%
	@echo "===========> Pushing image $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) push $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION)
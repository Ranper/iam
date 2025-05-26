.PHONY : build

GO := go

BUILD_DIR = /opt/iam

ifneq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"


# $* 表示当前规则的目标模式中 % 匹配的部分
# $(subst ., ,$*) 将 $* 中的点号 . 替换为空格
# $(word 2,...) 取替换后的字符串中的第二个单词
.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 1,$(subst ., ,$*)))
	@echo "Building go binary"
	@echo "$(ROOT_DIR)/cmd/$(COMMAND)"
	@$(GO) build $(GO_BUILD_FLAGS)  -o $(OUTPUT_DIR)/$(COMMAND) $(ROOT_DIR)/cmd/$(COMMAND)
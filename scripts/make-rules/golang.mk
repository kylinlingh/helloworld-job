# Makefile helper functions for golang
#

GO := go
GO_SUPPORTED_VERSIONS ?= 1.13|1.14|1.15|1.16|1.17|1.18|1.19|1.20|1.21
GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
ifneq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

# 如果当前的编译系统是windows
ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
    $(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

#假设cmd目录结构如下：
#cmd
#├── apiserver
#│   └── apiserver.go
#└── module2
#    └── module2.go

# wildcard：列出目录下的文件
# COMMANDS：/Users/kylin/Code/go-template/go-template-v1/cmd/apiserver /Users/kylin/Code/go-template/go-template-v1/cmd/module2
COMMANDS ?= $(filter-out %.md, $(wildcard ${ROOT_DIR}/cmd/*))
# 取出绝对路径的最后一个文件名
# BINS: apiserver module2
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))

ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif


.PHONY: go.build.verify
go.build.verify:
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif
# 验证go版本是否在支持的范围内

# % 的意思是匹配零或若干字符
.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/$(COMMAND)

.PHONY: go.build
go.build: go.build.verify $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))
#	@echo $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))
# 解析：
# 这里的真实规则是：go.build: go.build.verify go.build.darwin_arm64.apiserver go.build.darwin_arm64.module2
# 所以这里依赖了三个规则，而go.build.darwin_arm64.apiserver会匹配并执行规则go.build.%里的构建命令

.PHONY: go.build.multiarch
go.build.multiarch: go.build.verify $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))
# 解析：用于编译其他平台的可执行文件，需要在命令参数里传入PLATFORMS参数，参考使用说明：make help


.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

# 不需要执行单元测试的目录
EXCLUDE_TESTS=go-template-v1/test go-template-v1/pkg/log go-template-v1/third_party go-template-v1/internal/pump/storage go-template-v1/internal/pump go-template-v1/internal/pkg/logger

# set -o pipefail 命令的作用是当管道（|）中的任何一个命令返回非零退出状态码时，整个管道命令会立即以非零状态码退出，这可以帮助检测管道中的错误，防止错误被忽略而导致问题。
# 执行go test时设置了超时时间、竞态检查，开启了代码覆盖率检查，覆盖率测试数据保存在了coverage.out文件中
# go-junit-report将 go test 的结果转化成了 xml 格式的报告文件，该报告文件会被一些 CI 系统，例如 Jenkins 拿来解析并展示结果
.PHONY: go.test
go.test: tools.verify.go-junit-report
	@echo "===========> Run unit test"
	#@set -o pipefail;
	$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short -v `go list ./...|\
		egrep -v $(subst $(SPACE),'|',$(sort $(EXCLUDE_TESTS)))` 2>&1 | \
		tee >(go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml)
	# Mock 的代码是不需要编写测试用例的，需要将 Mock 代码的单元测试覆盖率数据从coverage.out文件中删除掉，会执行失败，不知道原因
	#@sed -i '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

.PHONY: go.test.cover
go.test.cover: go.test
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
		awk -v target=$(COVERAGE) -f $(ROOT_DIR)/scripts/coverage.awk

.PHONY: go.updates
go.updates: tools.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct

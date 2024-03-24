# Copyright 2023 Kylin Lin<kylinlingh@foxmail.com>. All rights reserved.

SHELL := /bin/bash

# 特殊变量MAKEFILE_LIST：make命令所需要处理的makefile文件列表，当前makefile的文件名总是位于列表的最后，文件名之间以空格进行分隔
# COMMON_SELF_DIR = scripts/make-rules/
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

# 如果ROOT_DIR这个变量从来没有被定义过
ifeq ($(origin ROOT_DIR), undefined)
	# ROOT_DIR是工程的根目录
	ROOT_DIR :=  $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif
ifeq ($(origin OUTPUT_DIR),undefined)
	OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif
ifeq ($(origin TOOLS_DIR),undefined)
	TOOLS_DIR := $(OUTPUT_DIR)/tools
$(shell mkdir -p $(TOOLS_DIR))
endif
ifeq ($(origin TMP_DIR),undefined)
	TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

# 设定版本号
ifeq ($(origin VERSION), undefined)
	# 如果当前版本已经有tag则直接输出此tag名，没有的话就输出上一个commit的id
	VERSION := $(shell git describe --tags --always --match='v*')
endif
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
# 显示HEAD提交的SHA1值
GIT_COMMIT:=$(shell git rev-parse HEAD)

# 测试覆盖率阈值为60%
ifeq ($(origin COVERAGE),undefined)
COVERAGE := 10
endif

# 只能用linux的系统来构建docker镜像
# PLATFORMS 只有在没设置值时才会等于 linux_amd64 linux_arm64
PLATFORMS ?= linux_amd64 linux_arm64

ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
	# GOOS的默认值是我们当前的操作系统，注意mac os操作的上的值是darwin
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
	# GOARCH则表示CPU架构，譬如：arm64
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# IMAGE_PLAT = linux_arm64
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# Linux command settings
FIND := find . ! -path './third_party/*' ! -path './vendor/*'
XARGS := xargs -r

# Makefile settings
ifndef V
MAKEFLAGS += --no-print-directory
endif

# githooks用于检查提交的commit信息是否规范
#COPY_GITHOOK:=$(shell cp -f githooks/* .git/hooks/)

# Specify components which need certificate
ifeq ($(origin CERTIFICATES),undefined)
CERTIFICATES=iam-apiserver iam-authz-server admin
endif

# Specify tools severity, include: BLOCKER_TOOLS, CRITICAL_TOOLS, TRIVIAL_TOOLS.
# Missing BLOCKER_TOOLS can cause the CI flow execution failed, i.e. `make all` failed.
# Missing CRITICAL_TOOLS can lead to some necessary operations failed. i.e. `make release` failed.
# TRIVIAL_TOOLS are Optional tools, missing these tool have no affect.
BLOCKER_TOOLS ?= gsemver golines go-junit-report golangci-lint addlicense goimports codegen
CRITICAL_TOOLS ?= swagger mockgen gotests git-chglog github-release coscmd go-mod-outdated protoc-gen-go cfssl go-gitlint
TRIVIAL_TOOLS ?= depth go-callvis gothanks richgo rts kube-score

COMMA := ,
SPACE :=
SPACE +=


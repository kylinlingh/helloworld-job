# $(MAKE) 在Makefile 中表示make 程序的名称

.PHONY: all
all: gen

ROOT_PACKAGE=helloworld-job

include scripts/make-rules/common.mk # 此文件必须在第一行里include，include只是做单纯的文本替换
include scripts/make-rules/golang.mk
include scripts/make-rules/gen.mk
include scripts/make-rules/tools.mk # 安装要用到的工具


## gen: Generate all necessary files, such as error code files.
.PHONY: gen
gen:
	@$(MAKE) gen.run
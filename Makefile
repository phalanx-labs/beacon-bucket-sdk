# 变量定义，方便后续维护
MAIN_FILE = main.go
SWAG_CMD = swag
SWAG_FLAGS = --parseDependency
BUILD_SCRIPT = script/build-docker.sh
SCRIPT_DIR = script
PROTO_FILE ?= proto/normal_upload.proto
BASE_GO_MODULE_DIR := $(shell go list -m -f '{{.Dir}}' github.com/bamboo-services/bamboo-base-go)
XBASE_LINK := proto/link/base.proto

.DEFAULT_GOAL := help

.PHONY: help proto tidy proto-init connect-install tag tag-upload release

# 显示帮助信息
help:
	@echo "BeaconBucketSDK - 可用命令"
	@echo ""
	@echo "开发命令:"
	@echo "  make connect-install - 安装 Connect-Go 代码生成器"
	@echo "  make proto           - 生成指定 proto 的客户端代码"
	@echo "                        示例: make proto PROTO_FILE=proto/normal_upload.proto"
	@echo "  make proto-init      - 初始化 proto 符号链接"
	@echo ""
	@echo "发布命令:"
	@echo "  make tag        	- 创建带有时间戳的 tag（不推送）"
	@echo "                   	  格式: v{version}-{YYYYMMDDHHMM}"
	@echo "                   	  示例: v1.0.0-202602191755"
	@echo "  make tag-upload 	- 单独上传 tag"
	@echo "  make release    	- 创建 tag 并推送到远程仓库"
	@echo ""

# 安装 Connect-Go 代码生成器
connect-install:
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	@echo "✅ Connect-Go 代码生成器安装完成"

# 初始化 proto 符号链接
proto-init:
	@mkdir -p $(dir $(XBASE_LINK))
	@ln -sf $(BASE_GO_MODULE_DIR)/proto/base.proto $(XBASE_LINK)
	@echo "符号链接已创建: $(XBASE_LINK) -> $(BASE_GO_MODULE_DIR)/proto/base.proto"

# 生成 proto（自动初始化符号链接）
proto: proto-init
	buf generate --path $(PROTO_FILE)

tidy:
	go mod tidy

# 创建 tag（仅本地）
tag:
	@echo "创建 tag: $(TAG_NAME)"
	git tag -a $(TAG_NAME) -m "Release $(TAG_NAME)"
	@echo "✅ Tag $(TAG_NAME) 创建成功"

tag-upload:
	@echo "推送 tag 到远程仓库..."
	git push --tags
	@echo "✅ Tag $(TAG_NAME) 推送成功！"

# 创建 tag 并推送
release: tag tag-upload

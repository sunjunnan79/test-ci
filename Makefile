# 默认镜像标签
TAG ?= latest

# 构建目标
build:
	@echo "构建服务：镜像标签：$(TAG)"
	docker buildx build \
		--platform linux/amd64 \
		-t crpi-r1jx5ale23646u4w.cn-hongkong.personal.cr.aliyuncs.com/clairvoyance-project/test-ci:$(TAG) \
		--build-arg GOOS=linux \
		--build-arg GOARCH=amd64 \
		--no-cache \
		.
# 指定 linux/amd64 平台
# 使用指定标签
# 构建 Linux 平台的二进制文件
# 构建 64 位架构的二进制文件
# 禁用缓存，确保每次都重新构建
# 确保这里是构建上下文的路径

# 推送镜像目标
push:
	@echo "推送镜像：crpi-r1jx5ale23646u4w.cn-hongkong.personal.cr.aliyuncs.com/clairvoyance-project/test-ci:$(TAG)"
	docker push crpi-r1jx5ale23646u4w.cn-hongkong.personal.cr.aliyuncs.com/clairvoyance-project/test-ci:$(TAG)

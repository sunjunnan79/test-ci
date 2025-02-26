# 默认镜像标签
TAG ?= latest

# 构建目标
build:
	@echo "构建服务：镜像标签：$(TAG)"
	docker buildx build \
		--platform linux/amd64 \
		-t sunjunnan112/test-ci:$(TAG) \
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
	@echo "推送镜像：sunjunnan112/test-ci:$(TAG)"
	docker push sunjunnan112/test-ci:$(TAG)

pushLatest:
	@echo "推送$(TAG)版本镜像为latest"
	docker tag sunjunnan112/test-ci:$(TAG) sunjunnan112/test-ci:latest
	docker push sunjunnan112/test-ci:latest

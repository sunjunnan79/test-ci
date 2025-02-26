# 第一阶段：构建 Go 应用
FROM golang:1.23-alpine AS builder

# 设置 Go 代理为七牛云的代理
ENV GOPROXY=https://goproxy.cn,direct

# 切换到 /app 目录
WORKDIR /app

# 拷贝基础文件和项目目录
COPY . /app

# 构建应用
RUN go mod tidy && go build -o app

# 第二阶段：构建最终镜像
FROM alpine:3.17

# 安装 tzdata 来设置时区
RUN apk add --no-cache tzdata

# 设置时区为 Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录为 /app
WORKDIR /app
EXPOSE 8080

# 启动用户服务
CMD ["./app"]

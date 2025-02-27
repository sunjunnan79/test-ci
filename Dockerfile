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

# 设置工作目录为 /app
WORKDIR /app

# 将构建的应用从 builder 镜像复制到最终镜像
COPY --from=builder /app/app .

# 为可执行文件增加最高权限
RUN chmod 777 ./app

# 安装 tzdata 来设置时区
RUN apk add --no-cache tzdata

# 设置时区为 Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

EXPOSE 8080

# 启动用户服务
CMD ["./app"]

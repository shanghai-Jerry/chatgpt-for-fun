# 使用Go 1.21官方镜像作为构建环境
FROM golang:1.23-alpine AS builder

# 安装git（用于go mod下载依赖）
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod tidy

# 复制源代码
COPY . .

# 构建二进制文件
RUN GOOS=linux go build  -o main .

# 使用scratch镜像作为运行环境，减小镜像体积
FROM alpine:latest

# 安装ca-certificates，用于HTTPS请求
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从builder阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行二进制文件
CMD ["./main"]
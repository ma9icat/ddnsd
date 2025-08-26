FROM golang:1.25-alpine AS builder

LABEL org.opencontainers.image.description "DDNSD (Dynamic DNS Daemon) is a lightweight service that automatically updates DNS records for domains with dynamic IP addresses."

WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ddnsd .

# 使用scratch镜像作为运行环境
FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates

WORKDIR /app/

# 从builder中复制预构建的二进制文件
COPY --from=builder /app/ddnsd .

# 复制环境变量文件
# COPY .env /app/.env

# 运行二进制文件
CMD ["/app/ddnsd"]
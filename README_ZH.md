# DDNSD - 动态DNS服务

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

DDNSD (Dynamic DNS Daemon) 是一个轻量级服务，可自动更新动态IP地址的DNS记录。它支持多种DNS提供商，可以持续运行以监控和更新您的IP地址变化。

## 功能特性

- 支持多种DNS提供商：
  - 腾讯云DNSPod
  - Cloudflare
  - 阿里云(Alibaba Cloud)
- 支持IPv4和IPv6
- 自动检测IP地址
- 可配置的更新间隔
- 支持Docker部署
- 轻量级且高效

## 支持的DNS提供商

| 提供商         | 中国大陆 | 国际版 | 环境变量值              |
|----------------|----------|--------|-------------------------|
| DNSPod         | ✅       | ✅     | `dnspod`                |
| Cloudflare     | ✅       | ✅     | `cloudflare`            |
| 阿里云         | ✅       | ✅     | `aliyun`, `alibabacloud` |

## 快速开始

### 前置条件

- Go 1.25 或更高版本
- Docker (可选，用于容器化部署)
- 所选DNS提供商的API凭证

### 配置

根据 [.env.simple](file:///home/magicat/projects/ddnsd/.env.simple) 创建 `.env` 文件：

```bash
# DNS提供商信息
# 可选值: dnspod, cloudflare, aliyun, alibabacloud
DNS_PROVIDER=dnspod

# 您的DNS提供商凭证
SECRET_ID=your_secret_id
SECRET_KEY=your_secret_key

# IPv4/IPv6开关
IPV4_ENABLED=true
IPV6_ENABLED=true

# 域名信息
IPV4_DOMAIN=example.com
IPV4_SUBDOMAINS=ipv4,ipv4.www
IPV6_DOMAIN=example.com
IPV6_SUBDOMAINS=ipv6,ipv6.www

# 更新间隔时间(秒)，默认300秒(5分钟)
INTERVAL=300
```

### 使用Go运行

1. 克隆代码库：
   ```bash
   git clone https://github.com/yourusername/ddnsd.git
   cd ddnsd
   ```

2. 创建您的 `.env` 配置文件

3. 运行服务：
   ```bash
   go run main.go
   ```

### 使用Docker运行

1. 构建Docker镜像：
   ```bash
   docker build -t ddnsd .
   ```

2. 使用您的配置运行容器：
   ```bash
   docker run -d --env-file .env ddnsd
   ```

### 构建

构建项目：

```bash
go build -o ddnsd .
```

使用Docker构建：

```bash
./build.sh
```

## 配置详情

### DNSPod配置

对于腾讯云DNSPod，您需要从[腾讯云控制台](https://console.cloud.tencent.com/cam/capi)获取API凭证。

```env
DNS_PROVIDER=dnspod
SECRET_ID=your_tencent_cloud_secret_id
SECRET_KEY=your_tencent_cloud_secret_key
```

### Cloudflare配置

对于Cloudflare，您需要从[Cloudflare仪表板](https://dash.cloudflare.com/profile/api-tokens)获取邮箱和API密钥。

```env
DNS_PROVIDER=cloudflare
SECRET_ID=your_email@example.com
SECRET_KEY=your_cloudflare_api_key
```

### 阿里云配置

对于阿里云，您需要从[阿里云控制台](https://ram.console.aliyun.com/manage/ak)获取AccessKey ID和AccessKey Secret。

```env
# 中国大陆区域
DNS_PROVIDER=aliyun
SECRET_ID=your_access_key_id
SECRET_KEY=your_access_key_secret

# 国际区域
DNS_PROVIDER=alibabacloud
SECRET_ID=your_access_key_id
SECRET_KEY=your_access_key_secret
```

## 环境变量

| 变量名              | 描述                           | 默认值                                |
|---------------------|--------------------------------|---------------------------------------|
| DNS_PROVIDER        | 使用的DNS提供商                | `dnspod`                              |
| SECRET_ID           | 提供商特定的凭证               | (必填)                                |
| SECRET_KEY          | 提供商特定的凭证               | (必填)                                |
| IPV4_ENABLED        | 启用IPv4更新                   | `true`                                |
| IPV6_ENABLED        | 启用IPv6更新                   | `false`                               |
| IPV4_DOMAIN         | IPv4记录的主域名               | (IPv4启用时必填)                      |
| IPV4_SUBDOMAINS     | 逗号分隔的IPv4子域名           | (IPv4启用时必填)                      |
| IPV6_DOMAIN         | IPv6记录的主域名               | (IPv6启用时必填)                      |
| IPV6_SUBDOMAINS     | 逗号分隔的IPv6子域名           | (IPv6启用时必填)                      |
| INTERVAL            | 更新间隔（秒）                 | `300` (5分钟)                         |
| IPV4_CHECK_URL      | 检查IPv4地址的服务             | `https://iplark.com/ipapi/public/ip`  |
| IPV6_CHECK_URL      | 检查IPv6地址的服务             | `https://6.iplark.com/ip`             |

## 许可证

该项目基于MIT许可证 - 详见 [LICENSE](LICENSE) 文件。

## 致谢

- [腾讯云SDK for Go](https://github.com/TencentCloud/tencentcloud-sdk-go)
- [Godotenv](https://github.com/joho/godotenv)
- [Cron](https://github.com/robfig/cron)
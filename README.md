# DDNSD - Dynamic DNS Service

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

DDNSD (Dynamic DNS Daemon) is a lightweight service that automatically updates DNS records for domains with dynamic IP addresses. It supports multiple DNS providers and can run continuously to monitor and update your IP address changes.

[中文文档](README_ZH.md)

## Features

- Supports multiple DNS providers:
  - Tencent Cloud DNSPod
  - Cloudflare
  - Alibaba Cloud (Aliyun)
- IPv4 and IPv6 support
- Automatic IP address detection
- Configurable update intervals
- Docker support for easy deployment
- Lightweight and efficient

## Supported DNS Providers

| Provider        | China Mainland | International | Environment Variables |
|-----------------|----------------|---------------|------------------------|
| DNSPod          | ✅              | ✅             | `dnspod`               |
| Cloudflare      | ✅              | ✅             | `cloudflare`           |
| Alibaba Cloud   | ✅              | ✅             | `aliyun`, `alibabacloud` |

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Docker (optional, for containerized deployment)
- API credentials for your chosen DNS provider

### Configuration

Create a `.env` file based on [.env.simple](file:///home/magicat/projects/ddnsd/.env.simple):

```bash
# DNS provider information
# Available values: dnspod, cloudflare, aliyun, alibabacloud
DNS_PROVIDER=dnspod

# Your DNS provider credentials
SECRET_ID=your_secret_id
SECRET_KEY=your_secret_key

# IPv4/IPv6 switch
IPV4_ENABLED=true
IPV6_ENABLED=true

# Domain information
IPV4_DOMAIN=example.com
IPV4_SUBDOMAINS=ipv4,ipv4.www
IPV6_DOMAIN=example.com
IPV6_SUBDOMAINS=ipv6,ipv6.www

# Update interval (seconds), default 300 seconds (5 minutes)
INTERVAL=300
```

### Running with Go

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ddnsd.git
   cd ddnsd
   ```

2. Create your `.env` file with your configuration

3. Run the service:
   ```bash
   go run main.go
   ```

### Running with Docker

1. Build the Docker image:
   ```bash
   docker build -t ddnsd .
   ```

2. Run the container with your configuration:
   ```bash
   docker run -d --env-file .env ddnsd
   ```

### Building

To build the project:

```bash
go build -o ddnsd .
```

To build with Docker:

```bash
./build.sh
```

## Configuration Details

### DNSPod Configuration

For Tencent Cloud DNSPod, you need to obtain your API credentials from the [Tencent Cloud Console](https://console.cloud.tencent.com/cam/capi).

```env
DNS_PROVIDER=dnspod
SECRET_ID=your_tencent_cloud_secret_id
SECRET_KEY=your_tencent_cloud_secret_key
```

### Cloudflare Configuration

For Cloudflare, you need your email and API key from the [Cloudflare Dashboard](https://dash.cloudflare.com/profile/api-tokens).

```env
DNS_PROVIDER=cloudflare
SECRET_ID=your_email@example.com
SECRET_KEY=your_cloudflare_api_key
```

### Alibaba Cloud Configuration

For Alibaba Cloud, you need your AccessKey ID and AccessKey Secret from the [Alibaba Cloud Console](https://ram.console.aliyun.com/manage/ak).

```env
# For China region
DNS_PROVIDER=aliyun
SECRET_ID=your_access_key_id
SECRET_KEY=your_access_key_secret

# For international region
DNS_PROVIDER=alibabacloud
SECRET_ID=your_access_key_id
SECRET_KEY=your_access_key_secret
```

## Environment Variables

| Variable            | Description                        | Default Value                         |
|---------------------|------------------------------------|---------------------------------------|
| DNS_PROVIDER        | DNS provider to use                | `dnspod`                              |
| SECRET_ID           | Provider-specific credential       | (required)                            |
| SECRET_KEY          | Provider-specific credential       | (required)                            |
| IPV4_ENABLED        | Enable IPv4 updates                | `true`                                |
| IPV6_ENABLED        | Enable IPv6 updates                | `false`                               |
| IPV4_DOMAIN         | Main domain for IPv4 records       | (required if IPv4 enabled)            |
| IPV4_SUBDOMAINS     | Comma-separated IPv4 subdomains    | (required if IPv4 enabled)            |
| IPV6_DOMAIN         | Main domain for IPv6 records       | (required if IPv6 enabled)            |
| IPV6_SUBDOMAINS     | Comma-separated IPv6 subdomains    | (required if IPv6 enabled)            |
| INTERVAL            | Update interval in seconds         | `300` (5 minutes)                     |
| IPV4_CHECK_URL      | Service to check IPv4 address      | `https://iplark.com/ipapi/public/ip`  |
| IPV6_CHECK_URL      | Service to check IPv6 address      | `https://6.iplark.com/ip`             |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Tencent Cloud SDK for Go](https://github.com/TencentCloud/tencentcloud-sdk-go)
- [Godotenv](https://github.com/joho/godotenv)
- [Cron](https://github.com/robfig/cron)
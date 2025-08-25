package internal

import (
	"ddnsd/config"
	"fmt"
)

// DNSRecord represents a DNS record
type DNSRecord struct {
	RecordID string
	Value    string
}

// DNSProvider defines the interface for DNS providers
type DNSProvider interface {
	GetRecord(domain, subdomain, recordType string) (*DNSRecord, error)
	CreateRecord(domain, subdomain, recordType, value string) (string, error)
	UpdateRecord(recordID, domain, subdomain, recordType, value string) error
}

// NewDNSProvider creates a new DNS provider based on configuration
func NewDNSProvider(cfg *config.Config) (DNSProvider, error) {
	switch cfg.Provider {
	case "dnspod":
		return newDNSPodProvider(cfg.SecretID, cfg.SecretKey)
	case "cloudflare":
		return newCloudflareProvider(cfg.SecretID, cfg.SecretKey)
	case "aliyun":
		return newAliyunProvider(cfg.SecretID, cfg.SecretKey, true) // China edition
	case "alibabacloud":
		return newAliyunProvider(cfg.SecretID, cfg.SecretKey, false) // International edition
	default:
		return nil, fmt.Errorf("unsupported DNS provider: %s", cfg.Provider)
	}
}
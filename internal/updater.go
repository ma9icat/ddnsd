package internal

import (
	"ddnsd/config"
	"ddnsd/utils"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// RunSequentialUpdates performs IPv4 and IPv6 updates in sequence
func RunSequentialUpdates(provider DNSProvider, cfg *config.Config) {
	if cfg.IPv4Enabled {
		utils.WithLogPrefix("[IPv4] ", func() {
			updateIPv4Records(provider, cfg)
		})
	}

	if cfg.IPv6Enabled {
		utils.WithLogPrefix("[IPv6] ", func() {
			updateIPv6Records(provider, cfg)
		})
	}
}

// updateIPv4Records updates all IPv4 DNS records
func updateIPv4Records(provider DNSProvider, cfg *config.Config) {
	utils.LogInfo("Starting record update")
	ipv4, err := getPublicIP(cfg.IPv4CheckURL, "IPv4")
	if err != nil {
		utils.LogError("Update failed: Error getting IP address - %v", err)
		utils.LogInfo("Update completed")
		return
	}
	utils.LogInfo("Current IP address: %s", ipv4)

	for _, subDomain := range cfg.IPv4SubDomains {
		fullDomain := fmt.Sprintf("%s.%s", subDomain, cfg.IPv4Domain)
		utils.LogInfo("Processing subdomain: %s", fullDomain)

		if err := updateRecord(provider, cfg.IPv4Domain, subDomain, ipv4, "A"); err != nil {
			utils.LogError("Subdomain update failed: %s - %v", subDomain, err)
		}
	}

	utils.LogInfo("Update completed")
}

// updateIPv6Records updates all IPv6 DNS records
func updateIPv6Records(provider DNSProvider, cfg *config.Config) {
	utils.LogInfo("Starting record update")
	ipv6, err := getPublicIP(cfg.IPv6CheckURL, "IPv6")
	if err != nil {
		utils.LogError("Update failed: Error getting IP address - %v", err)
		utils.LogInfo("Update completed")
		return
	}
	utils.LogInfo("Current IP address: %s", ipv6)

	for _, subDomain := range cfg.IPv6SubDomains {
		fullDomain := fmt.Sprintf("%s.%s", subDomain, cfg.IPv6Domain)
		utils.LogInfo("Processing subdomain: %s", fullDomain)

		if err := updateRecord(provider, cfg.IPv6Domain, subDomain, ipv6, "AAAA"); err != nil {
			utils.LogError("Subdomain update failed: %s - %v", subDomain, err)
		}
	}

	utils.LogInfo("Update completed")
}

// getPublicIP retrieves public IP address from specified URL
func getPublicIP(url, ipType string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP response error: status code=%d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", fmt.Errorf("empty response, no %s address obtained", ipType)
	}

	return ip, nil
}

// updateRecord creates or updates a single DNS record
func updateRecord(provider DNSProvider, domain, subDomain, ip, recordType string) error {
	record, err := provider.GetRecord(domain, subDomain, recordType)
	if err != nil {
		return fmt.Errorf("failed to query record: %v", err)
	}

	if record != nil {
		if record.Value == ip {
			utils.LogInfo("IP address unchanged, no update needed")
			return nil
		}

		// Update existing record
		if err := provider.UpdateRecord(record.RecordID, domain, subDomain, recordType, ip); err != nil {
			return fmt.Errorf("failed to modify record: %v", err)
		}
		utils.LogInfo("Record updated successfully")
		return nil
	}

	// Create new record
	recordID, err := provider.CreateRecord(domain, subDomain, recordType, ip)
	if err != nil {
		return fmt.Errorf("failed to create record: %v", err)
	}
	utils.LogInfo("Record created successfully, ID=%s", recordID)
	return nil
}
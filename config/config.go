package config

import (
	"ddnsd/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	Provider       string
	SecretID       string
	SecretKey      string
	Interval       int
	IPv4Enabled    bool
	IPv6Enabled    bool
	IPv4Domain     string
	IPv4SubDomains []string
	IPv4CheckURL   string
	IPv6Domain     string
	IPv6SubDomains []string
	IPv6CheckURL   string
}

// LoadConfig loads and validates configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Provider:       getEnv("DNS_PROVIDER", "dnspod"),
		SecretID:       getEnv("SECRET_ID", ""),
		SecretKey:      getEnv("SECRET_KEY", ""),
		IPv4Domain:     getEnv("IPV4_DOMAIN", ""),
		IPv6Domain:     getEnv("IPV6_DOMAIN", ""),
		IPv4Enabled:    getEnvAsBool("IPV4_ENABLED", true),
		IPv6Enabled:    getEnvAsBool("IPV6_ENABLED", false),
		IPv4CheckURL:   getEnv("IPV4_CHECK_URL", "https://iplark.com/ipapi/public/ip"),
		IPv6CheckURL:   getEnv("IPV6_CHECK_URL", "https://6.iplark.com/ip"),
		IPv4SubDomains: parseSubDomains(getEnv("IPV4_SUBDOMAINS", "")),
		IPv6SubDomains: parseSubDomains(getEnv("IPV6_SUBDOMAINS", "")),
	}

	// Parse interval with validation
	intervalStr := getEnv("INTERVAL", "300")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval < 30 {
		return nil, fmt.Errorf("invalid INTERVAL value: must be an integer ≥30")
	}
	cfg.Interval = interval

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks configuration for required values
func (c *Config) validate() error {
	if c.SecretID == "" || c.SecretKey == "" {
		return fmt.Errorf("SECRET_ID and SECRET_KEY must be set")
	}

	// Check for Cloudflare specific configuration
	if c.Provider == "cloudflare" {
		// For Cloudflare, we might want to validate Zone ID is provided
	}

	// Check for Aliyun specific configuration
	if c.Provider == "aliyun" || c.Provider == "alibabacloud" {
		// Add any Aliyun-specific validation if needed
	}

	if !c.IPv4Enabled && !c.IPv6Enabled {
		return fmt.Errorf("at least one of IPv4 or IPv6 must be enabled")
	}

	if c.IPv4Enabled {
		if c.IPv4Domain == "" {
			return fmt.Errorf("IPV4_DOMAIN must be set when IPv4 is enabled")
		}
		if len(c.IPv4SubDomains) == 0 {
			return fmt.Errorf("IPV4_SUBDOMAINS must be set when IPv4 is enabled")
		}
	}

	if c.IPv6Enabled {
		if c.IPv6Domain == "" {
			return fmt.Errorf("IPV6_DOMAIN must be set when IPv6 is enabled")
		}
		if len(c.IPv6SubDomains) == 0 {
			return fmt.Errorf("IPV6_SUBDOMAINS must be set when IPv6 is enabled")
		}
	}

	return nil
}

// getEnv returns environment variable or default value
func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// getEnvAsBool converts environment variable to boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	val = strings.ToLower(val)
	return val == "true" || val == "1" || val == "yes" || val == "on"
}

// parseSubDomains converts comma-separated string to slice
func parseSubDomains(subDomainsStr string) []string {
	var subDomains []string
	for _, s := range strings.Split(subDomainsStr, ",") {
		if trimmed := strings.TrimSpace(s); trimmed != "" {
			subDomains = append(subDomains, trimmed)
		}
	}
	return subDomains
}

// PrintConfigSummary displays configuration overview
func PrintConfigSummary(cfg *Config) {
	utils.LogInfo("\n=== Configuration Summary ===")
	utils.LogInfo("DNS Provider: %s", cfg.Provider)
	utils.LogInfo("Update Interval: %d seconds", cfg.Interval)

	if cfg.IPv4Enabled {
		utils.LogInfo("IPv4: Enabled=%v, Domain=%s, Subdomains=%v",
			cfg.IPv4Enabled, cfg.IPv4Domain, cfg.IPv4SubDomains)
	}

	if cfg.IPv6Enabled {
		utils.LogInfo("IPv6: Enabled=%v, Domain=%s, Subdomains=%v",
			cfg.IPv6Enabled, cfg.IPv6Domain, cfg.IPv6SubDomains)
	}
}
package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CloudflareProvider implements DNSProvider for Cloudflare
type CloudflareProvider struct {
	apiKey   string
	apiEmail string
	zoneID   string
	client   *http.Client
}

// Cloudflare API structures
type cloudflareRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type cloudflareResponse struct {
	Success bool               `json:"success"`
	Errors  []cloudflareError  `json:"errors"`
	Result  []cloudflareRecord `json:"result"`
}

type cloudflareSingleResponse struct {
	Success bool              `json:"success"`
	Errors  []cloudflareError `json:"errors"`
	Result  cloudflareRecord  `json:"result"`
}

type cloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type cloudflareCreateRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
}

type cloudflareUpdateRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
}

// newCloudflareProvider creates a new Cloudflare provider instance
func newCloudflareProvider(apiKey, apiEmail string) (*CloudflareProvider, error) {
	return &CloudflareProvider{
		apiKey:   apiKey,
		apiEmail: apiEmail,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetRecord retrieves an existing DNS record
func (c *CloudflareProvider) GetRecord(domain, subdomain, recordType string) (*DNSRecord, error) {
	var fullName string
	if subdomain == "@" {
		fullName = domain
	} else {
		fullName = subdomain + "." + domain
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s&type=%s", c.zoneID, fullName, recordType)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Auth-Key", c.apiKey)
	req.Header.Set("X-Auth-Email", c.apiEmail)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !cfResp.Success {
		if len(cfResp.Errors) > 0 {
			return nil, fmt.Errorf("API request failed: %s", cfResp.Errors[0].Message)
		}
		return nil, fmt.Errorf("API request failed")
	}

	if len(cfResp.Result) == 0 {
		return nil, nil
	}

	record := cfResp.Result[0]
	return &DNSRecord{
		RecordID: record.ID,
		Value:    record.Content,
	}, nil
}

// CreateRecord creates a new DNS record
func (c *CloudflareProvider) CreateRecord(domain, subdomain, recordType, value string) (string, error) {
	var fullName string
	if subdomain == "@" {
		fullName = domain
	} else {
		fullName = subdomain + "." + domain
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", c.zoneID)

	createReq := cloudflareCreateRequest{
		Type:    recordType,
		Name:    fullName,
		Content: value,
		Proxied: false,
	}

	body, err := json.Marshal(createReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Auth-Key", c.apiKey)
	req.Header.Set("X-Auth-Email", c.apiEmail)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var cfResp cloudflareSingleResponse
	if err := json.Unmarshal(respBody, &cfResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if !cfResp.Success {
		if len(cfResp.Errors) > 0 {
			return "", fmt.Errorf("API request failed: %s", cfResp.Errors[0].Message)
		}
		return "", fmt.Errorf("API request failed")
	}

	return cfResp.Result.ID, nil
}

// UpdateRecord updates an existing DNS record
func (c *CloudflareProvider) UpdateRecord(recordID, domain, subdomain, recordType, value string) error {
	var fullName string
	if subdomain == "@" {
		fullName = domain
	} else {
		fullName = subdomain + "." + domain
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", c.zoneID, recordID)

	updateReq := cloudflareUpdateRequest{
		Type:    recordType,
		Name:    fullName,
		Content: value,
		Proxied: false,
	}

	body, err := json.Marshal(updateReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Auth-Key", c.apiKey)
	req.Header.Set("X-Auth-Email", c.apiEmail)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var cfResp cloudflareSingleResponse
	if err := json.Unmarshal(respBody, &cfResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if !cfResp.Success {
		if len(cfResp.Errors) > 0 {
			return fmt.Errorf("API request failed: %s", cfResp.Errors[0].Message)
		}
		return fmt.Errorf("API request failed")
	}

	return nil
}
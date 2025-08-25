package internal

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AliyunProvider implements DNSProvider for Alibaba Cloud
type AliyunProvider struct {
	accessKeyID     string
	accessKeySecret string
	endpoint        string
	client          *http.Client
}

// Aliyun API structures
type aliyunDescribeResponse struct {
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	Records    struct {
		Record []struct {
			RecordId   string `json:"RecordId"`
			DomainName string `json:"DomainName"`
			RR         string `json:"RR"`
			Type       string `json:"Type"`
			Value      string `json:"Value"`
		} `json:"Record"`
	} `json:"DomainRecords"`
}

type aliyunCreateResponse struct {
	RequestId string `json:"RequestId"`
	RecordId  string `json:"RecordId"`
}

type aliyunUpdateResponse struct {
	RequestId string `json:"RequestId"`
}

// newAliyunProvider creates a new Aliyun provider instance
func newAliyunProvider(accessKeyID, accessKeySecret string, isChina bool) (*AliyunProvider, error) {
	endpoint := "alidns.cn-hangzhou.aliyuncs.com" // China edition
	if !isChina {
		endpoint = "alidns.ap-northeast-1.aliyuncs.com" // International edition
	}

	return &AliyunProvider{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		endpoint:        endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetRecord retrieves an existing DNS record
func (a *AliyunProvider) GetRecord(domain, subdomain, recordType string) (*DNSRecord, error) {
	params := map[string]string{
		"Action":           "DescribeSubDomainRecords",
		"SubDomain":        subdomain + "." + domain,
		"Type":             recordType,
		"Format":           "JSON",
		"Version":          "2015-01-09",
		"AccessKeyId":      a.accessKeyID,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	signature := a.generateSignature("GET", params)
	params["Signature"] = signature

	// Build URL
	baseURL := fmt.Sprintf("http://%s/", a.endpoint)
	var queryParts []string
	for k, v := range params {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	sort.Strings(queryParts)
	requestURL := baseURL + "?" + strings.Join(queryParts, "&")

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var aliResp aliyunDescribeResponse
	if err := json.Unmarshal(body, &aliResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if aliResp.TotalCount == 0 {
		return nil, nil
	}

	record := aliResp.Records.Record[0]
	return &DNSRecord{
		RecordID: record.RecordId,
		Value:    record.Value,
	}, nil
}

// CreateRecord creates a new DNS record
func (a *AliyunProvider) CreateRecord(domain, subdomain, recordType, value string) (string, error) {
	params := map[string]string{
		"Action":           "AddDomainRecord",
		"DomainName":       domain,
		"RR":               subdomain,
		"Type":             recordType,
		"Value":            value,
		"Format":           "JSON",
		"Version":          "2015-01-09",
		"AccessKeyId":      a.accessKeyID,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	signature := a.generateSignature("GET", params)
	params["Signature"] = signature

	// Build URL
	baseURL := fmt.Sprintf("http://%s/", a.endpoint)
	var queryParts []string
	for k, v := range params {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	sort.Strings(queryParts)
	requestURL := baseURL + "?" + strings.Join(queryParts, "&")

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var aliResp aliyunCreateResponse
	if err := json.Unmarshal(body, &aliResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return aliResp.RecordId, nil
}

// UpdateRecord updates an existing DNS record
func (a *AliyunProvider) UpdateRecord(recordID, domain, subdomain, recordType, value string) error {
	params := map[string]string{
		"Action":           "UpdateDomainRecord",
		"RecordId":         recordID,
		"RR":               subdomain,
		"Type":             recordType,
		"Value":            value,
		"Format":           "JSON",
		"Version":          "2015-01-09",
		"AccessKeyId":      a.accessKeyID,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	signature := a.generateSignature("GET", params)
	params["Signature"] = signature

	// Build URL
	baseURL := fmt.Sprintf("http://%s/", a.endpoint)
	var queryParts []string
	for k, v := range params {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	sort.Strings(queryParts)
	requestURL := baseURL + "?" + strings.Join(queryParts, "&")

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var aliResp aliyunUpdateResponse
	if err := json.Unmarshal(body, &aliResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	return nil
}

// generateSignature generates the signature for Aliyun API requests
func (a *AliyunProvider) generateSignature(method string, params map[string]string) string {
	// Sort parameters
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build canonicalized query string
	var queryParts []string
	for _, k := range keys {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(params[k])))
	}
	canonicalizedQueryString := strings.Join(queryParts, "&")

	// Build string to sign
	stringToSign := fmt.Sprintf("%s&%s&%s", method, url.QueryEscape("/"), url.QueryEscape(canonicalizedQueryString))

	// Sign the string
	key := a.accessKeySecret + "&"
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}
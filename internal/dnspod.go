package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

// DNSPodProvider implements DNSProvider for DNSPod
type DNSPodProvider struct {
	client *dnspod.Client
}

// newDNSPodProvider creates a new DNSPod provider instance
func newDNSPodProvider(secretID, secretKey string) (*DNSPodProvider, error) {
	credential := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	cpf.HttpProfile.ReqTimeout = 5 // 5 second timeout

	client, err := dnspod.NewClient(credential, "", cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNSPod client: %v", err)
	}

	return &DNSPodProvider{
		client: client,
	}, nil
}

// GetRecord retrieves an existing DNS record
func (d *DNSPodProvider) GetRecord(domain, subdomain, recordType string) (*DNSRecord, error) {
	req := dnspod.NewDescribeRecordListRequest()
	req.Domain = common.StringPtr(domain)
	req.Subdomain = common.StringPtr(subdomain)
	req.RecordType = common.StringPtr(recordType)

	resp, err := d.client.DescribeRecordList(req)
	if err != nil {
		// Ignore "no records" errors
		if strings.Contains(err.Error(), "No records") ||
			strings.Contains(err.Error(), "记录列表为空") ||
			strings.Contains(err.Error(), "RecordListEmpty") {
			return nil, nil
		}
		return nil, fmt.Errorf("API request failed: %v", err)
	}

	if len(resp.Response.RecordList) == 0 {
		return nil, nil
	}

	firstRecord := resp.Response.RecordList[0]
	return &DNSRecord{
		RecordID: fmt.Sprintf("%d", *firstRecord.RecordId),
		Value:    *firstRecord.Value,
	}, nil
}

// CreateRecord creates a new DNS record
func (d *DNSPodProvider) CreateRecord(domain, subdomain, recordType, value string) (string, error) {
	req := dnspod.NewCreateRecordRequest()
	req.Domain = common.StringPtr(domain)
	req.SubDomain = common.StringPtr(subdomain)
	req.RecordType = common.StringPtr(recordType)
	req.RecordLine = common.StringPtr("默认")
	req.Value = common.StringPtr(value)

	resp, err := d.client.CreateRecord(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}

	return fmt.Sprintf("%d", *resp.Response.RecordId), nil
}

// UpdateRecord updates an existing DNS record
func (d *DNSPodProvider) UpdateRecord(recordID, domain, subdomain, recordType, value string) error {
	recordIDUint, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid record ID: %v", err)
	}

	req := dnspod.NewModifyRecordRequest()
	req.Domain = common.StringPtr(domain)
	req.SubDomain = common.StringPtr(subdomain)
	req.RecordType = common.StringPtr(recordType)
	req.RecordLine = common.StringPtr("默认")
	req.Value = common.StringPtr(value)
	req.RecordId = common.Uint64Ptr(recordIDUint)

	_, err = d.client.ModifyRecord(req)
	if err != nil {
		return fmt.Errorf("API request failed: %v", err)
	}

	return nil
}
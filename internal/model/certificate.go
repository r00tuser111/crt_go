package model

// Certificate 表示从crt.sh获取的证书信息
type Certificate struct {
	IssuerCAID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
	ResultCount    int    `json:"result_count"`
}

// SubdomainInfo 表示子域名信息
type SubdomainInfo struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
}

// Result 表示处理后的结果
type Result struct {
	TargetDomain string          `json:"target_domain"`
	Subdomains   []SubdomainInfo `json:"subdomains"`
}

// History 表示域名监控的历史记录
type History struct {
	LastUpdate      string                       `json:"last_update"`
	KnownSubdomains map[string]map[string]string `json:"known_subdomains"` // domain -> subdomain -> timestamp
}

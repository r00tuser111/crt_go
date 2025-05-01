package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"

	"crt_go/internal/config"
	"crt_go/internal/history"
	"crt_go/internal/logger"
	"crt_go/internal/model"
	"crt_go/internal/output"
)

// Run 执行主要的处理逻辑
func Run(cfg *config.Config) error {
	if cfg.Domain == "" && cfg.DomainList == "" {
		return fmt.Errorf("必须指定域名(-d)或域名列表文件(-l)")
	}

	var domains []string
	if cfg.Domain != "" {
		domains = append(domains, cfg.Domain)
	}

	if cfg.DomainList != "" {
		content, err := ioutil.ReadFile(cfg.DomainList)
		if err != nil {
			return fmt.Errorf("读取域名列表文件失败: %v", err)
		}
		for _, domain := range strings.Split(string(content), "\n") {
			if domain = strings.TrimSpace(domain); domain != "" {
				domains = append(domains, domain)
			}
		}
	}

	if cfg.Monitor {
		return monitorDomains(cfg, domains)
	}

	return processDomains(cfg, domains)
}

// processDomains 处理域名列表
func processDomains(cfg *config.Config, domains []string) error {
	results := make([]model.Result, 0)
	semaphore := make(chan struct{}, cfg.Threads)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, domain := range domains {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(d string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			result, err := processOneDomain(d)
			if err != nil {
				if !cfg.Silent {
					logger.Error(fmt.Sprintf("处理域名 %s 失败: %v", d, err))
				}
				return
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()

			if cfg.SimpleOutput {
				// 简单输出模式下直接打印子域名
				for _, subdomain := range result.Subdomains {
					fmt.Println(subdomain.Name)
				}
			} else if !cfg.Silent {
				logger.Info(fmt.Sprintf("发现 %d 个子域名: %s", len(result.Subdomains), d))
			}
		}(domain)
	}

	wg.Wait()

	// 如果不是简单输出模式，则保存结果到文件
	if !cfg.SimpleOutput && (cfg.JSONOutput || cfg.XLSXOutput) {
		return output.SaveResults(cfg, results)
	}

	return nil
}

// isValidSubdomain 检查子域名是否有效
func isValidSubdomain(subdomain string) bool {
	// 过滤包含通配符的域名
	if strings.Contains(subdomain, "*") {
		return false
	}

	// 过滤邮箱地址
	if strings.Contains(subdomain, "@") {
		return false
	}

	// 过滤可能的无效字符
	invalidChars := []string{
		" ",      // 空格
		"{", "}", // 花括号
		"[", "]", // 方括号
		"<", ">", // 尖括号
		"|",  // 竖线
		"\\", // 反斜杠
		"\"", // 双引号
		"'",  // 单引号
		";",  // 分号
	}
	for _, char := range invalidChars {
		if strings.Contains(subdomain, char) {
			return false
		}
	}

	// 检查域名格式
	parts := strings.Split(subdomain, ".")
	for _, part := range parts {
		// 每个部分不能为空
		if part == "" {
			return false
		}
		// 每个部分不能以连字符开始或结束
		if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") {
			return false
		}
		// 每个部分只能包含字母、数字和连字符
		for _, r := range part {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '-') {
				return false
			}
		}
	}

	return true
}

// processOneDomain 处理单个域名
func processOneDomain(domain string) (model.Result, error) {
	url := fmt.Sprintf("https://crt.sh/json?q=%%.%s", domain)
	resp, err := http.Get(url)
	if err != nil {
		return model.Result{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Result{}, fmt.Errorf("读取响应失败: %v", err)
	}

	if len(body) == 0 {
		return model.Result{}, fmt.Errorf("空响应")
	}

	var certs []model.Certificate
	if err := json.Unmarshal(body, &certs); err != nil {
		return model.Result{}, fmt.Errorf("解析JSON失败: %v, 响应内容: %s", err, string(body))
	}

	// 使用 map 进行子域名去重，同时保留最新的时间戳
	subdomains := make(map[string]string)
	domainSuffix := "." + domain
	for _, cert := range certs {
		// 处理 CommonName
		if strings.HasSuffix(cert.CommonName, domainSuffix) &&
			cert.CommonName != domain &&
			isValidSubdomain(cert.CommonName) {
			if timestamp, exists := subdomains[cert.CommonName]; !exists || cert.EntryTimestamp > timestamp {
				subdomains[cert.CommonName] = cert.EntryTimestamp
			}
		}

		// 处理 NameValue
		for _, name := range strings.Split(cert.NameValue, "\n") {
			name = strings.TrimSpace(name)
			if strings.HasSuffix(name, domainSuffix) &&
				name != domain &&
				isValidSubdomain(name) {
				if timestamp, exists := subdomains[name]; !exists || cert.EntryTimestamp > timestamp {
					subdomains[name] = cert.EntryTimestamp
				}
			}
		}
	}

	// 将 map 转换为切片以便排序
	var subdomainList []model.SubdomainInfo
	for name, timestamp := range subdomains {
		subdomainList = append(subdomainList, model.SubdomainInfo{
			Name:      name,
			Timestamp: timestamp,
		})
	}

	// 按时间戳降序排序（最新的在前面）
	sort.Slice(subdomainList, func(i, j int) bool {
		return subdomainList[i].Timestamp > subdomainList[j].Timestamp
	})

	var result model.Result
	result.TargetDomain = domain
	result.Subdomains = subdomainList

	return result, nil
}

// monitorDomains 监控域名的变化
func monitorDomains(cfg *config.Config, domains []string) error {
	// 初始化历史记录管理器
	historyManager, err := history.NewManager("")
	if err != nil {
		return fmt.Errorf("初始化历史记录管理器失败: %v", err)
	}

	results := make([]model.Result, 0)
	for _, domain := range domains {
		result, err := processOneDomain(domain)
		if err != nil {
			if !cfg.Silent {
				logger.Error(fmt.Sprintf("监控域名 %s 失败: %v", domain, err))
			}
			continue
		}

		var newResult model.Result
		newResult.TargetDomain = domain

		// 检查新的子域名
		for _, subdomain := range result.Subdomains {
			if historyManager.IsNewSubdomain(domain, subdomain.Name, subdomain.Timestamp) {
				newResult.Subdomains = append(newResult.Subdomains, subdomain)
				if cfg.SimpleOutput {
					// 简单输出模式下直接打印新发现的子域名
					fmt.Println(subdomain.Name)
				}
			}
		}

		if len(newResult.Subdomains) > 0 {
			results = append(results, newResult)
			if !cfg.SimpleOutput && !cfg.Silent {
				logger.Info(fmt.Sprintf("发现 %d 个新的子域名: %s", len(newResult.Subdomains), domain))
			}
		}
	}

	// 保存历史记录
	if err := historyManager.SaveChanges(); err != nil {
		logger.Error(fmt.Sprintf("保存历史记录失败: %v", err))
	}

	// 如果不是简单输出模式，则保存结果到文件
	if !cfg.SimpleOutput && len(results) > 0 && (cfg.JSONOutput || cfg.XLSXOutput) {
		if err := output.SaveResults(cfg, results); err != nil {
			logger.Error(fmt.Sprintf("保存结果失败: %v", err))
		}
	}

	return nil
}

package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"crt_go/internal/model"
)

const historyFileName = "crt_go_history.json"

// Manager 管理子域名历史记录
type Manager struct {
	history     model.History
	historyPath string
}

// NewManager 创建一个新的历史记录管理器
func NewManager(dataDir string) (*Manager, error) {
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("无法获取用户主目录: %v", err)
		}
		dataDir = filepath.Join(homeDir, ".crt_go")
	}

	// 确保目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %v", err)
	}

	m := &Manager{
		historyPath: filepath.Join(dataDir, historyFileName),
		history: model.History{
			KnownSubdomains: make(map[string]map[string]string),
		},
	}

	// 加载现有历史记录
	if err := m.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("加载历史记录失败: %v", err)
		}
	}

	return m, nil
}

// load 从文件加载历史记录
func (m *Manager) load() error {
	data, err := os.ReadFile(m.historyPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.history)
}

// save 保存历史记录到文件
func (m *Manager) save() error {
	m.history.LastUpdate = time.Now().Format(time.RFC3339)
	data, err := json.MarshalIndent(m.history, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	return os.WriteFile(m.historyPath, data, 0644)
}

// IsNewSubdomain 检查是否为新的子域名
func (m *Manager) IsNewSubdomain(domain, subdomain, timestamp string) bool {
	if m.history.KnownSubdomains[domain] == nil {
		m.history.KnownSubdomains[domain] = make(map[string]string)
	}

	lastTimestamp, exists := m.history.KnownSubdomains[domain][subdomain]
	if !exists || timestamp > lastTimestamp {
		m.history.KnownSubdomains[domain][subdomain] = timestamp
		return true
	}
	return false
}

// SaveChanges 保存当前的变更
func (m *Manager) SaveChanges() error {
	return m.save()
}

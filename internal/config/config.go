package config

// Config 存储所有的配置选项
type Config struct {
	// 查询选项
	Domain     string
	DomainList string

	// 输出选项
	JSONOutput   bool
	XLSXOutput   bool
	SimpleOutput bool // 简单输出模式，只输出子域名到控制台
	OutputFile   string

	// 运行选项
	Threads int
	Silent  bool
	Monitor bool
}

// 默认配置
const (
	DefaultThreads    = 10
	DefaultOutputFile = "results"
)

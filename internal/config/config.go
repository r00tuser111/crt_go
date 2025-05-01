package config

import (
	"flag"
)

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

	Delay int // Delay between requests in milliseconds
}

// 默认配置
const (
	DefaultThreads    = 10
	DefaultOutputFile = "results"
)

func InitConfig() *Config {
	cfg := &Config{}

	// 查询选项
	flag.StringVar(&cfg.Domain, "d", "", "Target domain")
	flag.StringVar(&cfg.Domain, "domain", "", "Target domain")
	flag.StringVar(&cfg.DomainList, "l", "", "File containing list of domains")
	flag.StringVar(&cfg.DomainList, "list", "", "File containing list of domains")

	// 输出选项
	flag.BoolVar(&cfg.JSONOutput, "j", false, "Output results in JSON format")
	flag.BoolVar(&cfg.JSONOutput, "json", false, "Output results in JSON format")
	flag.BoolVar(&cfg.XLSXOutput, "x", false, "Output results in XLSX format")
	flag.BoolVar(&cfg.XLSXOutput, "xlsx", false, "Output results in XLSX format")
	flag.BoolVar(&cfg.SimpleOutput, "s", false, "Output results in simple format")
	flag.BoolVar(&cfg.SimpleOutput, "simple", false, "Output results in simple format")
	flag.StringVar(&cfg.OutputFile, "o", DefaultOutputFile, "Output file")
	flag.StringVar(&cfg.OutputFile, "output", DefaultOutputFile, "Output file")

	// 运行选项
	flag.IntVar(&cfg.Threads, "t", DefaultThreads, "Number of threads")
	flag.IntVar(&cfg.Threads, "threads", DefaultThreads, "Number of threads")
	flag.BoolVar(&cfg.Silent, "silent", false, "Run in silent mode")
	flag.BoolVar(&cfg.Monitor, "m", false, "Run in monitor mode")
	flag.BoolVar(&cfg.Monitor, "monitor", false, "Run in monitor mode")
	flag.IntVar(&cfg.Delay, "delay", 0, "Delay between requests in milliseconds")

	flag.Parse()

	return cfg
}

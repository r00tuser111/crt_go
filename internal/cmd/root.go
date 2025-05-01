package cmd

import (
	"crt_go/internal/banner"
	"crt_go/internal/config"
	"crt_go/internal/core"

	"github.com/spf13/cobra"
)

var (
	cfg     = &config.Config{}
	rootCmd = &cobra.Command{
		Use:   "crt_go",
		Short: "A tool for querying and monitoring subdomains from crt.sh",
		Long: `crt_go is a powerful command line tool that helps you discover and monitor 
subdomains using SSL certificate information from crt.sh.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果指定了 JSON 或 XLSX 输出，则关闭简单输出模式
			if cfg.JSONOutput || cfg.XLSXOutput {
				cfg.SimpleOutput = false
			}

			if !cfg.Silent {
				banner.Show()
			}
			return core.Run(cfg)
		},
	}
)

func init() {
	// 单域名查询
	rootCmd.Flags().StringVarP(&cfg.Domain, "domain", "d", "", "Single domain to query")

	// 域名列表查询
	rootCmd.Flags().StringVarP(&cfg.DomainList, "list", "l", "", "File containing list of domains")

	// 输出格式
	rootCmd.Flags().BoolVar(&cfg.SimpleOutput, "simple", true, "Simple output mode, only print subdomains to stdout (default)")
	rootCmd.Flags().BoolVar(&cfg.JSONOutput, "json", false, "Output results in JSON format")
	rootCmd.Flags().BoolVar(&cfg.XLSXOutput, "xlsx", false, "Output results in XLSX format")

	// 输出文件
	rootCmd.Flags().StringVarP(&cfg.OutputFile, "output", "o", "results", "Output file name (without extension)")

	// 并发数
	rootCmd.Flags().IntVarP(&cfg.Threads, "threads", "t", 10, "Number of concurrent threads")

	// 静默模式
	rootCmd.Flags().BoolVarP(&cfg.Silent, "silent", "s", false, "Silent mode, suppress all logs")

	// 监控模式
	rootCmd.Flags().BoolVarP(&cfg.Monitor, "monitor", "m", false, "Enable monitoring mode")
}

func Execute() error {
	return rootCmd.Execute()
}

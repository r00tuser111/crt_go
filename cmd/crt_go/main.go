package main

import (
	"os"

	"crt_go/internal/cmd"
	"crt_go/internal/logger"
)

func main() {
	// 初始化日志
	log := logger.InitLogger()
	defer log.Sync()

	// 执行根命令
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

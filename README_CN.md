# crt_go

基于 crt.sh 的证书透明度日志的快速子域名发现工具。

[English Documentation](README.md)

## 功能特点

- 使用证书透明度日志进行快速子域名枚举
- 自动去重子域名
- 按发现时间排序（最新优先）
- 多种输出格式（简单、JSON、XLSX）
- 域名监控模式，发现新的子域名
- 多线程支持，加快处理速度
- 支持单域名和域名列表处理

## 安装

```bash
go install github.com/r00tuser111/crt_go@latest
```

或从源码构建：

```bash
git clone https://github.com/r00tuser111/crt_go.git
cd crt_go
go build -o crt_go cmd/crt_go/main.go
```

## 使用方法

### 基本用法

简单输出（仅打印子域名）：
```bash
./crt_go -d example.com -s
```

保存为 JSON 文件：
```bash
./crt_go -d example.com --json
```

保存为 XLSX 文件：
```bash
./crt_go -d example.com --xlsx
```

### 高级用法

从文件处理多个域名：
```bash
./crt_go -l domains.txt -s
```

监控模式（检查新的子域名）：
```bash
./crt_go -d example.com -m -s
```

增加并发线程数：
```bash
./crt_go -d example.com -t 10 -s
```

### 监控模式说明

监控模式（`-m`）用于持续跟踪子域名的变化。启用后：

1. 工具会在 `~/.crt_go/history.json` 维护历史记录文件
2. 只显示新发现的子域名（之前运行中未发现过的）
3. 所有发现都会自动记录到历史文件中

监控示例：

```bash
# 基础监控（只显示新子域名）
./crt_go -d example.com -m -s

# 监控并保存新发现到 JSON 文件
./crt_go -d example.com -m --json

# 监控文件中的多个域名
./crt_go -l domains.txt -m -s

# 监控并将新发现追加到日志文件
./crt_go -d example.com -m -s >> new_subdomains.log

# 监控并通过邮件发送新发现
./crt_go -d example.com -m -s | mail -s "example.com 的新子域名" your@email.com

# 监控多个域名，仅在发现新子域名时发送邮件
./crt_go -l domains.txt -m -s | grep . && ./crt_go -l domains.txt -m -s | mail -s "发现新的子域名" your@email.com
```

使用 cron 实现自动化：
```bash
# 每6小时检查一次
0 */6 * * * /path/to/crt_go -d example.com -m -s >> /path/to/findings.log

# 每天凌晨2点检查多个域名
0 2 * * * /path/to/crt_go -l /path/to/domains.txt -m --json

# 每周日凌晨3点检查
0 3 * * 0 /path/to/crt_go -d example.com -m --json

# 每12小时检查一次，发现新子域名时发送邮件通知
0 */12 * * * /path/to/crt_go -d example.com -m -s | grep . && /path/to/crt_go -d example.com -m -s | mail -s "发现新的子域名" your@email.com

# 每天早上8点生成多域名报告并以JSON附件形式发送（仅包含最新结果）
0 8 * * * REPORT_FILE=$(date +\%Y\%m\%d)_domains_report.json && /path/to/crt_go -l domains.txt -m --json -o $REPORT_FILE && cat $REPORT_FILE | mail -s "每日子域名报告" -a $REPORT_FILE your@email.com && rm -f $REPORT_FILE
```

历史记录文件可以备份或在团队成员之间共享，以在不同机器上保持一致的监控记录。

### 输出示例

简单输出：
```bash
$ ./crt_go -d example.com -s
www.example.com
dev.example.com
mail.example.com
```

JSON 输出格式：
```json
{
  "target_domain": "example.com",
  "subdomains": [
    {
      "name": "www.example.com",
      "timestamp": "2024-01-27T14:59:38.909"
    }
  ]
}
```

## 参数选项

- `-d, --domain`: 目标域名
- `-l, --list`: 包含域名列表的文件
- `-s, --simple`: 简单输出模式（仅子域名）
- `-j, --json`: 保存为 JSON 格式
- `-x, --xlsx`: 保存为 XLSX 格式
- `-m, --monitor`: 监控模式，用于发现新的子域名
- `-t, --threads`: 并发线程数（默认：5）
- `--silent`: 抑制信息输出

## 贡献

欢迎提交 Pull Requests！对于重大更改，请先开 issue 讨论您想要更改的内容。

## 许可证

[MIT](LICENSE) 
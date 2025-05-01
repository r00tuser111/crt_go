# crt_go

A fast subdomain discovery tool based on the Certificate Transparency logs from crt.sh.

[中文文档](README_CN.md)

## Features

- Fast subdomain enumeration using Certificate Transparency logs
- Automatic deduplication of subdomains
- Results sorted by discovery time (newest first)
- Multiple output formats (Simple, JSON, XLSX)
- Domain monitoring mode for new subdomain discovery
- Multi-threading support for faster processing
- Support for both single domain and domain list processing

## Installation

```bash
go install github.com/r00tuser111/crt_go@latest
```

Or build from source:

```bash
git clone https://github.com/r00tuser111/crt_go.git
cd crt_go
go build -o crt_go cmd/crt_go/main.go
```

## Usage

### Basic Usage

Simple output (prints subdomains only):
```bash
./crt_go -d example.com -s
```

Save results to JSON file:
```bash
./crt_go -d example.com --json
```

Save results to XLSX file:
```bash
./crt_go -d example.com --xlsx
```

### Advanced Usage

Process multiple domains from a file:
```bash
./crt_go -l domains.txt -s
```

Monitor mode (check for new subdomains):
```bash
./crt_go -d example.com -m -s
```

Increase concurrent threads:
```bash
./crt_go -d example.com -t 10 -s
```

### Monitoring Mode

The monitoring mode (`-m`) is designed to track new subdomains over time. When enabled:

1. The tool maintains a history file at `~/.crt_go/history.json`
2. Only newly discovered subdomains (not seen in previous runs) are displayed
3. All discoveries are automatically recorded in the history file

Example monitoring setups:

```bash
# Basic monitoring (show only new subdomains)
./crt_go -d example.com -m -s

# Monitor and save new findings to JSON
./crt_go -d example.com -m --json

# Monitor multiple domains from a file
./crt_go -l domains.txt -m -s

# Monitor and append new findings to a log file
./crt_go -d example.com -m -s >> new_subdomains.log

# Monitor and send new findings via email
./crt_go -d example.com -m -s | mail -s "New subdomains for example.com" your@email.com

# Monitor multiple domains and send email only if new subdomains found
./crt_go -l domains.txt -m -s | grep . && ./crt_go -l domains.txt -m -s | mail -s "New subdomains found" your@email.com
```

Automation with cron:
```bash
# Check every 6 hours
0 */6 * * * /path/to/crt_go -d example.com -m -s >> /path/to/findings.log

# Daily check at 2 AM for multiple domains
0 2 * * * /path/to/crt_go -l /path/to/domains.txt -m --json

# Weekly check on Sunday at 3 AM
0 3 * * 0 /path/to/crt_go -d example.com -m --json

# Send email notification for new subdomains every 12 hours
0 */12 * * * /path/to/crt_go -d example.com -m -s | grep . && /path/to/crt_go -d example.com -m -s | mail -s "New subdomains discovered" your@email.com

# Daily report for multiple domains with JSON attachment (latest results only)
0 8 * * * REPORT_FILE=$(date +\%Y\%m\%d)_domains_report.json && /path/to/crt_go -l domains.txt -m --json -o $REPORT_FILE && cat $REPORT_FILE | mail -s "Daily subdomain report" -a $REPORT_FILE your@email.com && rm -f $REPORT_FILE
```

The history file can be backed up or shared between team members to maintain consistent monitoring across different machines.

### Output Examples

Simple output:
```bash
$ ./crt_go -d example.com -s
www.example.com
dev.example.com
mail.example.com
```

JSON output format:
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

## Options

- `-d, --domain`: Target domain to enumerate
- `-l, --list`: File containing list of domains
- `-s, --simple`: Simple output mode (subdomains only)
- `-j, --json`: Save results in JSON format
- `-x, --xlsx`: Save results in XLSX format
- `-m, --monitor`: Monitor mode for new subdomains
- `-t, --threads`: Number of concurrent threads (default: 5)
- `--silent`: Suppress information output

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE) 
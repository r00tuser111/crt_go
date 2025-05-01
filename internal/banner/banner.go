package banner

import "github.com/fatih/color"

const banner = `
   ______  ____  ______    ______  ____
  / ____/ / __/ /_  __/   / ____/ / __ \
 / /     / /_    / /     / / __  / / / /
/ /___  / __/   / /     / /_/ / / /_/ /
\____/ /_/     /_/      \____/  \____/

CRT.sh Subdomain Discovery Tool
Version: 1.0.0
Author: kking
`

// Show 显示banner信息
func Show() {
	color.Cyan(banner)
}

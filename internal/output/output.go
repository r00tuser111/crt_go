package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crt_go/internal/config"
	"crt_go/internal/model"

	"github.com/xuri/excelize/v2"
)

// SaveResults 保存结果到文件
func SaveResults(cfg *config.Config, results []model.Result) error {
	if len(results) == 0 {
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")

	if cfg.JSONOutput {
		return saveJSON(cfg.OutputFile+"_"+timestamp+".json", results)
	}

	return saveXLSX(cfg.OutputFile+"_"+timestamp+".xlsx", results)
}

// saveJSON 保存为JSON格式
func saveJSON(filename string, results []model.Result) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// saveXLSX 保存为XLSX格式
func saveXLSX(filename string, results []model.Result) error {
	f := excelize.NewFile()
	defer f.Close()

	for i, result := range results {
		// 使用目标域名作为sheet名，确保sheet名称有效
		sheetName := sanitizeSheetName(result.TargetDomain)

		// 创建新的sheet
		newIndex, err := f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("创建工作表失败: %v", err)
		}

		// 设置标题样式
		titleStyle, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 12,
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#E0E0E0"},
				Pattern: 1,
			},
		})
		if err != nil {
			return fmt.Errorf("创建样式失败: %v", err)
		}

		// 设置标题
		f.SetCellValue(sheetName, "A1", "子域名")
		f.SetCellValue(sheetName, "B1", "发现时间")
		f.SetCellStyle(sheetName, "A1", "B1", titleStyle)

		// 写入数据
		for i, subdomain := range result.Subdomains {
			row := i + 2
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), subdomain.Name)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), subdomain.Timestamp)
		}

		// 设置列宽
		f.SetColWidth(sheetName, "A", "A", 50)
		f.SetColWidth(sheetName, "B", "B", 30)

		// 冻结标题行
		f.SetPanes(sheetName, &excelize.Panes{
			Freeze:      true,
			Split:       false,
			XSplit:      0,
			YSplit:      1,
			TopLeftCell: "A2",
			ActivePane:  "bottomLeft",
		})

		// 设置为活动sheet
		if i == 0 {
			f.SetActiveSheet(newIndex)
		}
	}
	// 删除默认的 Sheet1
	f.DeleteSheet("Sheet1")

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	return f.SaveAs(filename)
}

// sanitizeSheetName 确保sheet名称符合Excel的要求
func sanitizeSheetName(name string) string {
	// Excel sheet名称的限制：
	// 1. 长度不能超过31个字符
	// 2. 不能包含以下字符: [ ] * ? / \
	// 3. 不能为空

	// 替换无效字符
	invalid := []string{"[", "]", "*", "?", "/", "\\", ":"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}

	// 限制长度
	if len(result) > 31 {
		result = result[:31]
	}

	// 确保不为空
	if result == "" {
		result = "Sheet1"
	}

	return result
}

package helpers

import (
	"fmt"
	"html/template"
	"strings"
)

var FuncMap = template.FuncMap{
	"formatNumber": func(n int64) string {
		if n >= 1000000000 {
			return fmt.Sprintf("%.1fB", float64(n)/1000000000)
		} else if n >= 1000000 {
			return fmt.Sprintf("%.1fM", float64(n)/1000000)
		} else if n >= 1000 {
			return fmt.Sprintf("%.1fK", float64(n)/1000)
		}
		return fmt.Sprintf("%d", n)
	},
	"formatBytes": func(b int64) string {
		if b >= 1073741824 {
			return fmt.Sprintf("%.2f GB", float64(b)/1073741824)
		} else if b >= 1048576 {
			return fmt.Sprintf("%.2f MB", float64(b)/1048576)
		} else if b >= 1024 {
			return fmt.Sprintf("%.2f KB", float64(b)/1024)
		}
		return fmt.Sprintf("%d B", b)
	},
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"title": strings.Title,
	"add": func(a, b int) int {
		return a + b
	},
}

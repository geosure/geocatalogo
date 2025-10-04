package helpers

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

// extractMetadata extracts metadata from abstract text (e.g., "Rows: 123 | Size: 45 MB | Columns: 10")
// Also handles lowercase formats like "8 rows, 4 columns"
func extractMetadata(abstract string, key string) string {
	// Try with colon first (e.g., "Rows: 123")
	pattern := regexp.MustCompile(`(?i)` + key + `:\s*([^|,]+)`)
	matches := pattern.FindStringSubmatch(abstract)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Try without colon (e.g., "8 rows")
	pattern = regexp.MustCompile(`(\d+)\s+` + strings.ToLower(key) + `\b`)
	matches = pattern.FindStringSubmatch(strings.ToLower(abstract))
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

// cleanAbstract removes metadata suffix from abstract, keeping only the descriptive part
func cleanAbstract(abstract string) string {
	// Remove JSON metadata (everything from {' to end)
	if idx := strings.Index(abstract, "{'"); idx != -1 {
		return strings.TrimSpace(abstract[:idx])
	}
	if idx := strings.Index(abstract, `{"`); idx != -1 {
		return strings.TrimSpace(abstract[:idx])
	}
	// Remove everything after "| Notes:" or just "Notes:"
	if idx := strings.Index(abstract, "| Notes:"); idx != -1 {
		return strings.TrimSpace(abstract[:idx])
	}
	if idx := strings.Index(abstract, "Notes:"); idx != -1 {
		return strings.TrimSpace(abstract[:idx])
	}
	return abstract
}

// hasMetadata checks if abstract contains metadata like "Rows:", "Size:", "Columns:" or "8 rows, 4 columns" or JSON
func hasMetadata(abstract string) bool {
	lower := strings.ToLower(abstract)
	return strings.Contains(lower, "rows:") ||
	       strings.Contains(lower, "size:") ||
	       strings.Contains(lower, "columns:") ||
	       strings.Contains(lower, "schema:") ||
	       strings.Contains(abstract, "{'") || // JSON metadata
	       regexp.MustCompile(`\d+\s+(rows|columns|features|fields)`).MatchString(lower)
}

// extractJSON extracts a specific field from JSON-like metadata in abstract
func extractJSON(abstract string, key string) string {
	// Look for pattern like 'key': 'value' or "key": "value"
	pattern := regexp.MustCompile(`['"]` + key + `['"]\s*:\s*['"]([^'"]+)['"]`)
	matches := pattern.FindStringSubmatch(abstract)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// hasJSON checks if abstract contains JSON-like metadata
func hasJSON(abstract string) bool {
	return strings.Contains(abstract, "{'") || strings.Contains(abstract, `{"`)
}

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
	"extractMetadata": extractMetadata,
	"cleanAbstract":   cleanAbstract,
	"hasMetadata":     hasMetadata,
	"extractJSON":     extractJSON,
	"hasJSON":         hasJSON,
}

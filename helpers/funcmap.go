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

// CountryCodeToFlag converts ISO 3166-1 alpha-2 country code to flag emoji
func CountryCodeToFlag(code string) string {
	if code == "" {
		return ""
	}

	code = strings.ToUpper(code)

	// Convert to regional indicator symbols (Unicode)
	// A = U+1F1E6, so offset from 'A' (0x41)
	var flag strings.Builder
	for _, c := range code {
		if c >= 'A' && c <= 'Z' {
			flag.WriteRune(rune(0x1F1E6 + (c - 'A')))
		}
	}

	return flag.String()
}

// CountryCodeToName converts country code to full name
func CountryCodeToName(code string) string {
	names := map[string]string{
		"ai": "Anguilla", "ar": "Argentina", "au": "Australia", "bd": "Bangladesh",
		"br": "Brazil", "ca": "Canada", "cd": "DR Congo", "cl": "Chile",
		"cn": "China", "co": "Colombia", "cu": "Cuba", "de": "Germany",
		"dm": "Dominica", "ee": "Estonia", "eg": "Egypt", "es": "Spain",
		"et": "Ethiopia", "eu": "European Union", "fr": "France", "gh": "Ghana",
		"gt": "Guatemala", "hu": "Hungary", "id": "Indonesia", "in": "India",
		"iq": "Iraq", "ir": "Iran", "it": "Italy", "jp": "Japan",
		"km": "Comoros", "kr": "South Korea", "mc": "Monaco", "ml": "Mali",
		"mm": "Myanmar", "mx": "Mexico", "na": "Namibia", "ng": "Nigeria",
		"nz": "New Zealand", "pe": "Peru", "ph": "Philippines", "pk": "Pakistan",
		"pt": "Portugal", "ru": "Russia", "sa": "Saudi Arabia", "sc": "Seychelles",
		"sd": "Sudan", "sg": "Singapore", "sl": "Sierra Leone", "sy": "Syria",
		"th": "Thailand", "tr": "Turkey", "tz": "Tanzania", "uk": "United Kingdom",
		"us": "United States", "vn": "Vietnam", "za": "South Africa",
	}

	if name, ok := names[strings.ToLower(code)]; ok {
		return name
	}
	return strings.ToUpper(code)
}

// ContinentToEmoji returns emoji for continent
func ContinentToEmoji(continent string) string {
	emojis := map[string]string{
		"africa":        "🌍",
		"asia":          "🌏",
		"europe":        "🌍",
		"north-america": "🌎",
		"south-america": "🌎",
		"oceania":       "🌏",
		"global":        "🌐",
	}

	if emoji, ok := emojis[strings.ToLower(continent)]; ok {
		return emoji
	}
	return "🌍"
}

// ContinentToName capitalizes continent name nicely
func ContinentToName(continent string) string {
	names := map[string]string{
		"africa":        "Africa",
		"asia":          "Asia",
		"europe":        "Europe",
		"north-america": "North America",
		"south-america": "South America",
		"oceania":       "Oceania",
		"global":        "Global",
	}

	if name, ok := names[strings.ToLower(continent)]; ok {
		return name
	}
	return strings.Title(continent)
}

// V6JobToGitHubURL converts a v6 job file path to GitHub URL
func V6JobToGitHubURL(v6JobFile string) string {
	if v6JobFile == "" {
		return ""
	}
	// Remove v6/ prefix to get jobs/...
	githubPath := strings.Replace(v6JobFile, "v6/", "", 1)
	return fmt.Sprintf("https://github.com/geosure/v6/blob/main/%s", githubPath)
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
	// Geography helpers
	"countryFlag":     CountryCodeToFlag,
	"countryName":     CountryCodeToName,
	"continentEmoji":  ContinentToEmoji,
	"continentName":   ContinentToName,
	// v6 helpers
	"v6GithubURL":     V6JobToGitHubURL,
}

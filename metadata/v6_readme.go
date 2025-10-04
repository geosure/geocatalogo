package metadata

// V6README represents geographic context from v6 README files
type V6README struct {
	Path      string                 `json:"path"`
	Content   string                 `json:"content"`
	Geography map[string]interface{} `json:"geography"`
	SizeBytes int                    `json:"size_bytes"`
}

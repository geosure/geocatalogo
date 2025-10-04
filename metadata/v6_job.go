package metadata

// V6Job represents a v6 job definition from YAML
type V6Job struct {
	Path       string                 `json:"path"`
	RawYAML    string                 `json:"raw_yaml"`
	Parsed     map[string]interface{} `json:"parsed"`
	DatasetID  string                 `json:"dataset_id"`
}

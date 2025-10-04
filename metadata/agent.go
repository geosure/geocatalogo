package metadata

// AgentIntrospection represents the complete agent introspection data
type AgentIntrospection struct {
	Agents            []Agent         `json:"agents"`
	LeadershipProjects []Agent         `json:"leadership_projects"`
	Metadata          AgentMetadata   `json:"metadata"`
}

// Agent represents a single AI agent's comprehensive introspection
type Agent struct {
	AgentID      string                  `json:"agent_id"`
	Name         string                  `json:"name"`
	Type         string                  `json:"type"`
	Status       string                  `json:"status"`
	DataFormat   string                  `json:"data_format"`
	Location     AgentLocation           `json:"location"`
	Execution    *AgentExecution         `json:"execution,omitempty"`
	OrgContext   *AgentOrgContext        `json:"organizational_context,omitempty"`
	Description  string                  `json:"description"`
	Capabilities *AgentCapabilities      `json:"capabilities,omitempty"`
	Playbook     *AgentPlaybook          `json:"playbook,omitempty"`
	Owner        string                  `json:"owner"`
	LastUpdated  string                  `json:"last_updated"`
	Notes        string                  `json:"notes"`
	PageURL      string                  `json:"page_url"`
	APIURL       string                  `json:"api_url"`
}

// AgentLocation represents the agent's location in the catalog
type AgentLocation struct {
	Continent string `json:"continent"`
	Country   string `json:"country"`
	State     string `json:"state"`
	City      string `json:"city"`
}

// AgentExecution represents how the agent runs
type AgentExecution struct {
	Platform          string   `json:"platform,omitempty"`
	GitHubRepo        string   `json:"github_repo,omitempty"`
	RepoPath          string   `json:"repo_path,omitempty"`
	HowToRun          string   `json:"how_to_run,omitempty"`
	Runtime           string   `json:"runtime,omitempty"`
	Model             string   `json:"model,omitempty"`
	ContextSize       string   `json:"context_size,omitempty"`
	WorkingDirectory  string   `json:"working_directory,omitempty"`
	MaintenanceScript string   `json:"maintenance_script,omitempty"`
	Dependencies      []string `json:"dependencies,omitempty"`
}

// AgentOrgContext represents the agent's organizational context
type AgentOrgContext struct {
	CreatedBy       string   `json:"created_by,omitempty"`
	CreationDate    string   `json:"creation_date,omitempty"`
	CurrentBoss     string   `json:"current_boss,omitempty"`
	CurrentUsers    []string `json:"current_users,omitempty"`
	ReportsTo       []string `json:"reports_to,omitempty"`
	CoordinatesWith []string `json:"coordinates_with,omitempty"`
	Manages         []string `json:"manages,omitempty"`
	Purpose         string   `json:"purpose,omitempty"`
}

// AgentCapabilities represents the agent's capabilities
type AgentCapabilities struct {
	Primary      []string `json:"primary,omitempty"`
	Secondary    []string `json:"secondary,omitempty"`
	Integrations []string `json:"integrations,omitempty"`
	DataSources  []string `json:"data_sources,omitempty"`
	Outputs      []string `json:"outputs,omitempty"`
}

// AgentPlaybook represents the agent's playbook
type AgentPlaybook struct {
	Name            string `json:"name,omitempty"`
	Location        string `json:"location,omitempty"`
	Purpose         string `json:"purpose,omitempty"`
	Workflow        string `json:"workflow,omitempty"`
	UpdateFrequency string `json:"update_frequency,omitempty"`
	EstimatedTime   string `json:"estimated_time,omitempty"`
}

// AgentMetadata represents metadata about the introspection
type AgentMetadata struct {
	IntrospectionDate  string `json:"introspection_date"`
	TotalAgents        int    `json:"total_agents"`
	TotalClaudeProjects int   `json:"total_claude_projects,omitempty"`
	Version            string `json:"version"`
	Schema             string `json:"schema"`
}

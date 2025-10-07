package webui

import (
	"github.com/go-spatial/geocatalogo/helpers"
	"github.com/go-spatial/geocatalogo/metadata"
)

type App struct {
	tc    *helpers.TemplateCache
	meta  *metadata.MetadataStore
}

func NewApp(tc *helpers.TemplateCache, meta *metadata.MetadataStore) *App {
	return &App{
		tc:   tc,
		meta: meta,
	}
}

type PageData struct {
	Records []Record
	Stats   CatalogStats
}

type Record struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Links      []Link     `json:"links,omitempty"`
}

type Link struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url"`
}

type Properties struct {
	Title                  string                 `json:"title"`
	Abstract               string                 `json:"abstract,omitempty"`
	Collection             string                 `json:"collection"`
	GROMetadata            GROMetadata            `json:"gro_metadata"`
	GeoCatalogo            GeoCatalogoMeta        `json:"_geocatalogo"`
	Operational            map[string]interface{} `json:"operational,omitempty"`
	Database               map[string]interface{} `json:"database,omitempty"`
	FileMetadata           map[string]interface{} `json:"file_metadata,omitempty"`
	Execution              map[string]interface{} `json:"execution,omitempty"`
	OrganizationalContext  map[string]interface{} `json:"organizational_context,omitempty"`
	Capabilities           map[string]interface{} `json:"capabilities,omitempty"`
}

type GeoCatalogoMeta struct {
	Inserted string `json:"inserted,omitempty"`
	Source   string `json:"source,omitempty"`
	Schema   string `json:"schema,omitempty"`
}

type GROMetadata struct {
	ImplementationStatus string `json:"implementation_status,omitempty"`
	DataFormat           string `json:"data_format,omitempty"`
	GeographicScope      string `json:"geographic_scope,omitempty"`
	Continent            string `json:"continent,omitempty"`
	Country              string `json:"country,omitempty"`
	StateProvince        string `json:"state_province,omitempty"`
	Admin2               string `json:"admin2,omitempty"`
	City                 string `json:"city,omitempty"`
	Owner                string `json:"owner,omitempty"`
	UpdateFrequency      string `json:"update_frequency,omitempty"`
	S3Path               string `json:"s3_path,omitempty"`
	DatabaseTable        string `json:"database_table,omitempty"`
	V6JobFile            string `json:"v6_job_file,omitempty"`
	V6JobType            string `json:"v6_job_type,omitempty"`
	FileSizeMB           string `json:"file_size_mb,omitempty"`
}

type CatalogStats struct {
	Total         int
	DataRecords   int  // Total - PotentialV6 (actual data sources)
	ExistingDB    int
	ExistingLocal int
	PotentialV6   int
	ExternalAPI   int
	ExternalNews  int
	ExternalGov   int
	ExternalOther int
	Continents    []ContinentStat
	Countries     []CountryStat
	Formats       []FormatStat
}

type FormatStat struct {
	Name  string
	Count int
}

type ContinentStat struct {
	Code  string
	Name  string
	Emoji string
	Count int
}

type CountryStat struct {
	Code  string
	Name  string
	Flag  string
	Count int
}

type StatsPageData struct {
	TotalRecords         int
	AgentCount           int
	BotCount             int
	ServiceCount         int
	DatabaseTableCount   int
	FileCount            int
	V6JobCount           int
	ExternalSourceCount  int
	IntrospectionStats   IntrospectionStats
}

type IntrospectionStats struct {
	DatabaseTables  int
	Agents          int
	CSVFiles        int
	ParquetFiles    int
	Shapefiles      int
	PDFs            int
	V6Jobs          int
}

type DatasetPageData struct {
	Record         Record
	RecordJSON     string
	DatabaseTable  *metadata.DatabaseTable
	CSVFile        *metadata.CSVFile
	ParquetFile    *metadata.ParquetFile
	ShapefileFile  *metadata.ShapefileFile
	GeoPackageFile *metadata.GeoPackageFile
	ExcelFile      *metadata.ExcelFile
	JSONFile       *metadata.JSONFile
	FileGDBFile    *metadata.FileGDBFile
	PNGFile        *metadata.PNGFile
	PDFFile        *metadata.PDFFile
	V6Job          *metadata.V6Job
	Agent          *metadata.Agent
}

type CollectionCounts struct {
	V6Jobs               int
	Database             int
	Files                int
	APIs                 int
	News                 int
	Government           int
	AIAgents             int // Claude, GPT, Gemini, etc.
	ClaudeProjects       int // Claude Projects (cloud-based agents)
	OperationalServices  int // Lambda functions (news-ingestor, location-matcher, etc.)
	DataInspectionBots   int // CSV inspector, Parquet analyzer, Database inspector, etc.
	CatalogManagementBots int // Catalog updater, converter, rebuilder, etc.
	DataBots             int // Legacy/generic data processing bots
	ScraperBots          int // News scraper, ACLED harvester, etc.
	AutoBots             int // Automation bots (orchestrators, job schedulers, etc.)
	HistoricalAgents     int // Way Barrios' archived prompt engineering examples
	VerbApps             int // User-facing verb applications
	InternalTools        int // Internal tools (introspect, validate)
	APIServices          int // API services (tile, search, catalog)
	TeamMembers          int // GRO team members
	Infrastructure       int // AWS infrastructure components
	Other                int
}

type SubGeography struct {
	Code  string
	Name  string
	Emoji string
	Count int
	URL   string // Link to this sub-geography's page
}

type GeographyPageData struct {
	Level            string               // city, county, state, country, continent
	Name             string               // Geographic name
	README           *metadata.V6README   // README content for this geography
	Jobs             []Record             // All jobs matching this geography
	JobCount         int                  // Number of jobs
	CollectionCounts CollectionCounts     // Breakdown by collection type
	SubGeographies   []SubGeography       // Countries (for continent), States (for country), Cities (for state)
	City             string
	County           string
	State            string
	Country          string
	Continent        string
}

type CollectionDetailPageData struct {
	CollectionCode        string
	CollectionName        string
	CollectionEmoji       string
	CollectionDescription string
	TotalCount            int
	Records               []Record
	ByStatus              map[string][]Record // Grouped by implementation_status
}

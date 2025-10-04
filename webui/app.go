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
	Title         string            `json:"title"`
	Abstract      string            `json:"abstract,omitempty"`
	Collection    string            `json:"collection"`
	GROMetadata   GROMetadata       `json:"gro_metadata"`
	GeoCatalogo   GeoCatalogoMeta   `json:"_geocatalogo"`
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
	V6Jobs          int
	Database        int
	Files           int
	APIs            int
	News            int
	Government      int
	AIAgents        int // Claude, GPT, Gemini, etc.
	DataBots        int // CSV inspector, Parquet analyzer, etc.
	ScraperBots     int // News scraper, ACLED harvester, etc.
	AutoBots        int // Catalog updater, S3 sync, job scheduler, etc.
	HistoricalAgents int // Way Barrios' archived prompt engineering examples
	VerbApps        int // User-facing verb applications
	InternalTools   int // Internal tools (introspect, validate)
	APIServices     int // API services (tile, search, catalog)
	TeamMembers     int // GRO team members
	Infrastructure  int // AWS infrastructure components
	Other           int
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

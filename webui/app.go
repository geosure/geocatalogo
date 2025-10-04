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
	Country              string `json:"country,omitempty"`
	Continent            string `json:"continent,omitempty"`
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
	ExistingDB    int
	ExistingLocal int
	PotentialV6   int
	ExternalAPI   int
	ExternalNews  int
	ExternalGov   int
	ExternalOther int
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
}

package webui

import (
	"github.com/go-spatial/geocatalogo/helpers"
)

type App struct {
	tc *helpers.TemplateCache
}

func NewApp(tc *helpers.TemplateCache) *App {
	return &App{
		tc: tc,
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
}

type Properties struct {
	Title       string      `json:"title"`
	Abstract    string      `json:"abstract,omitempty"`
	Collection  string      `json:"collection"`
	GROMetadata GROMetadata `json:"gro_metadata"`
}

type GROMetadata struct {
	ImplementationStatus string `json:"implementation_status,omitempty"`
	DataFormat           string `json:"data_format,omitempty"`
	Country              string `json:"country,omitempty"`
	Continent            string `json:"continent,omitempty"`
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

package webui

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-spatial/geocatalogo/helpers"
	"github.com/go-spatial/geocatalogo/metadata"
)

func (a *App) HandleCatalog(w http.ResponseWriter, r *http.Request) {
	// Load all records from JSON
	catalogPath := os.Getenv("CATALOG_JSON_PATH")
	if catalogPath == "" {
		catalogPath = "/Users/jjohnson/projects/geosure/catalog/data/geocatalogo_records.json"
	}

	data, err := os.ReadFile(catalogPath)
	if err != nil {
		http.Error(w, "Failed to load catalog", http.StatusInternalServerError)
		log.Printf("Error loading catalog: %v", err)
		return
	}

	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		http.Error(w, "Failed to parse catalog", http.StatusInternalServerError)
		log.Printf("Error parsing catalog: %v", err)
		return
	}

	// Calculate stats
	stats := CatalogStats{Total: len(records)}
	continentCounts := make(map[string]int)
	countryCounts := make(map[string]int)

	for _, rec := range records {
		switch rec.Properties.Collection {
		case "existing_db":
			stats.ExistingDB++
		case "existing_local":
			stats.ExistingLocal++
		case "potential_v6":
			stats.PotentialV6++
		case "external_api":
			stats.ExternalAPI++
		case "external_news":
			stats.ExternalNews++
		case "external_government":
			stats.ExternalGov++
		case "external_other":
			stats.ExternalOther++
		}

		// Count continents and countries
		if rec.Properties.GROMetadata.Continent != "" {
			continentCounts[rec.Properties.GROMetadata.Continent]++
		}
		if rec.Properties.GROMetadata.Country != "" {
			countryCounts[rec.Properties.GROMetadata.Country]++
		}
	}

	// Convert maps to sorted slices
	for code, count := range continentCounts {
		stats.Continents = append(stats.Continents, ContinentStat{
			Code:  code,
			Name:  helpers.ContinentToName(code),
			Emoji: helpers.ContinentToEmoji(code),
			Count: count,
		})
	}

	for code, count := range countryCounts {
		stats.Countries = append(stats.Countries, CountryStat{
			Code:  code,
			Name:  helpers.CountryCodeToName(code),
			Flag:  helpers.CountryCodeToFlag(code),
			Count: count,
		})
	}

	pageData := PageData{
		Records: records,
		Stats:   stats,
	}

	if err := a.tc.Render(w, "layout_catalog", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleDataset(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/dataset/"):]

	// Load all records from JSON
	catalogPath := os.Getenv("CATALOG_JSON_PATH")
	if catalogPath == "" {
		catalogPath = "/Users/jjohnson/projects/geosure/catalog/data/geocatalogo_records.json"
	}

	data, err := os.ReadFile(catalogPath)
	if err != nil {
		http.Error(w, "Failed to load catalog", http.StatusInternalServerError)
		log.Printf("Error loading catalog: %v", err)
		return
	}

	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		http.Error(w, "Failed to parse catalog", http.StatusInternalServerError)
		log.Printf("Error parsing catalog: %v", err)
		return
	}

	// Find the record by ID
	var record *Record
	for i := range records {
		if records[i].ID == id {
			record = &records[i]
			break
		}
	}

	if record == nil {
		http.Error(w, "Dataset not found", http.StatusNotFound)
		return
	}

	// Build page data with introspection metadata
	pageData := DatasetPageData{
		Record: *record,
	}

	// Convert record to JSON for display
	recordJSON, err := json.MarshalIndent(record, "", "  ")
	if err == nil {
		pageData.RecordJSON = string(recordJSON)
	}

	// Lookup introspection data based on type
	if a.meta != nil {
		meta := a.meta.Lookup(
			record.ID,
			record.Properties.GROMetadata.S3Path,
			record.Properties.GROMetadata.DatabaseTable,
			record.Properties.GROMetadata.DataFormat,
		)

		// Type assert to specific metadata types
		switch m := meta.(type) {
		case *metadata.DatabaseTable:
			pageData.DatabaseTable = m
		case *metadata.CSVFile:
			pageData.CSVFile = m
		case *metadata.ParquetFile:
			pageData.ParquetFile = m
		case *metadata.ShapefileFile:
			pageData.ShapefileFile = m
		case *metadata.GeoPackageFile:
			pageData.GeoPackageFile = m
		case *metadata.ExcelFile:
			pageData.ExcelFile = m
		case *metadata.JSONFile:
			pageData.JSONFile = m
		case *metadata.FileGDBFile:
			pageData.FileGDBFile = m
		case *metadata.PNGFile:
			pageData.PNGFile = m
		case *metadata.PDFFile:
			pageData.PDFFile = m
		}
	}

	if err := a.tc.Render(w, "layout_dataset", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Render statistics page
	w.Write([]byte("Statistics page - coming soon!"))
}

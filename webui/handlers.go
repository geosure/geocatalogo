package webui

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	// TODO: Load record by ID and render detail page
	w.Write([]byte("Dataset detail: " + id))
}

func (a *App) HandleStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Render statistics page
	w.Write([]byte("Statistics page - coming soon!"))
}

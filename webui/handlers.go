package webui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

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

	// Calculate DataRecords (everything except jobs)
	stats.DataRecords = stats.Total - stats.PotentialV6

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
			record.Properties.GROMetadata.V6JobFile,
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
		case *metadata.V6Job:
			pageData.V6Job = m
		}
	}

	if err := a.tc.Render(w, "layout_dataset", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleGeography(w http.ResponseWriter, r *http.Request) {
	// Parse URL: /geography/{level}/{name}
	// Example: /geography/city/dallas
	path := r.URL.Path[len("/geography/"):]
	parts := []string{}
	if path != "" {
		// Split by / but preserve URL encoding
		for _, p := range r.URL.Query()["city"] {
			parts = append(parts, p)
		}
		for _, p := range r.URL.Query()["county"] {
			parts = append(parts, p)
		}
		for _, p := range r.URL.Query()["state"] {
			parts = append(parts, p)
		}
		for _, p := range r.URL.Query()["country"] {
			parts = append(parts, p)
		}
		for _, p := range r.URL.Query()["continent"] {
			parts = append(parts, p)
		}
	}

	// Extract geography params from query string
	city := r.URL.Query().Get("city")
	county := r.URL.Query().Get("county")
	state := r.URL.Query().Get("state")
	country := r.URL.Query().Get("country")
	continent := r.URL.Query().Get("continent")

	// Determine level and name
	level := ""
	name := ""
	if city != "" {
		level = "city"
		name = city
	} else if county != "" {
		level = "county"
		name = county
	} else if state != "" {
		level = "state"
		name = state
	} else if country != "" {
		level = "country"
		name = country
	} else if continent != "" {
		level = "continent"
		name = continent
	} else {
		http.Error(w, "No geography specified", http.StatusBadRequest)
		return
	}

	// Find matching README
	var readme *metadata.V6README
	if a.meta != nil {
		readme = a.meta.FindREADMEForGeography(city, county, state, country, continent)
	}

	// Load all records and filter by geography
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

	// Filter ALL records by geography (not just jobs!)
	var matchingRecords []Record
	var counts CollectionCounts

	for _, rec := range records {
		// Match geography
		matches := true
		if city != "" && rec.Properties.GROMetadata.City != city {
			matches = false
		}
		if county != "" && rec.Properties.GROMetadata.Admin2 != county {
			matches = false
		}
		if state != "" && rec.Properties.GROMetadata.StateProvince != state {
			matches = false
		}
		if country != "" && rec.Properties.GROMetadata.Country != country {
			matches = false
		}
		if continent != "" {
			// Special case: 'global' continent should match geographic_scope='global'
			// instead of continent field (which is empty for global datasets)
			if continent == "global" {
				if rec.Properties.GROMetadata.GeographicScope != "global" {
					matches = false
				}
			} else {
				if rec.Properties.GROMetadata.Continent != continent {
					matches = false
				}
			}
		}

		if matches {
			matchingRecords = append(matchingRecords, rec)

			// Count by collection type
			switch rec.Properties.Collection {
			case "potential_v6":
				counts.V6Jobs++
			case "existing_db":
				counts.Database++
			case "existing_local":
				counts.Files++
			case "external_api":
				counts.APIs++
			case "external_news":
				counts.News++
			case "external_government":
				counts.Government++
			case "ai_agent":
				counts.AIAgents++
			case "data_bot":
				counts.DataBots++
			case "scraper_bot":
				counts.ScraperBots++
			case "automation_bot":
				counts.AutoBots++
			case "verb_app":
				counts.VerbApps++
			case "internal_tool":
				counts.InternalTools++
			case "api_service":
				counts.APIServices++
			default:
				counts.Other++
			}
		}
	}

	// Extract sub-geographies and counts
	var subGeographies []SubGeography
	switch level {
	case "continent":
		// Extract unique countries
		countryCounts := make(map[string]int)
		for _, rec := range matchingRecords {
			if rec.Properties.GROMetadata.Country != "" {
				countryCounts[rec.Properties.GROMetadata.Country]++
			}
		}
		for code, count := range countryCounts {
			subGeographies = append(subGeographies, SubGeography{
				Code:  code,
				Name:  CountryName(code),
				Emoji: CountryEmoji(code),
				Count: count,
				URL:   fmt.Sprintf("/geography/?continent=%s&country=%s", continent, code),
			})
		}
		// Sort by count (descending)
		sort.Slice(subGeographies, func(i, j int) bool {
			return subGeographies[i].Count > subGeographies[j].Count
		})

	case "country":
		// Extract unique states
		stateCounts := make(map[string]int)
		for _, rec := range matchingRecords {
			if rec.Properties.GROMetadata.StateProvince != "" {
				stateCounts[rec.Properties.GROMetadata.StateProvince]++
			}
		}
		for code, count := range stateCounts {
			subGeographies = append(subGeographies, SubGeography{
				Code:  code,
				Name:  strings.ToUpper(code),
				Emoji: "📍",
				Count: count,
				URL:   fmt.Sprintf("/geography/?continent=%s&country=%s&state=%s", continent, country, code),
			})
		}
		// Sort by count (descending)
		sort.Slice(subGeographies, func(i, j int) bool {
			return subGeographies[i].Count > subGeographies[j].Count
		})

	case "state":
		// Extract unique cities
		cityCounts := make(map[string]int)
		for _, rec := range matchingRecords {
			if rec.Properties.GROMetadata.City != "" {
				cityCounts[rec.Properties.GROMetadata.City]++
			}
		}
		for code, count := range cityCounts {
			subGeographies = append(subGeographies, SubGeography{
				Code:  code,
				Name:  code,
				Emoji: "🏙️",
				Count: count,
				URL:   fmt.Sprintf("/geography/?continent=%s&country=%s&state=%s&city=%s", continent, country, state, code),
			})
		}
		// Sort by count (descending)
		sort.Slice(subGeographies, func(i, j int) bool {
			return subGeographies[i].Count > subGeographies[j].Count
		})
	}

	pageData := GeographyPageData{
		Level:            level,
		Name:             name,
		README:           readme,
		Jobs:             matchingRecords,
		JobCount:         len(matchingRecords),
		CollectionCounts: counts,
		SubGeographies:   subGeographies,
		City:             city,
		County:           county,
		State:            state,
		Country:          country,
		Continent:        continent,
	}

	if err := a.tc.Render(w, "layout_geography", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Render statistics page
	w.Write([]byte("Statistics page - coming soon!"))
}

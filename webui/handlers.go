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
	// Only handle exact "/" path - let other handlers handle their routes
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

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
	stats := CatalogStats{
		Total:        len(records),
		StatusCounts: make(map[string]int),
	}
	continentCounts := make(map[string]int)
	countryCounts := make(map[string]int)
	formatCounts := make(map[string]int)
	ownerCounts := make(map[string]int)

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

		// Count data formats
		if rec.Properties.GROMetadata.DataFormat != "" {
			formatCounts[rec.Properties.GROMetadata.DataFormat]++
		}

		// Count implementation statuses
		if rec.Properties.GROMetadata.ImplementationStatus != "" {
			stats.StatusCounts[rec.Properties.GROMetadata.ImplementationStatus]++
		}

		// Count owners
		if rec.Properties.GROMetadata.Owner != "" {
			ownerCounts[rec.Properties.GROMetadata.Owner]++
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

	for format, count := range formatCounts {
		stats.Formats = append(stats.Formats, FormatStat{
			Name:  format,
			Count: count,
		})
	}

	// Sort formats by count (descending)
	sort.Slice(stats.Formats, func(i, j int) bool {
		return stats.Formats[i].Count > stats.Formats[j].Count
	})

	for owner, count := range ownerCounts {
		stats.Owners = append(stats.Owners, OwnerStat{
			Name:  owner,
			Count: count,
		})
	}

	// Sort owners by count (descending)
	sort.Slice(stats.Owners, func(i, j int) bool {
		return stats.Owners[i].Count > stats.Owners[j].Count
	})

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
	// Extract ID from URL - supports both /dataset/id and /collection/id formats
	var id string
	path := r.URL.Path

	if strings.HasPrefix(path, "/dataset/") {
		id = path[len("/dataset/"):]
	} else {
		// Collection-based URL: /ai_agent/clankr_catalog_agent
		// Strip leading slash and split on first /
		parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
		if len(parts) == 2 {
			id = parts[1]
		}
	}

	if id == "" {
		http.Error(w, "Invalid dataset URL", http.StatusBadRequest)
		return
	}

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

	// DEBUG: Log what we loaded
	if record.Properties.Operational != nil {
		log.Printf("DEBUG: Operational data loaded for %s: %v", record.ID, record.Properties.Operational)
	} else {
		log.Printf("DEBUG: No operational data for %s", record.ID)
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
		case *metadata.Agent:
			pageData.Agent = m
		}
	}

	if err := a.tc.Render(w, "layout_dataset", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleGeography(w http.ResponseWriter, r *http.Request) {
	// Parse URL: /geography/{level}/{code}
	// Examples: /geography/country/mx, /geography/continent/north_america
	path := r.URL.Path[len("/geography/"):]

	var level, name string
	var city, county, state, country, continent string

	// Try to parse path-based format first: /geography/country/mx
	if path != "" {
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			level = parts[0]
			name = parts[1]

			// Set the appropriate variable based on level
			switch level {
			case "city":
				city = name
			case "county":
				county = name
			case "state":
				state = name
			case "country":
				country = name
			case "continent":
				continent = name
			default:
				http.Error(w, "Invalid geography level", http.StatusBadRequest)
				return
			}
		}
	}

	// Fall back to query string format for backward compatibility: ?country=mx
	if level == "" {
		city = r.URL.Query().Get("city")
		county = r.URL.Query().Get("county")
		state = r.URL.Query().Get("state")
		country = r.URL.Query().Get("country")
		continent = r.URL.Query().Get("continent")

		// Determine level and name from query params
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
			// Match by continent field (includes 'global' for database tables and local files)
			if rec.Properties.GROMetadata.Continent != continent {
				matches = false
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
			case "claude_projects":
				counts.ClaudeProjects++
			case "operational_service":
				counts.OperationalServices++
			case "data_inspection_bot":
				counts.DataInspectionBots++
			case "catalog_management_bot":
				counts.CatalogManagementBots++
			case "data_bot":
				counts.DataBots++
			case "scraper_bot":
				counts.ScraperBots++
			case "automation_bot":
				counts.AutoBots++
			case "historical_agent":
				counts.HistoricalAgents++
			case "verb_app":
				counts.VerbApps++
			case "team_member":
				counts.TeamMembers++
			case "infrastructure":
				counts.Infrastructure++
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

func (a *App) HandleFormat(w http.ResponseWriter, r *http.Request) {
	// Parse URL: /format/{format_name}
	// Example: /format/csv, /format/api
	formatName := strings.TrimPrefix(r.URL.Path, "/format/")

	if formatName == "" {
		http.Error(w, "No format specified", http.StatusBadRequest)
		return
	}

	// Load all records and filter by format
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

	// Filter records by format
	var matchingRecords []Record
	var counts CollectionCounts

	for _, rec := range records {
		if rec.Properties.GROMetadata.DataFormat == formatName {
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
			case "data_inspection_bot":
				counts.DataInspectionBots++
			case "scraper_bot":
				counts.ScraperBots++
			case "automation_bot":
				counts.AutoBots++
			}
		}
	}

	pageData := GeographyPageData{
		Level:            "format",
		Name:             formatName,
		Jobs:             matchingRecords,
		JobCount:         len(matchingRecords),
		CollectionCounts: counts,
	}

	if err := a.tc.Render(w, "layout_geography", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleStatus(w http.ResponseWriter, r *http.Request) {
	// Parse URL: /status/{status_name}
	// Example: /status/implemented, /status/draft
	statusName := strings.TrimPrefix(r.URL.Path, "/status/")

	if statusName == "" {
		http.Error(w, "No status specified", http.StatusBadRequest)
		return
	}

	// Load all records and filter by implementation status
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

	// Filter records by implementation status
	var matchingRecords []Record
	var counts CollectionCounts

	for _, rec := range records {
		if rec.Properties.GROMetadata.ImplementationStatus == statusName {
			matchingRecords = append(matchingRecords, rec)

			// Count by collection type (same logic as HandleFormat)
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
			case "external_download":
				counts.News++
			case "ai_agent":
				counts.AIAgents++
			case "claude_projects":
				counts.ClaudeProjects++
			case "operational_service":
				counts.OperationalServices++
			case "data_inspection_bot":
				counts.DataInspectionBots++
			case "catalog_management_bot":
				counts.CatalogManagementBots++
			case "data_bot":
				counts.DataBots++
			case "scraper_bot":
				counts.ScraperBots++
			case "automation_bot":
				counts.AutoBots++
			case "historical_agent":
				counts.HistoricalAgents++
			case "verb_app":
				counts.VerbApps++
			case "internal_tool":
				counts.InternalTools++
			case "api_service":
				counts.APIServices++
			case "team_member":
				counts.TeamMembers++
			case "infrastructure":
				counts.Infrastructure++
			default:
				counts.Other++
			}
		}
	}

	pageData := GeographyPageData{
		Level:            "status",
		Name:             statusName,
		Jobs:             matchingRecords,
		JobCount:         len(matchingRecords),
		CollectionCounts: counts,
	}

	if err := a.tc.Render(w, "layout_geography", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleOwner(w http.ResponseWriter, r *http.Request) {
	// Parse URL: /owner/{owner_name}
	// Example: /owner/@data-platform, /owner/@clankr
	ownerName := strings.TrimPrefix(r.URL.Path, "/owner/")

	if ownerName == "" {
		http.Error(w, "No owner specified", http.StatusBadRequest)
		return
	}

	// Load all records and filter by owner
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

	// Filter records by owner
	var matchingRecords []Record
	var counts CollectionCounts

	for _, rec := range records {
		if rec.Properties.GROMetadata.Owner == ownerName {
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
			case "external_download":
				counts.News++
			case "ai_agent":
				counts.AIAgents++
			case "claude_projects":
				counts.ClaudeProjects++
			case "operational_service":
				counts.OperationalServices++
			case "data_inspection_bot":
				counts.DataInspectionBots++
			case "catalog_management_bot":
				counts.CatalogManagementBots++
			case "data_bot":
				counts.DataBots++
			case "scraper_bot":
				counts.ScraperBots++
			case "automation_bot":
				counts.AutoBots++
			case "historical_agent":
				counts.HistoricalAgents++
			case "verb_app":
				counts.VerbApps++
			case "internal_tool":
				counts.InternalTools++
			case "api_service":
				counts.APIServices++
			case "team_member":
				counts.TeamMembers++
			case "infrastructure":
				counts.Infrastructure++
			default:
				counts.Other++
			}
		}
	}

	pageData := GeographyPageData{
		Level:            "owner",
		Name:             ownerName,
		Jobs:             matchingRecords,
		JobCount:         len(matchingRecords),
		CollectionCounts: counts,
	}

	if err := a.tc.Render(w, "layout_geography", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleStats(w http.ResponseWriter, r *http.Request) {
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

	// Build comprehensive statistics
	stats := buildDetailedStats(records, a.meta)

	if err := a.tc.Render(w, "layout_stats", stats); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func buildDetailedStats(records []Record, meta *metadata.MetadataStore) StatsPageData {
	stats := StatsPageData{
		TotalRecords: len(records),
	}

	// Count by collection type
	collectionCounts := make(map[string]int)
	agentCount := 0
	botCount := 0
	serviceCount := 0
	dbTableCount := 0
	fileCount := 0

	for _, rec := range records {
		collectionCounts[rec.Properties.Collection]++

		// Count specific types
		switch rec.Properties.Collection {
		case "ai_agent":
			agentCount++
		case "operational_service":
			serviceCount++
		case "data_inspection_bot", "catalog_management_bot", "automation_bot", "scraper_bot", "data_bot":
			botCount++
		case "existing_db":
			dbTableCount++
		case "existing_local":
			fileCount++
		}
	}

	stats.AgentCount = agentCount
	stats.BotCount = botCount
	stats.ServiceCount = serviceCount
	stats.DatabaseTableCount = dbTableCount
	stats.FileCount = fileCount
	stats.V6JobCount = collectionCounts["potential_v6"]
	stats.ExternalSourceCount = collectionCounts["external_api"] + collectionCounts["external_news"] + collectionCounts["external_government"] + collectionCounts["external_other"]

	// Get metadata counts from introspection
	if meta != nil {
		if meta.Database != nil {
			stats.IntrospectionStats.DatabaseTables = len(meta.Database.Tables)
		}
		stats.IntrospectionStats.Agents = len(meta.Agents)
		stats.IntrospectionStats.CSVFiles = len(meta.CSVFiles)
		stats.IntrospectionStats.ParquetFiles = len(meta.Parquet)
		stats.IntrospectionStats.Shapefiles = len(meta.Shapefile)
		stats.IntrospectionStats.PDFs = len(meta.PDF)
		stats.IntrospectionStats.V6Jobs = len(meta.V6Jobs)
	}

	return stats
}

func (a *App) HandleAPIDocs(w http.ResponseWriter, r *http.Request) {
	if err := a.tc.Render(w, "layout_api_docs", nil); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleCollections(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleCollections called for path: %s", r.URL.Path)
	if err := a.tc.Render(w, "layout_collections", nil); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("❌ Collections template error: %v", err)
	} else {
		log.Printf("✅ Collections page rendered successfully")
	}
}

func (a *App) HandleOwners(w http.ResponseWriter, r *http.Request) {
	// Load all records to count owners
	catalogPath := os.Getenv("CATALOG_JSON_PATH")
	if catalogPath == "" {
		catalogPath = "/Users/jjohnson/projects/geosure/catalog/data/geocatalogo_records.json"
	}

	data, err := os.ReadFile(catalogPath)
	if err != nil {
		http.Error(w, "Failed to load catalog", http.StatusInternalServerError)
		return
	}

	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		http.Error(w, "Failed to parse catalog", http.StatusInternalServerError)
		return
	}

	// Count owners
	ownerCounts := make(map[string]int)
	for _, rec := range records {
		if rec.Properties.GROMetadata.Owner != "" {
			ownerCounts[rec.Properties.GROMetadata.Owner]++
		}
	}

	// Convert to sorted slice
	var owners []OwnerStat
	for owner, count := range ownerCounts {
		owners = append(owners, OwnerStat{
			Name:  owner,
			Count: count,
		})
	}

	// Sort by count descending
	sort.Slice(owners, func(i, j int) bool {
		return owners[i].Count > owners[j].Count
	})

	pageData := struct {
		Owners     []OwnerStat
		TotalOwners int
		TotalRecords int
	}{
		Owners:      owners,
		TotalOwners: len(owners),
		TotalRecords: len(records),
	}

	if err := a.tc.Render(w, "layout_owners", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleStatuses(w http.ResponseWriter, r *http.Request) {
	// Load all records to count statuses
	catalogPath := os.Getenv("CATALOG_JSON_PATH")
	if catalogPath == "" {
		catalogPath = "/Users/jjohnson/projects/geosure/catalog/data/geocatalogo_records.json"
	}

	data, err := os.ReadFile(catalogPath)
	if err != nil {
		http.Error(w, "Failed to load catalog", http.StatusInternalServerError)
		return
	}

	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		http.Error(w, "Failed to parse catalog", http.StatusInternalServerError)
		return
	}

	// Count statuses
	statusCounts := make(map[string]int)
	for _, rec := range records {
		if rec.Properties.GROMetadata.ImplementationStatus != "" {
			statusCounts[rec.Properties.GROMetadata.ImplementationStatus]++
		}
	}

	// Convert to sorted slice with proper ordering
	statusOrder := []string{"implemented", "draft", "potential", "active", "archived"}
	var statuses []StatusStat

	// Add statuses in the defined order
	for _, status := range statusOrder {
		if count, ok := statusCounts[status]; ok {
			statuses = append(statuses, StatusStat{
				Name:  status,
				Count: count,
			})
			delete(statusCounts, status)
		}
	}

	// Add any remaining statuses not in the predefined order
	for status, count := range statusCounts {
		statuses = append(statuses, StatusStat{
			Name:  status,
			Count: count,
		})
	}

	pageData := struct {
		Statuses      []StatusStat
		TotalStatuses int
		TotalRecords  int
	}{
		Statuses:      statuses,
		TotalStatuses: len(statuses),
		TotalRecords:  len(records),
	}

	if err := a.tc.Render(w, "layout_statuses", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleGeographies(w http.ResponseWriter, r *http.Request) {
	// Load all records to count continents
	catalogPath := os.Getenv("CATALOG_JSON_PATH")
	if catalogPath == "" {
		catalogPath = "/Users/jjohnson/projects/geosure/catalog/data/geocatalogo_records.json"
	}

	data, err := os.ReadFile(catalogPath)
	if err != nil {
		http.Error(w, "Failed to load catalog", http.StatusInternalServerError)
		return
	}

	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		http.Error(w, "Failed to parse catalog", http.StatusInternalServerError)
		return
	}

	// Count continents
	continentCounts := make(map[string]int)
	for _, rec := range records {
		if rec.Properties.GROMetadata.Continent != "" {
			continentCounts[rec.Properties.GROMetadata.Continent]++
		}
	}

	// Convert to sorted slice
	var continents []ContinentStat
	for code, count := range continentCounts {
		continents = append(continents, ContinentStat{
			Code:  code,
			Name:  helpers.ContinentToName(code),
			Emoji: helpers.ContinentToEmoji(code),
			Count: count,
		})
	}

	// Sort by count descending
	sort.Slice(continents, func(i, j int) bool {
		return continents[i].Count > continents[j].Count
	})

	pageData := struct {
		Continents      []ContinentStat
		TotalContinents int
		TotalRecords    int
	}{
		Continents:      continents,
		TotalContinents: len(continents),
		TotalRecords:    len(records),
	}

	if err := a.tc.Render(w, "layout_geographies", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (a *App) HandleCollectionDetail(w http.ResponseWriter, r *http.Request) {
	// Extract collection name from URL: /collection/{name}
	path := r.URL.Path
	collectionName := path[len("/collection/"):]
	if collectionName == "" {
		http.Redirect(w, r, "/collections", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("HandleCollectionDetail called for collection: %s", collectionName)

	// Load all records
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

	var allRecords []Record
	if err := json.Unmarshal(data, &allRecords); err != nil {
		http.Error(w, "Failed to parse catalog", http.StatusInternalServerError)
		log.Printf("Error parsing catalog: %v", err)
		return
	}

	// Filter records by collection
	var records []Record
	for _, rec := range allRecords {
		if rec.Properties.Collection == collectionName {
			records = append(records, rec)
		}
	}

	// Group by implementation status OR by category for infrastructure OR build hierarchy for v6
	byStatus := make(map[string][]Record)
	byCategory := make(map[string][]Record)
	var hierarchy *TreeNode

	if collectionName == "infrastructure" {
		// For infrastructure, group by data_format (which contains the actual AWS service type)
		for _, rec := range records {
			category := rec.Properties.GROMetadata.DataFormat
			if category == "" {
				category = "uncategorized"
			}
			byCategory[category] = append(byCategory[category], rec)
		}
	} else if collectionName == "potential_v6" {
		// For v6 jobs, build hierarchical tree from v6_job_file paths
		hierarchy = buildV6Hierarchy(records)
	} else {
		// For other collections, group by implementation status
		for _, rec := range records {
			status := rec.Properties.GROMetadata.ImplementationStatus
			if status == "" {
				status = "unspecified"
			}
			byStatus[status] = append(byStatus[status], rec)
		}
	}

	// Collection metadata
	collectionMetadata := map[string]struct {
		Name        string
		Emoji       string
		Description string
	}{
		"ai_agent":                {"AI Agents", "🤖", "Autonomous AI agents managing different parts of the GRO ecosystem"},
		"existing_db":             {"Database Tables", "🗄️", "PostgreSQL + PostGIS tables in the production database"},
		"existing_local":          {"Local Files", "📁", "Data files stored locally (CSV, Parquet, Shapefile, GeoJSON, etc.)"},
		"potential_v6":            {"v6 Job Definitions", "⚙️", "YAML job definitions for the future automation pipeline"},
		"external_api":            {"External APIs", "🔌", "Third-party API endpoints we integrate with"},
		"external_government":     {"Government Data", "🏛️", "Government open data portals and official statistics"},
		"external_news":           {"News & Media", "📰", "News APIs, RSS feeds, and media monitoring sources"},
		"external_news_active":    {"Active News Sources", "📡", "News sources actively collected via GNews API"},
		"external_academic":       {"Academic Data", "🎓", "Academic datasets and research institution data"},
		"external_download":       {"Downloadable Datasets", "⬇️", "Datasets available for direct download"},
		"external_other":          {"Other External Sources", "🌐", "Other external data sources"},
		"infrastructure":          {"Infrastructure", "⚙️", "AWS infrastructure components (Lambda, RDS, S3, etc.)"},
		"internal_tool":           {"Internal Tools", "🛠️", "Internal tools and utilities for development"},
		"verb_app":                {"Verb Applications", "🔤", "Verb-based applications (explore, curate, chronicle, etc.)"},
		"team_member":             {"Team Members", "👥", "Team members and their roles"},
		"claude_projects":         {"Claude Projects", "💬", "Claude Projects used in development"},
		"historical_agent":        {"Historical Agents", "📦", "Agents that are no longer active but kept for reference"},
		"automation_bot":          {"Automation Bots", "🔧", "Bots that automate repetitive tasks and workflows"},
		"data_inspection_bot":     {"Data Inspection Bots", "🔍", "Bots that introspect and validate data quality"},
		"catalog_management_bot":  {"Catalog Management Bots", "📋", "Bots that maintain and update the catalog"},
		"operational_service":     {"Operational Services", "⚡", "Running services powering the GRO platform"},
		"api_service":             {"API Services", "🔌", "API endpoints for data access and integration"},
	}

	meta := collectionMetadata[collectionName]
	if meta.Name == "" {
		meta.Name = collectionName
		meta.Emoji = "📦"
		meta.Description = "Resources in this collection"
	}

	pageData := CollectionDetailPageData{
		CollectionCode:        collectionName,
		CollectionName:        meta.Name,
		CollectionEmoji:       meta.Emoji,
		CollectionDescription: meta.Description,
		TotalCount:            len(records),
		Records:               records,
		ByStatus:              byStatus,
		ByCategory:            byCategory,
		Hierarchy:             hierarchy,
	}

	if err := a.tc.Render(w, "layout_collection_detail", pageData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("❌ Collection detail template error: %v", err)
	}
}


func (a *App) HandleQuery(w http.ResponseWriter, r *http.Request) {
	if err := a.tc.Render(w, "layout_query", nil); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("❌ Query page template error: %v", err)
	}
}

// buildV6Hierarchy constructs a hierarchical tree from v6 job file paths
// Example path: "v6/jobs/africa/ml/bamako/bootstrap/040_import_bamako_air_quality.yml"
// This creates a tree structure like:
//   jobs/
//     ├── africa/
//     │   └── ml/
//     │       └── bamako/
//     │           └── bootstrap/
//     │               └── 040_import_bamako_air_quality.yml
func buildV6Hierarchy(records []Record) *TreeNode {
	root := &TreeNode{
		Name:     "jobs",
		Path:     "jobs",
		Level:    0,
		Children: make(map[string]*TreeNode),
	}

	for i := range records {
		rec := &records[i]
		jobFile := rec.Properties.GROMetadata.V6JobFile
		if jobFile == "" {
			continue
		}

		// Parse path: "v6/jobs/africa/ml/bamako/bootstrap/040_import_bamako_air_quality.yml"
		// Split and skip "v6/" prefix, start from "jobs"
		parts := strings.Split(jobFile, "/")
		if len(parts) < 2 || parts[0] != "v6" || parts[1] != "jobs" {
			continue
		}

		// Start from "jobs" (parts[1])
		current := root

		// Traverse/create tree structure for all parts after "jobs"
		for j := 2; j < len(parts); j++ {
			part := parts[j]
			isLeaf := (j == len(parts)-1) // Last part is the filename

			if current.Children == nil {
				current.Children = make(map[string]*TreeNode)
			}

			node, exists := current.Children[part]
			if !exists {
				node = &TreeNode{
					Name:     part,
					Path:     strings.Join(parts[1:j+1], "/"), // Full path from "jobs"
					Level:    j - 1,
					IsLeaf:   isLeaf,
					Children: make(map[string]*TreeNode),
				}
				current.Children[part] = node
			}

			if isLeaf {
				node.Record = rec
				node.Count = 1
			}

			current = node
		}
	}

	// Calculate counts for all non-leaf nodes
	calculateCounts(root)

	// Convert Children maps to sorted slices for stable template rendering
	sortChildrenRecursive(root)

	return root
}

// calculateCounts recursively calculates the total count of leaf nodes under each node
func calculateCounts(node *TreeNode) int {
	if node.IsLeaf {
		return 1
	}

	totalCount := 0
	for _, child := range node.Children {
		totalCount += calculateCounts(child)
	}
	node.Count = totalCount
	return totalCount
}

// sortChildrenRecursive converts Children maps to sorted slices recursively
func sortChildrenRecursive(node *TreeNode) {
	if node.IsLeaf || len(node.Children) == 0 {
		return
	}

	// Convert map to slice
	children := make([]*TreeNode, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, child)
	}

	// Sort by name
	sort.Slice(children, func(i, j int) bool {
		// Directories first, then files
		if !children[i].IsLeaf && children[j].IsLeaf {
			return true
		}
		if children[i].IsLeaf && !children[j].IsLeaf {
			return false
		}
		return children[i].Name < children[j].Name
	})

	node.ChildrenSorted = children

	// Recursively sort children
	for _, child := range children {
		sortChildrenRecursive(child)
	}
}

package webui

import "net/http"

func NewMux(app *App) *http.ServeMux {
	mux := http.NewServeMux()

	// Register specific routes BEFORE catch-all "/" route
	mux.HandleFunc("/dataset/", app.HandleDataset) // Legacy URL support
	mux.HandleFunc("/geography/", app.HandleGeography)
	mux.HandleFunc("/format/", app.HandleFormat)
	mux.HandleFunc("/status/", app.HandleStatus)
	mux.HandleFunc("/owner/", app.HandleOwner)
	mux.HandleFunc("/collection/", app.HandleCollectionDetail)
	mux.HandleFunc("GET /collections", app.HandleCollections)   // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /statuses", app.HandleStatuses)         // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /geographies", app.HandleGeographies)   // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /owners", app.HandleOwners)             // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /query", app.HandleQuery)               // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /stats", app.HandleStats)               // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /api-docs", app.HandleAPIDocs)          // Go 1.22+ exact match syntax

	// Collection-based dataset URLs (e.g., /ai_agent/clankr_catalog_agent)
	mux.HandleFunc("/ai_agent/", app.HandleDataset)
	mux.HandleFunc("/data_inspection_bot/", app.HandleDataset)
	mux.HandleFunc("/scraper_bot/", app.HandleDataset)
	mux.HandleFunc("/automation_bot/", app.HandleDataset)
	mux.HandleFunc("/data_bot/", app.HandleDataset)
	mux.HandleFunc("/catalog_management_bot/", app.HandleDataset)
	mux.HandleFunc("/operational_service/", app.HandleDataset)
	mux.HandleFunc("/claude_projects/", app.HandleDataset)
	mux.HandleFunc("/historical_agent/", app.HandleDataset)
	mux.HandleFunc("/verb_app/", app.HandleDataset)
	mux.HandleFunc("/team_member/", app.HandleDataset)
	mux.HandleFunc("/infrastructure/", app.HandleDataset)
	mux.HandleFunc("/internal_tool/", app.HandleDataset)
	mux.HandleFunc("/api_service/", app.HandleDataset)
	mux.HandleFunc("/potential_v6/", app.HandleDataset)
	mux.HandleFunc("/existing_db/", app.HandleDataset)
	mux.HandleFunc("/existing_local/", app.HandleDataset)
	mux.HandleFunc("/external_api/", app.HandleDataset)
	mux.HandleFunc("/external_news/", app.HandleDataset)
	mux.HandleFunc("/external_government/", app.HandleDataset)
	mux.HandleFunc("/external_academic/", app.HandleDataset)
	mux.HandleFunc("/external_download/", app.HandleDataset)
	mux.HandleFunc("/external_other/", app.HandleDataset)

	// Redirect /api to API documentation page
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api-docs", http.StatusTemporaryRedirect)
	})

	// Register "/" last as it's a catch-all
	mux.HandleFunc("/", app.HandleCatalog)

	return mux
}

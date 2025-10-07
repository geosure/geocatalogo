package webui

import "net/http"

func NewMux(app *App) *http.ServeMux {
	mux := http.NewServeMux()

	// Register specific routes BEFORE catch-all "/" route
	mux.HandleFunc("/dataset/", app.HandleDataset)
	mux.HandleFunc("/geography/", app.HandleGeography)
	mux.HandleFunc("/collection/", app.HandleCollectionDetail)
	mux.HandleFunc("GET /collections", app.HandleCollections) // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /query", app.HandleQuery)             // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /stats", app.HandleStats)             // Go 1.22+ exact match syntax
	mux.HandleFunc("GET /api-docs", app.HandleAPIDocs)        // Go 1.22+ exact match syntax

	// Redirect /api to existing GeoCatalogo STAC API
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8000/api", http.StatusTemporaryRedirect)
	})

	// Register "/" last as it's a catch-all
	mux.HandleFunc("/", app.HandleCatalog)

	return mux
}

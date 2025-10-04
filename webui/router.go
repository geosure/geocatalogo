package webui

import "net/http"

func NewMux(app *App) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.HandleCatalog)
	mux.HandleFunc("/dataset/", app.HandleDataset)
	mux.HandleFunc("/stats", app.HandleStats)

	// Redirect /api to existing GeoCatalogo STAC API
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8000/api", http.StatusTemporaryRedirect)
	})

	return mux
}

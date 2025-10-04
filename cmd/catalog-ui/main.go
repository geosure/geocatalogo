package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-spatial/geocatalogo/helpers"
	"github.com/go-spatial/geocatalogo/webui"
)

func main() {
	// Initialize template cache
	fsys := os.DirFS("/Users/jjohnson/projects/geocatalogo")
	tc := helpers.NewTemplateCache(fsys, helpers.FuncMap)

	// Create app and router
	app := webui.NewApp(tc)
	mux := webui.NewMux(app)

	// Start server
	port := ":3000"
	log.Printf("🚀 GRO Catalog UI starting on http://localhost%s", port)
	log.Printf("📊 Catalog browser at http://localhost%s/", port)
	log.Printf("📈 Statistics at http://localhost%s/stats", port)
	log.Printf("🔌 API docs at http://localhost%s/api", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

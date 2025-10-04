package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-spatial/geocatalogo/helpers"
	"github.com/go-spatial/geocatalogo/metadata"
	"github.com/go-spatial/geocatalogo/webui"
)

func main() {
	// Initialize template cache
	fsys := os.DirFS("/Users/jjohnson/projects/geocatalogo")
	tc := helpers.NewTemplateCache(fsys, helpers.FuncMap)

	// Load all introspection metadata into memory
	catalogDataPath := os.Getenv("CATALOG_DATA_PATH")
	if catalogDataPath == "" {
		catalogDataPath = "/Users/jjohnson/projects/geosure/catalog/data"
	}

	log.Printf("📦 Loading introspection metadata from %s...", catalogDataPath)
	meta, err := metadata.LoadAll(catalogDataPath)
	if err != nil {
		log.Printf("⚠️  Warning: Could not load metadata: %v", err)
	} else {
		log.Printf("✅ Metadata loaded successfully")
		if meta.Database != nil {
			log.Printf("   - Database: %d tables", len(meta.Database.Tables))
		}
		log.Printf("   - CSV/TSV: %d files", len(meta.CSVFiles))
		log.Printf("   - Parquet: %d files", len(meta.Parquet))
		log.Printf("   - Shapefile: %d files", len(meta.Shapefile))
		log.Printf("   - GeoPackage: %d files", len(meta.GeoPackage))
		log.Printf("   - Excel: %d files", len(meta.Excel))
		log.Printf("   - JSON: %d files", len(meta.JSON))
		log.Printf("   - File GDB: %d files", len(meta.FileGDB))
		log.Printf("   - PNG: %d files", len(meta.PNG))
		log.Printf("   - PDF: %d files", len(meta.PDF))
		log.Printf("   - V6 Jobs: %d jobs", len(meta.V6Jobs))
		log.Printf("   - V6 READMEs: %d files", len(meta.V6READMEs))
	}

	// Create app and router
	app := webui.NewApp(tc, meta)
	mux := webui.NewMux(app)

	// Start server
	port := ":3000"
	log.Printf("🚀 GRO Catalog UI starting on http://localhost%s", port)
	log.Printf("📊 Catalog browser at http://localhost%s/", port)
	log.Printf("📈 Statistics at http://localhost%s/stats", port)
	log.Printf("🔌 API docs at http://localhost%s/api", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

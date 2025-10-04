package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// MetadataStore holds all introspection data in memory
type MetadataStore struct {
	Database   *DatabaseSchema
	CSVFiles   map[string]*CSVFile
	Parquet    map[string]*ParquetFile
	Shapefile  map[string]*ShapefileFile
	GeoPackage map[string]*GeoPackageFile
	Excel      map[string]*ExcelFile
	JSON       map[string]*JSONFile
	FileGDB    map[string]*FileGDBFile
	PNG        map[string]*PNGFile
	PDF        map[string]*PDFFile
	V6Jobs     map[string]*V6Job
	V6READMEs  []V6README
}

// LoadAll loads all introspection files into memory
func LoadAll(basePath string) (*MetadataStore, error) {
	store := &MetadataStore{
		CSVFiles:   make(map[string]*CSVFile),
		Parquet:    make(map[string]*ParquetFile),
		Shapefile:  make(map[string]*ShapefileFile),
		GeoPackage: make(map[string]*GeoPackageFile),
		Excel:      make(map[string]*ExcelFile),
		JSON:       make(map[string]*JSONFile),
		FileGDB:    make(map[string]*FileGDBFile),
		PNG:        make(map[string]*PNGFile),
		PDF:        make(map[string]*PDFFile),
		V6Jobs:     make(map[string]*V6Job),
	}

	// Load database schema
	if err := loadDatabaseSchema(basePath+"/database_schema_latest.json", store); err != nil {
		fmt.Printf("Warning: Could not load database schema: %v\n", err)
	}

	// Load CSV introspection
	if err := loadCSVIntrospection(basePath+"/csv_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load CSV introspection: %v\n", err)
	}

	// Load Parquet introspection
	if err := loadParquetIntrospection(basePath+"/parquet_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load Parquet introspection: %v\n", err)
	}

	// Load Shapefile introspection
	if err := loadShapefileIntrospection(basePath+"/shapefile_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load Shapefile introspection: %v\n", err)
	}

	// Load GeoPackage introspection
	if err := loadGeoPackageIntrospection(basePath+"/geopackage_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load GeoPackage introspection: %v\n", err)
	}

	// Load Excel introspection
	if err := loadExcelIntrospection(basePath+"/excel_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load Excel introspection: %v\n", err)
	}

	// Load JSON introspection
	if err := loadJSONIntrospection(basePath+"/json_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load JSON introspection: %v\n", err)
	}

	// Load File GDB introspection
	if err := loadFileGDBIntrospection(basePath+"/filegdb_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load File GDB introspection: %v\n", err)
	}

	// Load PNG introspection
	if err := loadPNGIntrospection(basePath+"/png_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load PNG introspection: %v\n", err)
	}

	// Load PDF introspection
	if err := loadPDFIntrospection(basePath+"/pdf_introspection.json", store); err != nil {
		fmt.Printf("Warning: Could not load PDF introspection: %v\n", err)
	}

	// Load v6 job metadata
	if err := loadV6JobMetadata(basePath+"/v6_job_metadata.json", store); err != nil {
		fmt.Printf("Warning: Could not load v6 job metadata: %v\n", err)
	}

	// Load v6 README metadata
	if err := loadV6READMEMetadata(basePath+"/v6_readme_metadata.json", store); err != nil {
		fmt.Printf("Warning: Could not load v6 README metadata: %v\n", err)
	}

	return store, nil
}

// Lookup finds metadata for a given dataset by catalog ID or path
func (s *MetadataStore) Lookup(catalogID string, s3Path string, databaseTable string, dataFormat string, v6JobFile string) interface{} {
	// V6 job lookup
	if v6JobFile != "" {
		if job, ok := s.V6Jobs[v6JobFile]; ok {
			return job
		}
	}

	// Database table lookup
	if databaseTable != "" && s.Database != nil {
		for i := range s.Database.Tables {
			if s.Database.Tables[i].Name == databaseTable {
				return &s.Database.Tables[i]
			}
		}
	}

	// File lookup by S3 path
	if s3Path != "" {
		// Extract S3 key from full path
		key := strings.TrimPrefix(s3Path, "s3://geosure-data-dev/")

		// Try each format
		switch strings.ToLower(dataFormat) {
		case "csv", "tsv", "txt":
			if meta, ok := s.CSVFiles[key]; ok {
				return meta
			}
		case "parquet":
			if meta, ok := s.Parquet[key]; ok {
				return meta
			}
		case "shapefile", "shp":
			if meta, ok := s.Shapefile[key]; ok {
				return meta
			}
		case "geopackage", "gpkg":
			if meta, ok := s.GeoPackage[key]; ok {
				return meta
			}
		case "excel", "xlsx", "xls":
			if meta, ok := s.Excel[key]; ok {
				return meta
			}
		case "json", "geojson":
			if meta, ok := s.JSON[key]; ok {
				return meta
			}
		case "filegdb", "gdb":
			if meta, ok := s.FileGDB[key]; ok {
				return meta
			}
		case "png", "image":
			if meta, ok := s.PNG[key]; ok {
				return meta
			}
		case "pdf":
			if meta, ok := s.PDF[key]; ok {
				return meta
			}
		}
	}

	return nil
}

func loadDatabaseSchema(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var schema DatabaseSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return err
	}

	store.Database = &schema
	return nil
}

func loadCSVIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		Files []CSVFile `json:"files"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.Files {
		store.CSVFiles[result.Files[i].S3Key] = &result.Files[i]
	}
	return nil
}

func loadParquetIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		Files []ParquetFile `json:"files"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.Files {
		store.Parquet[result.Files[i].S3Key] = &result.Files[i]
	}
	return nil
}

func loadShapefileIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		Shapefiles []ShapefileFile `json:"shapefiles"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.Shapefiles {
		store.Shapefile[result.Shapefiles[i].S3Key] = &result.Shapefiles[i]
	}
	return nil
}

func loadGeoPackageIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		GeoPackages []GeoPackageFile `json:"geopackages"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.GeoPackages {
		store.GeoPackage[result.GeoPackages[i].S3Key] = &result.GeoPackages[i]
	}
	return nil
}

func loadExcelIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		Files []ExcelFile `json:"files"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.Files {
		store.Excel[result.Files[i].S3Key] = &result.Files[i]
	}
	return nil
}

func loadJSONIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		Files []JSONFile `json:"files"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.Files {
		store.JSON[result.Files[i].S3Key] = &result.Files[i]
	}
	return nil
}

func loadFileGDBIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		FileGDBs []FileGDBFile `json:"filegdbs"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.FileGDBs {
		store.FileGDB[result.FileGDBs[i].S3Key] = &result.FileGDBs[i]
	}
	return nil
}

func loadPNGIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		Images []PNGFile `json:"images"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.Images {
		store.PNG[result.Images[i].S3Key] = &result.Images[i]
	}
	return nil
}

func loadPDFIntrospection(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result struct {
		PDFs []PDFFile `json:"pdfs"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	for i := range result.PDFs {
		store.PDF[result.PDFs[i].S3Key] = &result.PDFs[i]
	}
	return nil
}
func loadV6JobMetadata(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// v6_job_metadata.json is a map: path -> V6Job
	var jobsMap map[string]*V6Job
	if err := json.Unmarshal(data, &jobsMap); err != nil {
		return err
	}

	store.V6Jobs = jobsMap
	return nil
}

func loadV6READMEMetadata(path string, store *MetadataStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var readmes []V6README
	if err := json.Unmarshal(data, &readmes); err != nil {
		return err
	}

	store.V6READMEs = readmes
	return nil
}

// FindREADMEForGeography finds README matching geographic hierarchy
func (s *MetadataStore) FindREADMEForGeography(city, county, state, country, continent string) *V6README {
	// Try city level first (most specific)
	if city != "" {
		for i := range s.V6READMEs {
			geo := s.V6READMEs[i].Geography
			if geoCity, ok := geo["city"].(string); ok && geoCity == city {
				if level, ok := geo["level"].(string); ok && level == "city" {
					return &s.V6READMEs[i]
				}
			}
		}
	}

	// Try county
	if county != "" {
		for i := range s.V6READMEs {
			geo := s.V6READMEs[i].Geography
			if geoCounty, ok := geo["county"].(string); ok && geoCounty == county {
				if level, ok := geo["level"].(string); ok && level == "county" {
					return &s.V6READMEs[i]
				}
			}
		}
	}

	// Try state
	if state != "" {
		for i := range s.V6READMEs {
			geo := s.V6READMEs[i].Geography
			if geoState, ok := geo["state"].(string); ok && geoState == state {
				if level, ok := geo["level"].(string); ok && level == "state" {
					return &s.V6READMEs[i]
				}
			}
		}
	}

	// Try country
	if country != "" {
		for i := range s.V6READMEs {
			geo := s.V6READMEs[i].Geography
			if geoCountry, ok := geo["country"].(string); ok && geoCountry == country {
				if level, ok := geo["level"].(string); ok && level == "country" {
					return &s.V6READMEs[i]
				}
			}
		}
	}

	return nil
}

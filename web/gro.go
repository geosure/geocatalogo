///////////////////////////////////////////////////////////////////////////////
//
// The MIT License (MIT)
// Copyright (c) 2025 GeoSure / Jeff Johnson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
// USE OR OTHER DEALINGS IN THE SOFTWARE.
//
///////////////////////////////////////////////////////////////////////////////

// Package web - GRO RESTful API
package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-spatial/geocatalogo"
	"github.com/gorilla/mux"
)

// GRORouter provides RESTful GRO API Routing
// This router provides a RESTful interface with nested, grouped endpoints
// API v1 with 5-dimensional navigation: Geography, Collections, Formats, Status, Ownership
func GRORouter(cat *geocatalogo.GeoCatalogue) *mux.Router {
	router := mux.NewRouter()

	// API v1 routes - all endpoints nested under /api/v1/
	api := router.PathPrefix("/api/v1").Subrouter()

	// Root endpoint - API info
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		GRORoot(w, r, cat)
	}).Methods("GET")

	// Search endpoint
	api.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		GROSearch(w, r, cat)
	}).Methods("GET")

	// Collections - grouped resource
	api.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		GROListCollections(w, r, cat)
	}).Methods("GET")

	api.HandleFunc("/collections/{collection}", func(w http.ResponseWriter, r *http.Request) {
		GROCollection(w, r, cat)
	}).Methods("GET")

	// Formats - grouped resource
	api.HandleFunc("/formats", func(w http.ResponseWriter, r *http.Request) {
		GROListFormats(w, r, cat)
	}).Methods("GET")

	api.HandleFunc("/formats/{format}", func(w http.ResponseWriter, r *http.Request) {
		GROFormat(w, r, cat)
	}).Methods("GET")

	// Statuses - grouped resource
	api.HandleFunc("/statuses", func(w http.ResponseWriter, r *http.Request) {
		GROListStatuses(w, r, cat)
	}).Methods("GET")

	api.HandleFunc("/statuses/{status}", func(w http.ResponseWriter, r *http.Request) {
		GROStatus(w, r, cat)
	}).Methods("GET")

	// Owners - grouped resource
	api.HandleFunc("/owners", func(w http.ResponseWriter, r *http.Request) {
		GROListOwners(w, r, cat)
	}).Methods("GET")

	api.HandleFunc("/owners/{owner}", func(w http.ResponseWriter, r *http.Request) {
		GROOwner(w, r, cat)
	}).Methods("GET")

	// Geography - explicit hierarchy with list endpoints at each level

	// Level 0: All continents
	api.HandleFunc("/geography/continents", func(w http.ResponseWriter, r *http.Request) {
		GROListContinents(w, r, cat)
	}).Methods("GET")

	// Level 1: Continent resources + countries list
	api.HandleFunc("/geography/continents/{continent}", func(w http.ResponseWriter, r *http.Request) {
		GROGeographyContinent(w, r, cat)
	}).Methods("GET")

	api.HandleFunc("/geography/continents/{continent}/countries", func(w http.ResponseWriter, r *http.Request) {
		GROListCountries(w, r, cat)
	}).Methods("GET")

	// Level 2: Country resources + states list
	api.HandleFunc("/geography/continents/{continent}/countries/{country}", func(w http.ResponseWriter, r *http.Request) {
		GROGeographyCountry(w, r, cat)
	}).Methods("GET")

	api.HandleFunc("/geography/continents/{continent}/countries/{country}/states", func(w http.ResponseWriter, r *http.Request) {
		GROListStates(w, r, cat)
	}).Methods("GET")

	// Level 3: State resources + cities list
	api.HandleFunc("/geography/continents/{continent}/countries/{country}/states/{state}", func(w http.ResponseWriter, r *http.Request) {
		GROGeographyState(w, r, cat)
	}).Methods("GET")

	api.HandleFunc("/geography/continents/{continent}/countries/{country}/states/{state}/cities", func(w http.ResponseWriter, r *http.Request) {
		GROListCities(w, r, cat)
	}).Methods("GET")

	// Level 4: City resources
	api.HandleFunc("/geography/continents/{continent}/countries/{country}/states/{state}/cities/{city}", func(w http.ResponseWriter, r *http.Request) {
		GROGeographyCity(w, r, cat)
	}).Methods("GET")

	// Records - primary resource endpoint
	api.HandleFunc("/records/{id}", func(w http.ResponseWriter, r *http.Request) {
		GRORecord(w, r, cat)
	}).Methods("GET")

	// Resources - unified query endpoint with filters
	api.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		GROResources(w, r, cat)
	}).Methods("GET")

	return router
}

// GRORoot provides catalog overview
func GRORoot(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	// Get total count
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 1, map[string]string{})

	response := map[string]interface{}{
		"api":           "gro",
		"version":       "1.0.0",
		"description":   "GRO RESTful Catalog API - Nested & Grouped 5D Navigation",
		"total_records": results.Matches,
		"base_path":     "/api/v1",
		"dimensions": []string{
			"geography",
			"collections",
			"formats",
			"statuses",
			"owners",
		},
		"endpoints": map[string]interface{}{
			"search": map[string]string{
				"path":        "/api/v1/search",
				"description": "Search with text query and filters",
				"example":     "/api/v1/search?q=wildfire&size=10",
			},
			"resources": map[string]string{
				"path":        "/api/v1/resources",
				"description": "Unified resource query with filters",
				"example":     "/api/v1/resources?collection=ai_agent&format=csv",
			},
			"collections": map[string]string{
				"list":   "/api/v1/collections",
				"detail": "/api/v1/collections/{name}",
			},
			"formats": map[string]string{
				"list":   "/api/v1/formats",
				"detail": "/api/v1/formats/{name}",
			},
			"statuses": map[string]string{
				"list":   "/api/v1/statuses",
				"detail": "/api/v1/statuses/{name}",
			},
			"owners": map[string]string{
				"list":   "/api/v1/owners",
				"detail": "/api/v1/owners/{name}",
			},
			"geography": map[string]interface{}{
				"continents": map[string]string{
					"list":      "/api/v1/geography/continents",
					"detail":    "/api/v1/geography/continents/{continent}",
					"countries": "/api/v1/geography/continents/{continent}/countries",
				},
				"countries": map[string]string{
					"detail": "/api/v1/geography/continents/{continent}/countries/{country}",
					"states": "/api/v1/geography/continents/{continent}/countries/{country}/states",
				},
				"states": map[string]string{
					"detail": "/api/v1/geography/continents/{continent}/countries/{country}/states/{state}",
					"cities": "/api/v1/geography/continents/{continent}/countries/{country}/states/{state}/cities",
				},
				"cities": map[string]string{
					"detail": "/api/v1/geography/continents/{continent}/countries/{country}/states/{state}/cities/{city}",
				},
			},
			"records": map[string]string{
				"detail": "/api/v1/records/{id}",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GROSearch handles search queries
func GROSearch(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	var q string
	var from, size int
	var collections []string

	// Parse query parameters
	query := r.URL.Query()

	if qVal := query.Get("q"); qVal != "" {
		q = qVal
	}

	if fromVal := query.Get("from"); fromVal != "" {
		from, _ = strconv.Atoi(fromVal)
	}

	if sizeVal := query.Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	} else {
		size = 10 // default
	}

	if collVal := query.Get("collections"); collVal != "" {
		collections = strings.Split(collVal, ",")
	}

	// Extract property filters
	propertyFilters := make(map[string]string)
	propertyKeys := []string{
		"continent", "country", "state", "city", "admin2",
		"collection", "type", "owner", "data_format", "status",
		"geographic_scope", "database_table", "v6_job_file",
		"v6_job_type", "s3_path", "title", "implementation_status",
	}

	for _, key := range propertyKeys {
		if val := query.Get(key); val != "" {
			propertyFilters[key] = val
		}
	}

	// Perform search
	results := cat.Search(collections, q, []float64{}, []time.Time{}, from, size, propertyFilters)

	// Return results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"matched":  results.Matches,
		"returned": len(results.Records),
		"from":     from,
		"size":     size,
		"records":  results.Records,
	})
}

// GROListContinents lists all unique continents with counts
func GROListContinents(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 10000, map[string]string{})

	continentCounts := make(map[string]int)
	for _, rec := range results.Records {
		if rec.Properties.GROMetadata != nil && rec.Properties.GROMetadata.Continent != "" {
			continentCounts[rec.Properties.GROMetadata.Continent]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":      len(continentCounts),
		"continents": continentCounts,
	})
}

// GROGeographyContinent returns resources for a specific continent
func GROGeographyContinent(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	continent := vars["continent"]

	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	propertyFilters := map[string]string{"continent": continent}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"continent": continent,
		"matched":   results.Matches,
		"returned":  len(results.Records),
		"records":   results.Records,
	})
}

// GROListCountries lists all countries for a specific continent
func GROListCountries(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	continent := vars["continent"]

	propertyFilters := map[string]string{"continent": continent}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 10000, propertyFilters)

	countryCounts := make(map[string]int)
	for _, rec := range results.Records {
		if rec.Properties.GROMetadata != nil && rec.Properties.GROMetadata.Country != "" {
			countryCounts[rec.Properties.GROMetadata.Country]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"continent": continent,
		"total":     len(countryCounts),
		"countries": countryCounts,
	})
}

// GROGeographyCountry returns resources for a specific country
func GROGeographyCountry(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	continent := vars["continent"]
	country := vars["country"]

	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	propertyFilters := map[string]string{
		"continent": continent,
		"country":   country,
	}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"continent": continent,
		"country":   country,
		"matched":   results.Matches,
		"returned":  len(results.Records),
		"records":   results.Records,
	})
}

// GROListStates lists all states for a specific country
func GROListStates(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	continent := vars["continent"]
	country := vars["country"]

	propertyFilters := map[string]string{
		"continent": continent,
		"country":   country,
	}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 10000, propertyFilters)

	stateCounts := make(map[string]int)
	for _, rec := range results.Records {
		if rec.Properties.GROMetadata != nil && rec.Properties.GROMetadata.StateProvince != "" {
			stateCounts[rec.Properties.GROMetadata.StateProvince]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"continent": continent,
		"country":   country,
		"total":     len(stateCounts),
		"states":    stateCounts,
	})
}

// GROGeographyState returns resources for a specific state
func GROGeographyState(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	continent := vars["continent"]
	country := vars["country"]
	state := vars["state"]

	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	propertyFilters := map[string]string{
		"continent": continent,
		"country":   country,
		"state":     state,
	}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"continent": continent,
		"country":   country,
		"state":     state,
		"matched":   results.Matches,
		"returned":  len(results.Records),
		"records":   results.Records,
	})
}

// GROListCities lists all cities for a specific state
func GROListCities(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	continent := vars["continent"]
	country := vars["country"]
	state := vars["state"]

	propertyFilters := map[string]string{
		"continent": continent,
		"country":   country,
		"state":     state,
	}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 10000, propertyFilters)

	cityCounts := make(map[string]int)
	for _, rec := range results.Records {
		if rec.Properties.GROMetadata != nil && rec.Properties.GROMetadata.City != "" {
			cityCounts[rec.Properties.GROMetadata.City]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"continent": continent,
		"country":   country,
		"state":     state,
		"total":     len(cityCounts),
		"cities":    cityCounts,
	})
}

// GROGeographyCity returns resources for a specific city
func GROGeographyCity(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	continent := vars["continent"]
	country := vars["country"]
	state := vars["state"]
	city := vars["city"]

	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	propertyFilters := map[string]string{
		"continent": continent,
		"country":   country,
		"state":     state,
		"city":      city,
	}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"continent": continent,
		"country":   country,
		"state":     state,
		"city":      city,
		"matched":   results.Matches,
		"returned":  len(results.Records),
		"records":   results.Records,
	})
}

// GROFormat handles format-based filtering
func GROFormat(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	format := vars["format"]

	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	propertyFilters := map[string]string{"data_format": format}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"format":   format,
		"matched":  results.Matches,
		"returned": len(results.Records),
		"records":  results.Records,
	})
}

// GROStatus handles status-based filtering
func GROStatus(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	status := vars["status"]

	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	propertyFilters := map[string]string{"implementation_status": status}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   status,
		"matched":  results.Matches,
		"returned": len(results.Records),
		"records":  results.Records,
	})
}

// GROOwner handles owner-based filtering
func GROOwner(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	owner := vars["owner"]

	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	propertyFilters := map[string]string{"owner": owner}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"owner":    owner,
		"matched":  results.Matches,
		"returned": len(results.Records),
		"records":  results.Records,
	})
}

// GROCollection handles collection-based filtering
func GROCollection(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	collection := vars["collection"]

	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	propertyFilters := map[string]string{"collection": collection}
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"collection": collection,
		"matched":    results.Matches,
		"returned":   len(results.Records),
		"records":    results.Records,
	})
}

// GRORecord retrieves a specific record by ID
func GRORecord(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	id := vars["id"]

	results := cat.Get([]string{id})

	if len(results.Records) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Record not found",
			"id":    id,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results.Records[0])
}

// GROListCollections lists all unique collections with counts
func GROListCollections(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	// Get all records
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 10000, map[string]string{})

	// Count collections
	collectionCounts := make(map[string]int)
	for _, rec := range results.Records {
		if rec.Properties.Collection != "" {
			collectionCounts[rec.Properties.Collection]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":       len(collectionCounts),
		"collections": collectionCounts,
	})
}

// GROListFormats lists all unique formats with counts
func GROListFormats(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 10000, map[string]string{})

	formatCounts := make(map[string]int)
	for _, rec := range results.Records {
		if rec.Properties.GROMetadata != nil && rec.Properties.GROMetadata.DataFormat != "" {
			formatCounts[rec.Properties.GROMetadata.DataFormat]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":   len(formatCounts),
		"formats": formatCounts,
	})
}

// GROListStatuses lists all unique implementation statuses with counts
func GROListStatuses(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 10000, map[string]string{})

	statusCounts := make(map[string]int)
	for _, rec := range results.Records {
		if rec.Properties.GROMetadata != nil && rec.Properties.GROMetadata.ImplementationStatus != "" {
			statusCounts[rec.Properties.GROMetadata.ImplementationStatus]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":    len(statusCounts),
		"statuses": statusCounts,
	})
}

// GROListOwners lists all unique owners with counts
func GROListOwners(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 10000, map[string]string{})

	ownerCounts := make(map[string]int)
	for _, rec := range results.Records {
		if rec.Properties.GROMetadata != nil && rec.Properties.GROMetadata.Owner != "" {
			ownerCounts[rec.Properties.GROMetadata.Owner]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":  len(ownerCounts),
		"owners": ownerCounts,
	})
}

// GROResources provides unified resource querying with multiple filters
// This endpoint allows combining filters: ?collection=ai_agent&format=csv&status=implemented
func GROResources(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	query := r.URL.Query()
	propertyFilters := make(map[string]string)

	// Extract all supported property filters
	propertyKeys := []string{
		"continent", "country", "state", "city", "admin2",
		"collection", "type", "owner", "data_format", "status",
		"geographic_scope", "database_table", "v6_job_file",
		"v6_job_type", "s3_path", "title", "implementation_status",
	}

	for _, key := range propertyKeys {
		if val := query.Get(key); val != "" {
			propertyFilters[key] = val
		}
	}

	// Parse pagination
	var from, size int
	if fromVal := query.Get("from"); fromVal != "" {
		from, _ = strconv.Atoi(fromVal)
	}
	if sizeVal := query.Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	} else {
		size = 100 // default
	}

	// Optional text search
	q := query.Get("q")

	// Perform search
	results := cat.Search([]string{}, q, []float64{}, []time.Time{}, from, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"filters":  propertyFilters,
		"query":    q,
		"matched":  results.Matches,
		"returned": len(results.Records),
		"from":     from,
		"size":     size,
		"records":  results.Records,
	})
}

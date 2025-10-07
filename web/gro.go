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
// This router provides a RESTful interface matching the human UI structure
// with 5-dimensional navigation: Geography, Collections, Formats, Status, Ownership
func GRORouter(cat *geocatalogo.GeoCatalogue) *mux.Router {
	router := mux.NewRouter()

	// Root endpoint - catalog info
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		GRORoot(w, r, cat)
	}).Methods("GET")

	// Search endpoint (backward compatible with CSW3)
	router.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		GROSearch(w, r, cat)
	}).Methods("GET")

	// 5D Navigation Routes
	router.HandleFunc("/geography/{continent}", func(w http.ResponseWriter, r *http.Request) {
		GROGeography(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/geography/{continent}/{country}", func(w http.ResponseWriter, r *http.Request) {
		GROGeography(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/geography/{continent}/{country}/{state}", func(w http.ResponseWriter, r *http.Request) {
		GROGeography(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/geography/{continent}/{country}/{state}/{city}", func(w http.ResponseWriter, r *http.Request) {
		GROGeography(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/format/{format}", func(w http.ResponseWriter, r *http.Request) {
		GROFormat(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/status/{status}", func(w http.ResponseWriter, r *http.Request) {
		GROStatus(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/owner/{owner}", func(w http.ResponseWriter, r *http.Request) {
		GROOwner(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/collection/{collection}", func(w http.ResponseWriter, r *http.Request) {
		GROCollection(w, r, cat)
	}).Methods("GET")

	// Record lookup by ID
	router.HandleFunc("/record/{id}", func(w http.ResponseWriter, r *http.Request) {
		GRORecord(w, r, cat)
	}).Methods("GET")

	// List endpoints
	router.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		GROListCollections(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/formats", func(w http.ResponseWriter, r *http.Request) {
		GROListFormats(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/statuses", func(w http.ResponseWriter, r *http.Request) {
		GROListStatuses(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/owners", func(w http.ResponseWriter, r *http.Request) {
		GROListOwners(w, r, cat)
	}).Methods("GET")

	return router
}

// GRORoot provides catalog overview
func GRORoot(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	// Get total count
	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, 1, map[string]string{})

	response := map[string]interface{}{
		"api":         "gro",
		"version":     "1.0.0",
		"description": "GRO RESTful Catalog API - 5-Dimensional Navigation",
		"total_records": results.Matches,
		"dimensions": []string{
			"geography",
			"collections",
			"formats",
			"statuses",
			"owners",
		},
		"endpoints": map[string]string{
			"search":      "/search?q=term",
			"geography":   "/geography/{continent}/{country}/{state}/{city}",
			"format":      "/format/{format}",
			"status":      "/status/{status}",
			"owner":       "/owner/{owner}",
			"collection":  "/collection/{collection}",
			"record":      "/record/{id}",
			"collections": "/collections",
			"formats":     "/formats",
			"statuses":    "/statuses",
			"owners":      "/owners",
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

// GROGeography handles geography-based filtering
func GROGeography(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	vars := mux.Vars(r)
	propertyFilters := make(map[string]string)

	if continent := vars["continent"]; continent != "" {
		propertyFilters["continent"] = continent
	}
	if country := vars["country"]; country != "" {
		propertyFilters["country"] = country
	}
	if state := vars["state"]; state != "" {
		propertyFilters["state"] = state
	}
	if city := vars["city"]; city != "" {
		propertyFilters["city"] = city
	}

	// Get size from query param (default 100)
	size := 100
	if sizeVal := r.URL.Query().Get("size"); sizeVal != "" {
		size, _ = strconv.Atoi(sizeVal)
	}

	results := cat.Search([]string{}, "", []float64{}, []time.Time{}, 0, size, propertyFilters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"filters":  propertyFilters,
		"matched":  results.Matches,
		"returned": len(results.Records),
		"records":  results.Records,
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

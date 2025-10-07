# GeoCatalogo Property Filters

## Overview

GeoCatalogo now supports **property-based filtering** on the API endpoint, making it easy to query catalog entries by specific attributes like continent, country, collection type, data format, and more.

## Available Property Filters

### Geographic Filters
- `continent` - Filter by continent (e.g., `clankr`, `north_america`, `south_america`)
- `country` - Filter by country (e.g., `colombia`, `united_states`)
- `state` or `state_province` - Filter by state/province
- `city` - Filter by city name
- `admin2` or `county` - Filter by admin level 2 (county)

### Metadata Filters
- `collection` - Filter by collection type (e.g., `ai_agent`, `existing_db`, `automation_bot`)
- `type` - Filter by resource type (e.g., `dataset`, `service`)
- `owner` - Filter by owner (e.g., `@clankr`, `@data-platform`)
- `data_format` - Filter by data format (e.g., `database`, `python_script`, `yaml`)
- `status` or `implementation_status` - Filter by status (e.g., `active`, `implemented`, `draft`)
- `geographic_scope` - Filter by geographic scope (e.g., `global`, `regional`)

### System Filters
- `database_table` - Filter by database table name (partial match)
- `v6_job_file` - Filter by v6 job file path (partial match)
- `v6_job_type` - Filter by v6 job type (e.g., `bootstrap`, `collector`)
- `s3_path` - Filter by S3 path (partial match)

### General Filters
- `title` - Filter by title (partial match)

## Usage Examples

### 1. Filter by Continent

Get all catalog entries for the "clankr" continent (operational infrastructure):

```bash
curl 'http://localhost:8000/?continent=clankr&maxrecords=10'
```

### 2. Filter by Collection Type

Get all AI agents:

```bash
curl 'http://localhost:8000/?collection=ai_agent'
```

Get all database tables:

```bash
curl 'http://localhost:8000/?collection=existing_db&maxrecords=50'
```

### 3. Filter by Country

Get all catalog entries for Colombia:

```bash
curl 'http://localhost:8000/?country=colombia'
```

### 4. Filter by Status

Get all active resources:

```bash
curl 'http://localhost:8000/?status=active&maxrecords=50'
```

### 5. Combine Multiple Filters

Get all active AI agents in the clankr ecosystem:

```bash
curl 'http://localhost:8000/?continent=clankr&collection=ai_agent&status=active'
```

Get all database tables with "article" in the name:

```bash
curl 'http://localhost:8000/?collection=existing_db&database_table=article'
```

### 6. Combine Property Filters with Text Search

You can combine property filters with the traditional text search (`q` parameter):

```bash
curl 'http://localhost:8000/?q=news&collection=existing_db'
```

This finds all database tables that contain "news" in their title or abstract.

### 7. Filter by Owner

Get all resources managed by the clankr agent:

```bash
curl 'http://localhost:8000/?owner=@clankr'
```

### 8. Filter by Data Format

Get all Python scripts in the catalog:

```bash
curl 'http://localhost:8000/?data_format=python_script&maxrecords=20'
```

Get all database resources:

```bash
curl 'http://localhost:8000/?data_format=database'
```

### 9. Filter by v6 Job Type

Get all v6 bootstrap jobs:

```bash
curl 'http://localhost:8000/?v6_job_type=bootstrap&maxrecords=50'
```

Get all v6 collector jobs:

```bash
curl 'http://localhost:8000/?v6_job_type=collector&maxrecords=100'
```

## Filter Matching Behavior

### Exact Match Filters
These filters require an exact match (case-insensitive):
- `continent`, `country`, `state`, `city`, `admin2`
- `collection`, `type`, `owner`, `data_format`, `status`, `geographic_scope`

### Partial Match Filters
These filters use "contains" matching (case-insensitive):
- `title`
- `database_table`
- `v6_job_file`
- `v6_job_type`
- `s3_path`

## Common Use Cases

### Discovery: Find All AI Agents

```bash
curl 'http://localhost:8000/?collection=ai_agent'
```

Response includes all 6+ AI agents in the GRO ecosystem.

### Operations: Find All Active Infrastructure

```bash
curl 'http://localhost:8000/?continent=clankr&status=active'
```

Returns all active operational resources (agents, bots, services).

### Data Engineering: Find Database Tables

```bash
curl 'http://localhost:8000/?collection=existing_db'
```

Returns all ~47 database tables with schema information.

### Job Management: Find v6 Jobs by Type

```bash
# All bootstrap jobs
curl 'http://localhost:8000/?v6_job_type=bootstrap&maxrecords=200'

# All collector jobs
curl 'http://localhost:8000/?v6_job_type=collector&maxrecords=500'

# All inference jobs
curl 'http://localhost:8000/?v6_job_type=inference&maxrecords=100'
```

### Geographic Discovery: Find Country-Specific Data

```bash
# All Colombia data
curl 'http://localhost:8000/?country=colombia'

# All California data
curl 'http://localhost:8000/?state=california'

# All Los Angeles data
curl 'http://localhost:8000/?city=los_angeles'
```

## API Response Format

Responses follow the standard GeoCatalogo format:

```json
{
  "ElapsedTime": 0,
  "Matches": 6,
  "Returned": 6,
  "NextRecord": 6,
  "Records": [
    {
      "id": "clankr_explore_agent",
      "type": "Feature",
      "properties": {
        "title": "Explore Application Agent",
        "type": "dataset",
        "collection": "ai_agent",
        "gro_metadata": {
          "continent": "clankr",
          "implementation_status": "active",
          "data_format": "go_application",
          "owner": "@explore"
        }
      }
    }
  ]
}
```

## Implementation Notes

- All filters are case-insensitive
- Multiple filters are combined with AND logic (all must match)
- Property filters can be combined with text search (`q` parameter)
- Property filters work with both CSW3 (`/`) and STAC (`/stac/search`) endpoints
- Pagination works with `maxrecords` and `startposition` parameters

## Error Handling

If no query parameters are provided, the API returns an error:

```json
{
  "Code": 20001,
  "Description": "ERROR: one of q, recordids, or property filters are required"
}
```

At least one of the following must be provided:
- `q` (text search)
- `recordids` (specific record IDs)
- One or more property filters

## Migration from Text Search

**Before (text-only search):**
```bash
curl 'http://localhost:8000/?q=clankr'
```
Returns all records containing "clankr" in title/abstract (could be 100+ matches).

**After (property filter):**
```bash
curl 'http://localhost:8000/?continent=clankr'
```
Returns only records where `gro_metadata.continent = "clankr"` (precise, 78 matches).

Property filters provide **precise, structured querying** instead of fuzzy text search.

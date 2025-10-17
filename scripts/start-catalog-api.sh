#!/bin/bash
# Catalog API Startup Script
# Starts the geocatalogo API server on port 9000 with GRO profile (memory-based)

set -e

cd /Users/jjohnson/projects/geocatalogo

# Set environment variables for GRO catalog
export GEOCATALOGO_REPOSITORY_TYPE=memory
export GEOCATALOGO_REPOSITORY_URL="file://$HOME/projects/geosure/catalog/data/geocatalogo_records.json"
export GEOCATALOGO_SERVER_URL=http://localhost:9000
export GEOCATALOGO_LOGGING_LEVEL=INFO
export GEOCATALOGO_METADATA_IDENTIFICATION_TITLE="GRO Geospatial Data Catalog (Port 9000)"
export GEOCATALOGO_METADATA_IDENTIFICATION_ABSTRACT="Global Risk Observatory unified geospatial data catalog - Catalog Agent instance"
export GEOCATALOGO_METADATA_PROVIDER_NAME="Geosure"
export GEOCATALOGO_METADATA_PROVIDER_URL="https://geosure.ai"

echo "Starting Catalog API (GRO Profile)..."
echo "Working directory: $(pwd)"
echo "Port: 9000"
echo "API: gro"
echo "Repository: memory (${GEOCATALOGO_REPOSITORY_URL})"
echo ""

# Start the geocatalogo server
exec $HOME/projects/geocatalogo/cmd/geocatalogo/geocatalogo serve -port 9000 -api gro

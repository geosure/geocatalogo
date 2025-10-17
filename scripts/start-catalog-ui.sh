#!/bin/bash
# Catalog UI Startup Script
# Starts the catalog web interface on port 3000

set -e

cd /Users/jjohnson/projects/geocatalogo

echo "Starting Catalog UI..."
echo "Working directory: $(pwd)"
echo "Port: 3000"
echo ""

# Start the catalog-ui binary
exec ./catalog-ui

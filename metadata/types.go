package metadata

// DatabaseSchema represents the complete database schema
type DatabaseSchema struct {
	Metadata struct {
		Timestamp  string `json:"timestamp"`
		Database   string `json:"database"`
		TotalTables int   `json:"total_tables"`
		TotalRows   int64 `json:"total_rows"`
	} `json:"metadata"`
	Tables []DatabaseTable `json:"tables"`
}

// DatabaseTable represents a single database table
type DatabaseTable struct {
	Name     string           `json:"name"`
	Schema   string           `json:"schema"`
	RowCount int64            `json:"row_count"`
	Size     string           `json:"size"`
	Columns  []DatabaseColumn `json:"columns"`
	Indexes  []string         `json:"indexes"`
}

// DatabaseColumn represents a database column
type DatabaseColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// CSVFile represents CSV/TSV introspection data
type CSVFile struct {
	S3Key        string      `json:"s3_key"`
	S3Path       string      `json:"s3_path"`
	Success      bool        `json:"success"`
	RowCount     int64       `json:"row_count"`
	ColumnCount  int         `json:"column_count"`
	FileSizeMB   float64     `json:"file_size_mb"`
	Schema       []CSVColumn `json:"schema"`
	ErrorMessage string      `json:"error,omitempty"`
}

// CSVColumn represents a CSV column
type CSVColumn struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
	Nullable   bool   `json:"nullable"`
}

// ParquetFile represents Parquet introspection data
type ParquetFile struct {
	S3Key       string          `json:"s3_key"`
	S3Path      string          `json:"s3_path"`
	Success     bool            `json:"success"`
	RowCount    int64           `json:"row_count"`
	ColumnCount int             `json:"column_count"`
	FileSizeMB  float64         `json:"file_size_mb"`
	Schema      []ParquetColumn `json:"schema"`
}

// ParquetColumn represents a Parquet column
type ParquetColumn struct {
	ColumnName string `json:"name"`
	DataType   string `json:"type"`
	Nullable   bool   `json:"nullable"`
}

// ShapefileFile represents Shapefile introspection data
type ShapefileFile struct {
	S3Key        string            `json:"s3_key"`
	S3Path       string            `json:"s3_path"`
	Success      bool              `json:"success"`
	FeatureCount int64             `json:"feature_count"`
	FileSizeMB   float64           `json:"file_size_mb"`
	GeometryType string            `json:"geometry_type"`
	CRS          string            `json:"srs"`
	Extent       []float64         `json:"extent"`
	Fields       []ShapefileField  `json:"schema"`
	XMLMetadata  map[string]string `json:"xml_metadata,omitempty"`
}

// ShapefileField represents a Shapefile field
type ShapefileField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Width int    `json:"width,omitempty"`
}

// GeoPackageFile represents GeoPackage introspection data
type GeoPackageFile struct {
	S3Key      string            `json:"s3_key"`
	S3Path     string            `json:"s3_path"`
	Success    bool              `json:"success"`
	FileSizeMB float64           `json:"file_size_mb"`
	Layers     []GeoPackageLayer `json:"layers"`
}

// GeoPackageLayer represents a GeoPackage layer
type GeoPackageLayer struct {
	Name         string               `json:"name"`
	FeatureCount int64                `json:"feature_count"`
	GeometryType string               `json:"geometry_type"`
	CRS          string               `json:"crs"`
	Extent       []float64            `json:"extent"`
	Fields       []GeoPackageField    `json:"fields"`
}

// GeoPackageField represents a GeoPackage field
type GeoPackageField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ExcelFile represents Excel introspection data
type ExcelFile struct {
	S3Key      string       `json:"s3_key"`
	S3Path     string       `json:"s3_path"`
	Success    bool         `json:"success"`
	FileSizeMB float64      `json:"file_size_mb"`
	Sheets     []ExcelSheet `json:"sheets"`
}

// ExcelSheet represents an Excel sheet
type ExcelSheet struct {
	Name        string   `json:"name"`
	RowCount    int      `json:"row_count"`
	ColumnCount int      `json:"column_count"`
	Headers     []string `json:"headers"`
}

// JSONFile represents JSON/GeoJSON introspection data
type JSONFile struct {
	S3Key        string            `json:"s3_key"`
	S3Path       string            `json:"s3_path"`
	Success      bool              `json:"success"`
	FileSizeMB   float64           `json:"file_size_mb"`
	FileType     string            `json:"file_type"`
	FeatureCount int               `json:"feature_count,omitempty"`
	GeometryType string            `json:"geometry_type,omitempty"`
	Properties   []string          `json:"properties,omitempty"`
	Structure    map[string]string `json:"structure,omitempty"`
}

// FileGDBFile represents File Geodatabase introspection data
type FileGDBFile struct {
	S3Key      string         `json:"s3_key"`
	S3Path     string         `json:"s3_path"`
	Success    bool           `json:"success"`
	FileSizeMB float64        `json:"file_size_mb"`
	Layers     []FileGDBLayer `json:"layers"`
}

// FileGDBLayer represents a File Geodatabase layer
type FileGDBLayer struct {
	Name         string          `json:"name"`
	FeatureCount int64           `json:"feature_count"`
	GeometryType string          `json:"geometry_type"`
	CRS          string          `json:"crs"`
	Extent       []float64       `json:"extent"`
	Fields       []FileGDBField  `json:"fields"`
}

// FileGDBField represents a File Geodatabase field
type FileGDBField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// PNGFile represents PNG introspection data
type PNGFile struct {
	S3Key           string                 `json:"s3_key"`
	S3Path          string                 `json:"s3_path"`
	Success         bool                   `json:"success"`
	FileSizeMB      float64                `json:"file_size_mb"`
	ImageType       string                 `json:"image_type"`
	ContentAnalysis map[string]interface{} `json:"content_analysis,omitempty"`
	Keywords        []string               `json:"keywords,omitempty"`
}

// PDFFile represents PDF introspection data
type PDFFile struct {
	S3Key      string                 `json:"s3_key"`
	S3Path     string                 `json:"s3_path"`
	Success    bool                   `json:"success"`
	FileSizeMB float64                `json:"file_size_mb"`
	PageCount  int                    `json:"page_count"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Structure  map[string]interface{} `json:"structure,omitempty"`
}

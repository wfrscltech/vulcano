package config

const (
	DatabaseTypePostgres = "postgres"
	DatabaseTypeMssql    = "mssql"
)

var supportedDatabaseTypes = []string{DatabaseTypePostgres, DatabaseTypeMssql}

var supportedLogLevels = []string{"debug", "info", "warning", "error"}

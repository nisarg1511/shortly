package models

import "database/sql"

type Application struct {
	Config Config
	DB     *sql.DB
}

type Config struct {
	Address string
	DBSN    string
}

type APIResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type HealthCheckResponse struct {
	Database string `json:"database"`
}

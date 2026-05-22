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

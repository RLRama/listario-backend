package models

import "gorm.io/gorm"

type DBManager struct {
	TestDB *gorm.DB
	ProdDB *gorm.DB
}

type DBConfig struct {
	DSN  string
	Name string
}

type DBVersionsResponse struct {
	TestDBVersion string `json:"test_db_version"`
	ProdDBVersion string `json:"prod_db_version"`
	TestDBError   string `json:"test_db_error,omitempty"`
	ProdDBError   string `json:"prod_db_error,omitempty"`
}

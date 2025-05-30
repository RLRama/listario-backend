package config

import (
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

var TestDB *gorm.DB
var ProdDB *gorm.DB

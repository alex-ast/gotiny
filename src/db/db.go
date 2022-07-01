package db

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
)

type DbCfg struct {
	// Type of database to use:
	//	inproc - mock database, keeps data in memory
	//	MongoDB - use MongoDB
	Type string `mapstructure:"Type"`
	// Connection string for MongoDB
	Connection string `mapstructure:"Connection"`
	// Database name
	Database string `mapstructure:"Database"`
}

type Db interface {
	Connect() error
	Store(key string, urlInfo dbmodels.Url) error
	Load(key string, urlInfo *dbmodels.Url) error
	Delete(key string) error
}


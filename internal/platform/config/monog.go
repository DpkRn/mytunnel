package config

import (
	"log"
	"os"
)

type MongoDBCfg struct {
	URI       string
	DBName    string
	ColName   string
	ColPrefix string
}

func (c config) MongoDB() MongoDBCfg {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatalf("MONGODB_URI is not set")
	}
	db := os.Getenv("MONGODB_DB")
	if db == "" {
		db = "tunnel"
	}
	col := os.Getenv("MONGODB_COLLECTION")
	if col == "" {
		col = "events"
	}
	return MongoDBCfg{
		URI:     uri,
		DBName:  db,
		ColName: col,
	}
}

func (c MongoDBCfg) URIFunc() string {
	return c.URI
}

func (c MongoDBCfg) DBNameFunc() string {
	return c.DBName
}

func (c MongoDBCfg) ColNameFunc() string {
	return c.ColName
}

func (c MongoDBCfg) ColPrefixFunc() string {
	return c.ColPrefix
}

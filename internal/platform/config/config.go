package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config interface {
	TierFunc() string
	TCPServer() TCPCfg
	HTTPServer() HTTPServerCfg
	MongoDB() MongoDBCfg
}

type config struct {
	Tier string
}

func NewConfig() Config {
	_ = godotenv.Load()
	return config{
		Tier: os.Getenv("TIER"),
	}
}

func (c config) TierFunc() string {
	return c.Tier
}

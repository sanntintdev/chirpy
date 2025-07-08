package config

import (
	"os"
	"sync/atomic"

	"github.com/sanntintdev/chirpy/internal/database"
)

type APIConfig struct {
	FileServerHitCount *atomic.Int32
	DBQueries          *database.Queries
	Platform           string
}

func NewAPIConfig(dbQueries *database.Queries) *APIConfig {
	return &APIConfig{
		FileServerHitCount: new(atomic.Int32),
		DBQueries:          dbQueries,
		Platform:           os.Getenv("PLATFORM"),
	}
}

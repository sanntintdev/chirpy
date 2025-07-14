package config

import (
	"sync/atomic"

	"github.com/sanntintdev/chirpy/internal/database"
)

type APIConfig struct {
	FileServerHitCount *atomic.Int32
	DBQueries          *database.Queries
	Platform           string
	JWT_SECRET_KEY     string
	JWT_ISSUER         string
}

type APIConfigParams struct {
	DBQueries      *database.Queries
	Platform       string
	JWT_SECRET_KEY string
	JWT_ISSUER     string
}

func NewAPIConfig(params APIConfigParams) *APIConfig {
	return &APIConfig{
		FileServerHitCount: new(atomic.Int32),
		DBQueries:          params.DBQueries,
		Platform:           params.Platform,
		JWT_SECRET_KEY:     params.JWT_SECRET_KEY,
		JWT_ISSUER:         params.JWT_ISSUER,
	}
}

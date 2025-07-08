package middleware

import (
	"net/http"

	"github.com/sanntintdev/chirpy/pkg/config"
)

func MetricsInc(cfg *config.APIConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHitCount.Add(1)
		next.ServeHTTP(w, r)
	})
}
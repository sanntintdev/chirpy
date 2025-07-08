package handlers

import (
	"fmt"
	"net/http"

	"github.com/sanntintdev/chirpy/pkg/config"
	"github.com/sanntintdev/chirpy/pkg/utils"
)

func HealthzHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
}

func MetricsHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		hits := cfg.FileServerHitCount.Load()
		hitsString := fmt.Sprintf(`
		<html>
		    <body>
		        <h1>Welcome, Chirpy Admin</h1>
		        <p>Chirpy has been visited %d times!</p>
		    </body>
		</html>
		`, hits)
		w.Write([]byte(hitsString))
	}
}

func ResetHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			utils.RespondWithErr(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
			return
		}

		if cfg.Platform != "dev" {
			utils.RespondWithErr(w, http.StatusForbidden, "Forbidden", nil)
			return
		}

		cfg.FileServerHitCount.Store(0)

		ctx := r.Context()
		if err := cfg.DBQueries.DeleteAllUsers(ctx); err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to reset database", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, "OK")
	}
}

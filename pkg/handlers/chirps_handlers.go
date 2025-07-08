package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sanntintdev/chirpy/internal/database"
	"github.com/sanntintdev/chirpy/pkg/config"
	"github.com/sanntintdev/chirpy/pkg/models"
	"github.com/sanntintdev/chirpy/pkg/utils"
)

func CreateChirpHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := models.CreateChirpRequest{}

		if err := utils.DecodeJSON(w, r, &request); err != nil {
			return
		}

		if len(request.Body) > 140 {
			utils.RespondWithErr(w, http.StatusBadRequest, "Chirp is too long", nil)
			return
		}

		cleanedText := utils.ReplaceProfanity(request.Body)
		ctx := r.Context()

		chirp, err := cfg.DBQueries.CreateChirps(ctx, database.CreateChirpsParams{
			ID:     uuid.New(),
			Body:   cleanedText,
			UserID: request.UserID,
		})

		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to create chirp", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, models.ChirpResponse{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}
}

func GetChirpsHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		chirps, err := cfg.DBQueries.GetChirps(ctx)
		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to get chirps", err)
			return
		}

		var chirpResponses []models.ChirpResponse
		for _, chirp := range chirps {
			chirpResponses = append(chirpResponses, models.ChirpResponse{
				ID:        chirp.ID,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
			})
		}
		utils.RespondWithJSON(w, http.StatusOK, chirpResponses)
	}
}

func GetChirpHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.PathValue("chirpId"))
		if err != nil {
			utils.RespondWithErr(w, http.StatusBadRequest, "Invalid ID", err)
			return
		}

		chirp, err := cfg.DBQueries.GetChirp(ctx, id)
		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to get chirp", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, models.ChirpResponse{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}
}

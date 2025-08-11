package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sanntintdev/chirpy/internal/database"
	"github.com/sanntintdev/chirpy/pkg/config"
	"github.com/sanntintdev/chirpy/pkg/middleware"
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

		// Get user ID from authenticated context
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			utils.RespondWithErr(w, http.StatusUnauthorized, "User not authenticated", nil)
			return
		}

		cleanedText := utils.ReplaceProfanity(request.Body)
		ctx := r.Context()

		chirp, err := cfg.DBQueries.CreateChirps(ctx, database.CreateChirpsParams{
			ID:     uuid.New(),
			Body:   cleanedText,
			UserID: userID,
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
			utils.RespondWithErr(w, http.StatusNotFound, "Failed to get chirp", err)
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

func DeleteChirpHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get authenticated user ID
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			utils.RespondWithErr(w, http.StatusUnauthorized, "User not authenticated", nil)
			return
		}

		id, err := uuid.Parse(r.PathValue("chirpId"))
		if err != nil {
			utils.RespondWithErr(w, http.StatusBadRequest, "Invalid ID", err)
			return
		}

		// First check if chirp exists and get its owner
		chirp, err := cfg.DBQueries.GetChirp(ctx, id)
		if err != nil {
			utils.RespondWithErr(w, http.StatusNotFound, "Chirp not found", err)
			return
		}

		// Check if the authenticated user owns this chirp
		if chirp.UserID != userID {
			utils.RespondWithErr(w, http.StatusForbidden, "You can only delete your own chirps", nil)
			return
		}

		err = cfg.DBQueries.DeleteChirp(ctx, id)
		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to delete chirp", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

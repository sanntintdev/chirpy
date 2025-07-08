package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sanntintdev/chirpy/internal/database"
	"github.com/sanntintdev/chirpy/pkg/config"
	"github.com/sanntintdev/chirpy/pkg/models"
	"github.com/sanntintdev/chirpy/pkg/utils"
)

func CreateUserHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &models.CreateUserRequest{}
		if err := utils.DecodeJSON(w, r, request); err != nil {
			return
		}

		ctx := r.Context()
		user, err := cfg.DBQueries.CreateUser(ctx, database.CreateUserParams{
			ID:    uuid.New(),
			Email: request.Email,
		})
		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to create user", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}
